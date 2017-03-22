package main

import (
	"finder"
	"flag"
	"fmt"
	"log"
	"os"
)

var Usage = func() {
	fmt.Println("USAGE: mfind command [arguments] ...")
	fmt.Println("\nThe commands are:\n\tpath\t add file path.\n\tfilename\tthe filename to be matched.")
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
	flag.Parse()
	fmt.Println(dirpath, filename)
	if len(dirpath) > 0 && len(filename) > 0 {
		mf := finder.NewMfinderSimple(dirpath, filename)
		mf.Run()
	} else {
		Usage()
	}
}
