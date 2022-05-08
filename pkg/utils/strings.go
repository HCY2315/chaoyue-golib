package utils

import (
	"bytes"
	"github.com/google/uuid"
	"strings"
	"sync"
	"time"
	"unsafe"
)

var (
	_bufPool = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
)

func StringBuilder(strs ...string) (r string) {
	buf := _bufPool.Get().(*bytes.Buffer)
	for _, str := range strs {
		buf.WriteString(str)
	}
	r = buf.String()
	buf.Reset()
	_bufPool.Put(buf)
	return
}

func GetUuid() string {
	var id string
	v, err := uuid.NewRandom()
	if err != nil {
		id = time.Now().String()
	} else {
		id = v.String()
	}
	return strings.ToUpper(strings.Replace(id, "-", "", -1))
}

// It is only safe to use if you really know the byte slice won't be mutated in the lifetime of the string
func Bytes2String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

// the byte slice won't be mutated
func String2Bytes(s string) []byte {
	if len(s) == 0 {
		return []byte{}
	}
	ss := (*[2]uintptr)(unsafe.Pointer(&s))
	bs := [3]uintptr{ss[0], ss[1], ss[1]}
	return *(*[]byte)(unsafe.Pointer(&bs))
}
