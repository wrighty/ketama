package watcher

import (
	"io/ioutil"
	"log"
	"os"
	"runtime/debug"
	"testing"
	"time"
)

type basicReloadable struct {
	didReload bool
}

var filename = "testdata/foo"

func (r *basicReloadable) Reload() bool {
	r.didReload = true
	return false
}

func writeOriginalFile() {
	data := []byte("foo")
	write(data)
}

func alterFile() {
	data := []byte("bar")
	write(data)
}

func write(c []byte) {
	err := ioutil.WriteFile(filename, c, 0777)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteFile() {
	err := os.Remove(filename)
	if err != nil {
		log.Fatal(err)
	}
}

func TestWatcher(t *testing.T) {
	writeOriginalFile()
	defer writeOriginalFile()
	r := &basicReloadable{
		didReload: false,
	}

	w := Make(filename, r)
	if !w.watching {
		t.Log("watcher is not watching when it should be")
		t.Fail()
	}
	alterFile()

	waitForWatcher(t, w, func() bool { return r.didReload })

	if !r.didReload {
		t.Log("Failed to reload", filename)
		t.Fail()
	}
}

type waitFunc func() bool

func TestWatcherWhenFileGoesAway(t *testing.T) {
	writeOriginalFile()
	defer writeOriginalFile()
	r := &basicReloadable{
		didReload: false,
	}

	w := Make(filename, r)
	deleteFile()

	waitForWatcher(t, w, func() bool { return !w.watching })

	if w.watching {
		t.Log("file removed, watching should have quit")
		t.Fail()
	}
}

func waitForWatcher(t *testing.T, w *Watcher, f waitFunc) {
	timeout := time.After(1 * time.Millisecond)
	for {
		select {
		case <-timeout:
			t.Log(string(debug.Stack()))
			t.Log("hit timeout, check stack trace above to see more")
			return
		default:
			if f() {
				return
			}
		}
	}
}
