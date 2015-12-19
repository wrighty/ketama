package ketama

import "testing"

func TestSetHost(t *testing.T) {

	hosts := []string{"host1", "host2"}

	c := Make(hosts)
	h1 := c.GetHost("Hello World!")
	if h1 != "host1" {
		t.Fail()
	}
	h2 := c.GetHost("Hello World")
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
}

func TestFindMethodsMatch(t *testing.T) {
	c := MakeWithWeights(benchmarkHosts)

	for _, key := range benchmarkKeys {
		point := c.hash(key)
		p1 := c.findNearestPoint(point)
		p2 := c.findNearestPointBisect(point)
		if p1 != p2 {
			t.Log("points mismatch: array walking says", p1, "bisect says", p2, "when looking up", point)
			t.Fail()
		}
	}
}

func TestEdgeCases(t *testing.T) {
	c := MakeWithWeights(benchmarkHosts)
	points := []uint32{
		0,
		4294967295,
		c.points[0],
		c.points[len(c.points)-1],
		c.points[len(c.points)-2],
		c.points[len(c.points)/2],
	}

	for _, point := range points {
		p1 := c.findNearestPoint(point)
		p2 := c.findNearestPointBisect(point)
		if p1 != p2 {
			t.Log("points mismatch: array walking says", p1, "bisect says", p2, "when looking up", point)
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
	"that",
	"we",
	"try",
	"to",
	"find",
	"bugs",
	"with",
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
	var benchmarkPoints []uint32
	for _, k := range benchmarkKeys {
		benchmarkPoints = append(benchmarkPoints, c.hash(k))
	}
	for i := 0; i < b.N; i++ {
		for _, point := range benchmarkPoints {
			c.findNearestPointBisect(point)
		}
	}
}

func BenchmarkWalk(b *testing.B) {
	c := MakeWithWeights(benchmarkHosts)
	var benchmarkPoints []uint32
	for _, k := range benchmarkKeys {
		benchmarkPoints = append(benchmarkPoints, c.hash(k))
	}
	for i := 0; i < b.N; i++ {
		for _, point := range benchmarkPoints {
			c.findNearestPoint(point)
		}
	}
}
