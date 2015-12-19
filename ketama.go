package ketama

import (
	"crypto/md5"
	"math"
	"sort"
	"strconv"
)

//Continum models a sparse space that a given host occupies one or more points along
type Continuim struct {
	pointsMap map[uint32]*host
	points    []uint32
}

type host struct {
	name string
}

//BySize implements sort.Interface for unit32
type BySize []uint32

func (a BySize) Len() int           { return len(a) }
func (a BySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySize) Less(i, j int) bool { return a[i] < a[j] }

//GetHost looks up the host that is next nearest on the Continuim to where key hashes to and returns the name
func (c *Continuim) GetHost(key string) string {
	point := c.hash(key)

	//nearest := c.findNearestPoint(point)
	nearest := c.findNearestPointBisect(point)
	h := c.pointsMap[nearest]
	return h.name
}

func (c *Continuim) findNearestPoint(point uint32) uint32 {
	//this is hideous linear walk through the array to find the first biggest point
	var firstBiggest uint32
	for _, p := range c.points {
		if p > point {
			firstBiggest = p
			break
		}
	}
	//check for point that is outside Continuim
	//we tried every point and found nothing bigger, therefore wrap
	if firstBiggest == 0 {
		firstBiggest = c.points[0]
	}
	return firstBiggest
}

func (c *Continuim) findNearestPointBisect(point uint32) uint32 {
	//this is a hideous binary search through the array to find the first biggest point
	jump := len(c.points) / 2
	var firstBiggest uint32
	curPos := jump
	for {
		jump = jump / 2
		if jump == 0 {
			jump = 1
		}
		cmp := c.points[curPos]

		if cmp < point {
			curPos += jump
			if curPos >= len(c.points) {
				firstBiggest = c.points[0]
				break
			}
			continue
		}
		if cmp == point {
			if curPos < (len(c.points) - 1) {
				firstBiggest = c.points[curPos+1]
				break
			}
			firstBiggest = c.points[0]
			break
		}
		if cmp > point {
			if curPos == 0 {
				firstBiggest = c.points[0]
				break
			}
			if c.points[curPos-1] > point {
				curPos -= jump
				continue
			}
			firstBiggest = cmp
			break
		}
	}
	return firstBiggest
}

//Make instantiates a Continuim with a set of hosts of equal weights
func Make(hosts []string) *Continuim {
	c := &Continuim{}
	c.setHosts(hosts)
	return c
}

//MakeWithWeights instantiates a Continuim with a set of hosts, each with an explicit weighting
func MakeWithWeights(hosts map[string]uint) *Continuim {
	c := &Continuim{}
	c.setHostsWithWeights(hosts)
	return c
}

//setHosts is a convience function when you have hosts with equal weight
func (c *Continuim) setHosts(hosts []string) {
	var weights = make(map[string]uint)
	for _, host := range hosts {
		weights[host] = 1
	}
	c.setHostsWithWeights(weights)
}

//setHostsWithWeights mirrors the java implementation of ketama and uses each host's relative weight to determine the number of points it occupies across the continuim
// https://github.com/RJ/ketama/blob/18cf9a7717dad0d8106a5205900a17617043fe2c/java_ketama/SockIOPool.java#L587-L607
func (c *Continuim) setHostsWithWeights(hostnames map[string]uint) {
	c.pointsMap = make(map[uint32]*host)
	var totalWeight uint

	for _, weight := range hostnames {
		totalWeight += weight
	}

	for hostname, weight := range hostnames {
		h := &host{
			name: hostname,
		}
		factor := int(math.Floor((40 * float64(len(hostnames)) * float64(weight)) / float64(totalWeight)))
		//fmt.Println(factor)
		for i := 0; i < factor; i++ {
			key := h.name + "-" + strconv.Itoa(i)
			sum := md5.Sum([]byte(key))
			for j := 0; j < 4; j++ {
				point := uint32(sum[3+j*4])<<24 |
					uint32(sum[2+j*4])<<16 |
					uint32(sum[1+j*4])<<8 |
					uint32(sum[j*4])
				c.pointsMap[point] = h
				//fmt.Println("added point", point, "for host", h.name, "using key", key)
			}
		}
	}
	//use the points in the map to construct a sorted array to search later in GetHost
	for point, _ := range c.pointsMap {
		c.points = append(c.points, point)
		sort.Sort(BySize(c.points))
	}
}

func (c Continuim) hash(key string) uint32 {
	sum := md5.Sum([]byte(key))

	return uint32(sum[3])<<24 |
		uint32(sum[2])<<16 |
		uint32(sum[1])<<8 |
		uint32(sum[0])
}
