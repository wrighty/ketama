package ketama

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"strconv"
	"testing"
)

func TestSetHost(t *testing.T) {

	hosts := []string{"host1", "host2"}

	tests := []struct {
		k        string
		expected string
	}{
		{"Hello World!", "host1"},
		{"Hello World", "host2"},
	}

	c := Make(hosts)
	for _, test := range tests {
		actual := c.GetHost(test.k)
		if actual != test.expected {
			t.Log("failed testing", test.k, "expected", test.expected, "got", actual)
			t.Fail()
		}
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
	tests := []struct {
		p        uint32
		expected uint32
	}{
		{0, c.points[0]},
		{4294967295, c.points[0]},
		{c.points[0], c.points[1]},
		{c.points[len(c.points)-1], c.points[0]},
		{c.points[len(c.points)-2], c.points[len(c.points)-1]},
		{c.points[len(c.points)/2], c.points[(len(c.points)/2)+1]},
	}

	for _, test := range tests {
		p1 := c.findNearestPoint(test.p)
		p2 := c.findNearestPointBisect(test.p)
		if p1 != p2 {
			t.Log("points mismatch: array walking says", p1, "bisect says", p2, "when looking up", test.p)
			t.Fail()

		}
		if p1 != test.expected {
			t.Log("did not find expected point, got", p1, "expected", test.expected)
			t.Fail()
		}
	}

}

//TestCompatibilityWithLibketama recreates the test baked into libketama to ensure this package is compatible
//with other implementations
//It is a port of https://github.com/RJ/ketama/blob/master/libketama/ketama_test.c
func TestCompatibilityWithLibketama(t *testing.T) {
	var output bytes.Buffer

	//magic value comes from
	//https://github.com/RJ/ketama/blob/18cf9a7717dad0d8106a5205900a17617043fe2c/libketama/test.sh#L8
	expected := "5672b131391f5aa2b280936aec1eea74"

	ketamaHosts := map[string]uint{
		"10.0.1.1:11211": 600,
		"10.0.1.2:11211": 300,
		"10.0.1.3:11211": 200,
		"10.0.1.4:11211": 350,
		"10.0.1.5:11211": 1000,
		"10.0.1.6:11211": 800,
		"10.0.1.7:11211": 950,
		"10.0.1.8:11211": 100,
	}
	c := MakeWithWeights(ketamaHosts)

	output.WriteString("\n") //mirrors an empty string from libketama's ketama_error()

	var host, k string
	var p, nearest uint32

	for i := 0; i < 1000000; i++ {
		k = strconv.Itoa(i)
		p = c.hash(k)
		nearest = c.findNearestPointBisect(p)
		host = c.pointsMap[nearest].name
		output.WriteString(fmt.Sprintf("%d %d %s\n", p, nearest, host))
	}
	hash := md5.Sum(output.Bytes())
	actual := fmt.Sprintf("%x", hash)

	if actual != expected {
		t.Log("this package is no longer compatible with the original libketama")
		t.Fail()
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
