package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func RandUIntBelow(ceiling int) int {
	return RandIntRange(0, ceiling)
}

func RandIntRange(floor, ceiling int) int {
	return int(RandInt64Range(int64(floor), int64(ceiling)))
}

func RandInt64Range(floor, ceiling int64) int64 {
	if floor > ceiling {
		panic("blow > ceiling")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return floor + r.Int63n(ceiling-floor)
}

func RandInt64Below(ceiling int64) int64 {
	return RandInt64Range(0, ceiling)
}

func RandAscii(size int) string {
	if size < 0 {
		panic("size negative")
	}
	bs := make([]byte, size)
	for i, _ := range bs {
		bs[i] = byte(RandIntRange(32, 96))
	}
	return string(bs)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandDigitString(length int) string {
	if length <= 0 {
		panic("rand negative length")
	}
	// perf
	s := make([]string, length)
	for i := range s {
		d := RandUIntBelow(10)
		s[i] = strconv.Itoa(d)
	}
	return strings.Join(s, "")
}

func RandUInt64Below(ceiling uint64) uint64 {
	return uint64(rand.Float32() * float32(ceiling))
}

func RandLetterAndNumber(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandNoDupAsciiSet(setSize, strSize int) []string {
	m := make(map[string]bool, setSize)
	retryTimes := 10000
	for i := 0; i < setSize+retryTimes; i++ {
		s := RandAscii(strSize)
		if _, ok := m[s]; ok {
			continue
		}
		m[s] = true
		if len(m) >= setSize {
			break
		}
	}
	if len(m) < setSize {
		panic(fmt.Sprintf("after %d retry times, still dup, "+
			"setSize:%d and strSize:%d, try to incr strSize and decr setSize",
			retryTimes, setSize, strSize))
	}
	ret := make([]string, 0, setSize)
	for s := range m {
		ret = append(ret, s)
	}
	return ret
}

const (
	MaxRandBytes = 1 * MB
)

func RandBytes(length int) []byte {
	if length > MaxRandBytes {
		panic(fmt.Sprintf("RandBytes过大(>%d)", MaxRandBytes))
	}
	bs := make([]byte, length)
	rand.Read(bs)
	return bs
}
