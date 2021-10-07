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

func SplitShellCommand(cmd string) []string {
	// a b c => a b c
	// a "b" "c d" => a, b, c d
	r := newSplitShellCommand(cmd)
	findSpace := true
	findWord := 0 // 0 不到，1 word 开始，2 word 结束，也就是找 end-char
	findWordEndChar := ' '
	words := []string{}
	word := []int32{}
	for r.idx < len(r.runeList) {
		if findSpace {
			if r.dropChar(' ') {
				continue
			}
			findSpace = false
			findWord = 1
			continue
		}
		if findWord == 1 {
			if r.isChar('"') {
				findWord = 2
				findWordEndChar = '"'
			} else if r.isChar('\'') {
				findWord = 2
				findWordEndChar = '\''
			} else {
				findWord = 2
				findWordEndChar = ' '
				word = append(word, r.runeList[r.idx])
			}
			r.idx++
			continue
		}
		if findWord == 2 {
			if r.isChar('\\') {
				if r.idx+1 < len(r.runeList) {
					word = append(word, '\\', r.runeList[r.idx+1])
					r.idx += 2
				} else {
					r.idx++ // 忽略这个 \
				}
			} else {
				if r.isChar(findWordEndChar) {
					words = append(words, string(word))
					word = []int32{}
					// 一个单词结束，该 drop space
					findSpace = true
					findWord = 0
				} else {
					word = append(word, r.runeList[r.idx])
				}
				r.idx++
				continue
			}
		}
	}

	if len(word) > 0 {
		words = append(words, string(word))
	}
	return words
}

type splitShellCommand struct {
	idx      int
	runeList []int32
}

func newSplitShellCommand(cmd string) *splitShellCommand {
	return &splitShellCommand{
		idx:      0,
		runeList: []int32(cmd),
	}
}

func (r *splitShellCommand) dropChar(char int32) bool {
	if r.runeList[r.idx] == char {
		r.idx++
		return true
	}
	return false
}

func (r *splitShellCommand) isChar(char int32) bool {
	return r.runeList[r.idx] == char
}
