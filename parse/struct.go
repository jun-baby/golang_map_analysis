package parse

import "unsafe"

const (
	PtrSize = 4 << (^uintptr(0) >> 63)
)

type MapStruct[K comparable, V any] struct {
	count      int // # live cells == size of map.  Must be first (used by len() builtin)
	flags      uint8
	B          uint8           // log_2 of # of Buckets (can hold up to loadFactor * 2^B items)
	noverflow  uint16          // approximate number of overflow Buckets; see incrnoverflow for details
	hash0      uint32          // hash seed
	Buckets    []*bmap[K, V]   // array of 2^B Buckets. may be nil if count==0.
	oldbuckets []*bmap[K, V]   // previous bucket array of half the size, non-nil only when growing
	nevacuate  uintptr         // progress counter for evacuation (Buckets less than this have been evacuated)
	extra      *mapextra[K, V] // optional fields
}

func (s *MapStruct[K, V]) ConstructFrom(m unsafe.Pointer) {
	s.parseCount(m)
	s.parseFlags(m)
	s.parseB(m)
	s.parseNoverflow(m)
	s.parseHash0(m)
	s.parseBuckets(m)
	s.parseOldbuckets(m)
	s.parseNevacuate(m)
	s.parseExtra(m)
}

func (s *MapStruct[K, V]) parseCount(m unsafe.Pointer) {
	s.count = *(*int)(add(m, unsafe.Offsetof((*s).count)))
}

func (s *MapStruct[K, V]) parseFlags(m unsafe.Pointer) {
	s.flags = *(*uint8)(add(m, unsafe.Offsetof((*s).flags)))
}

func (s *MapStruct[K, V]) parseB(m unsafe.Pointer) {
	s.B = *(*uint8)(add(m, unsafe.Offsetof((*s).B)))
}

func (s *MapStruct[K, V]) parseNoverflow(m unsafe.Pointer) {
	s.noverflow = *(*uint16)(add(m, unsafe.Offsetof((*s).noverflow)))
}

func (s *MapStruct[K, V]) parseHash0(m unsafe.Pointer) {
	s.hash0 = *(*uint32)(add(m, unsafe.Offsetof((*s).hash0)))
}

func (s *MapStruct[K, V]) parseBuckets(m unsafe.Pointer) {
	if s.B == 0 {
		s.parseB(m)
	}
	bucketCount := 1 << s.B
	bucketSize := s.BucketSize()
	bucketOffset := unsafe.Offsetof((*s).Buckets)
	bucketPtr := unsafe.Pointer(uintptr(*(*int)(add(m, bucketOffset))))
	buckets := make([]*bmap[K, V], bucketCount)
	for i := 0; i < bucketCount; i++ {
		buckets[i] = (*bmap[K, V])(add(bucketPtr, bucketSize*uintptr(i)))
	}

	s.Buckets = buckets
}

func (s *MapStruct[K, V]) parseOldbuckets(m unsafe.Pointer) {
	if s.B == 0 {
		s.parseB(m)
	}
	bucketOffset := unsafe.Offsetof((*s).oldbuckets)
	bucketPtr := unsafe.Pointer(uintptr(*(*int)(add(m, bucketOffset))))
	if uintptr(bucketPtr) == 0 {
		return
	}
	bucketCount := 1 << (s.B - 1)
	bucketSize := s.BucketSize()
	buckets := make([]*bmap[K, V], bucketCount)
	for i := 0; i < bucketCount; i++ {
		buckets[i] = (*bmap[K, V])(add(bucketPtr, bucketSize*uintptr(i)))
	}

	s.oldbuckets = buckets
}

func (s *MapStruct[K, V]) parseNevacuate(m unsafe.Pointer) {
	s.nevacuate = *(*uintptr)(add(m, unsafe.Offsetof((*s).nevacuate)))
}

func (s *MapStruct[K, V]) parseExtra(m unsafe.Pointer) {
	off := add(m, unsafe.Offsetof((*s).extra))
	pointer := unsafe.Pointer(uintptr(*(*int)(off)))
	s.extra = (*mapextra[K, V])(pointer)
}

func (s *MapStruct[K, V]) KeySize() int {
	k := new(K)
	return int(unsafe.Sizeof(*k))
}

func (s *MapStruct[K, V]) ValueSize() int {
	v := new(V)
	return int(unsafe.Sizeof(*v))
}

func (s *MapStruct[K, V]) BucketSize() uintptr {
	return unsafe.Sizeof(bmap[K, V]{})
}

func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

type mapextra[K comparable, V any] struct {
	Overflow     *[]*bmap[K, V]
	Oldoverflow  *[]*bmap[K, V]
	NextOverflow *bmap[K, V]
}
type bmap[K comparable, V any] struct {
	Tophash  [8]uint8
	Keys     [8]K
	Vals     [8]V
	Overflow *bmap[K, V]
}
