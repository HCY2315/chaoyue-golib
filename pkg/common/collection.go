package common

import "github.com/HCY2315/chaoyue-golib/pkg/utils"

func DistinctIntList(listList ...[]int) []int {
	m := make(map[int]EmptyStruct, 1*utils.KB)
	for _, list := range listList {
		for _, i := range list {
			m[i] = EmptyStruct{}
		}
	}
	ret := make([]int, 0, len(m))
	for k := range m {
		ret = append(ret, k)
	}
	return ret
}

type UIntCollection struct {
	m map[uint]EmptyStruct
}

func NewUIntCollection() *UIntCollection {
	return &UIntCollection{
		m: make(map[uint]EmptyStruct, 1*utils.KB),
	}
}

func (i *UIntCollection) AddList(l []uint) *UIntCollection {
	for _, n := range l {
		i.m[n] = EmptyStruct{}
	}
	return i
}

func (i *UIntCollection) Add(n uint) *UIntCollection {
	i.m[n] = EmptyStruct{}
	return i
}

func (i *UIntCollection) ToList() []uint {
	ret := make([]uint, 0, len(i.m))
	for n := range i.m {
		ret = append(ret, n)
	}
	return ret
}

func (i *UIntCollection) Join(other *UIntCollection) *UIntCollection {
	joined := NewUIntCollection()
	otherList := other.ToList()
	for _, u := range otherList {
		_, find := i.m[u]
		if find {
			joined.Add(u)
		}
	}
	return joined
}
