package main

import (
	"fmt"
	"map_visit/parse"
	"unsafe"
)

type M struct {
	a string
	b int
}

type bmap struct {
	tophash  [8]uint8
	keys     [8]int
	values   [8]int
	overflow int
}

func main() {
	m := map[string]M{
		"1": {"a", 1},
		"2": {"b", 2},
		"3": {"c", 3},
		"4": {"d", 4},
		"5": {"e", 5},
		"6": {"f", 6},
		"7": {"g", 7},
		"8": {"h", 8},
		"9": {"i", 9},
	}
	for k, v := range m {
		fmt.Println(k, v)
	}

	parseMap := parse.ParseMap(m)
	fmt.Println(parseMap)
	//b := *(*bmap)(parseMap.Buckets)
	//fmt.Println(b)

	parseMap1 := parse.ParseMap(m)
	fmt.Println(parseMap1)
	//b1 := *(*bmap)(parseMap1.Buckets)
	//fmt.Println(b1)
	//b2 := *(*bmap)(add(parseMap1.Buckets, unsafe.Sizeof(bmap{})))
	//fmt.Println(b2)

	fmt.Println(m)

	//fmt.Println(count)
}

func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}
