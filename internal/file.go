package internal

import (
	"io/ioutil"
	"os"
	"strings"
)

func Walk(dir, dir2 string, f func(isDir bool, path, target string) error) error {
	joinPath := func(a, b string) string {
		if strings.HasSuffix(a, "/") {
			return a + b
		}
		return a + "/" + b
	}
	fileInfo, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !fileInfo.IsDir() {
		return f(false, dir, dir2)
	}
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, v := range fs {
		p1, p2 := joinPath(dir, v.Name()), joinPath(dir2, v.Name())
		if err = f(v.IsDir(), p1, p2); err != nil {
			return err
		}

		if v.IsDir() {
			if err = Walk(p1, p2, f); err != nil {
				return err
			}
		}
	}
	return nil
}
