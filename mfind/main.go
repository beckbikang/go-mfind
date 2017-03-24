package main

import (
	"finder"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var Usage = func() {
	fmt.Println("Usage of mfind:\n\t-filename string\tfile name\n\t-path string\t\tdir path")
}

func main() {
	//目录路径
	var dirpath string
	flag.StringVar(&dirpath, "path", "", "dir path")
	if dirpath == "." {
		newPath, err := os.Getwd()
		dirpath = newPath
		if err != nil {
			log.Fatal("get current dir path faild\n")
		}
	}

	//要查找的文件名
	var filename string
	flag.StringVar(&filename, "filename", "", "file name")

	//查找文件
	var isOnlyFindType int
	flag.IntVar(&isOnlyFindType, "isfile", 0, "which type fild 0 all 1 file  2 dir ")

	//查找文件的大小
	var fileSize string
	flag.StringVar(&fileSize, "size", "", "file size like +10m +100 -10M -100")
	fileSize = strings.ToLower(strings.Trim(fileSize, " "))

	flag.Parse()
	fmt.Printf("\n#########we will find %s  from %s#########\n\n", filename, dirpath)
	if len(dirpath) > 0 {
		mf := finder.NewMfinderSimple(dirpath, filename)
		mf.IsOnlyFindType = isOnlyFindType
		//设置大小
		mf.SetFileSize(fileSize)
		mf.Run()
	} else {
		Usage()
	}

}
