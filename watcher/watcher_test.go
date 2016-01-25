package watcher

import (
	"io/ioutil"
	"log"
	"os"
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

	//this sleep feels bad, but we need to give fsnotify some time to spot our write
	time.Sleep(1 * time.Millisecond)

	if !r.didReload {
		t.Log("Failed to reload", filename)
		t.Fail()
	}
}

func TestWatcherWhenFileGoesAway(t *testing.T) {
	writeOriginalFile()
	defer writeOriginalFile()
	r := &basicReloadable{
		didReload: false,
	}

	w := Make(filename, r)
	deleteFile()

	//this sleep feels bad, but we need to give fsnotify some time to spot our write
	time.Sleep(1 * time.Millisecond)

	if w.watching {
		t.Log("file removed, watching should have quit")
		t.Fail()
	}
}
