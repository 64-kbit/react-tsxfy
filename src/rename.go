package main

import (
	"io/fs"
	"io/ioutil"
)

func ScanFilesInDir(dir string) []fs.FileInfo {
	var files []fs.FileInfo
	var err error
	if files, err = ioutil.ReadDir(dir); err != nil {
		panic(err)
	}
	return files
}
