package internal

import (
	"io/fs"
	"strconv"
)

func GetFilePerm(fm fs.FileMode) string {
	s := make([]int32, 3)
	perm := []int32(fm.Perm().String())
	m := map[int32]int{'r': 4, 'w': 2, 'x': 1}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			p := perm[i*3+j+1]
			s[i] += int32(m[p])
		}
	}
	res := "0"
	for _, v := range s {
		res += strconv.FormatInt(int64(v), 10)
	}
	return res
}
