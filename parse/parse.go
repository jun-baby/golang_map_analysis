package parse

import (
	"unsafe"
)

func ParseMap[K comparable, V any](m map[K]V) MapStruct[K, V] {
	ms := MapStruct[K, V]{}
	ms.ConstructFrom(GetMapPtr(m))
	return ms
}

func GetMapPtr[K comparable, V any](m map[K]V) unsafe.Pointer {
	// 得到map的指针，m本来就是引用，&m相当于2级指针。
	// unsafe.Pointer(&m) 得到二级指针pp
	p := unsafe.Pointer(uintptr(*(*int)(unsafe.Pointer(&m))))
	return p
}
