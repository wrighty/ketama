package ketama

import (
	"fmt"
	"testing"
)

func TestSetHost(t *testing.T) {

	hosts := []string{"host1", "host2"}

	c := Make(hosts)
	h1 := c.GetHost("Hello World!")
	fmt.Println(h1)
	if h1 != "host1" {
		t.Fail()
	}
	h2 := c.GetHost("Hello World")
	fmt.Println(h2)
	if h2 != "host2" {
		t.Fail()
	}

}

func TestSetHostWithWeights(t *testing.T) {
	hosts := map[string]uint{
		"host1": 79,
		"host2": 1,
	}
	c := MakeWithWeights(hosts)
	h1 := c.GetHost("Hello World!")
	if h1 != "host1" {
		t.Fail()
	}
	fmt.Println(h1)
}

func TestFindMethodsMatch(t *testing.T) {
	c := MakeWithWeights(benchmarkHosts)

	for _, key := range benchmarkKeys {
		point := c.hash(key)
		a := c.findNearestHost(point)
		b := c.findNearestHostBisect(point)
		if a != b {
			t.Log("mismatch of points", a, b)
			t.Fail()
		}
	}
}

var benchmarkKeys = []string{
	"this",
	"is",
	"a",
	"test",
	"of",
	"searches",
}

var benchmarkHosts = map[string]uint{
	"host1":  30,
	"host2":  30,
	"host3":  30,
	"host4":  30,
	"host5":  30,
	"host6":  30,
	"host7":  30,
	"host8":  30,
	"host9":  30,
	"host10": 30,
	"host11": 30,
	"host12": 30,
	"host13": 30,
	"host14": 30,
	"host15": 30,
	"host16": 30,
	"host17": 30,
	"host18": 30,
	"host19": 30,
}

func BenchmarkBisect(b *testing.B) {
	c := MakeWithWeights(benchmarkHosts)
	for i := 0; i < b.N; i++ {
		for _, key := range benchmarkKeys {
			point := c.hash(key)
			c.findNearestHostBisect(point)
		}
	}
}

func BenchmarkWalk(b *testing.B) {
	c := MakeWithWeights(benchmarkHosts)
	for i := 0; i < b.N; i++ {
		for _, key := range benchmarkKeys {
			point := c.hash(key)
			c.findNearestHost(point)
		}
	}
}
