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
	fmt.Println("Usage of mfind:")
	fmt.Println("\t-a\n\t\tshow more info")
	fmt.Println("\t-filename string\n\t\tfile name")
	fmt.Println("\t-path string\n\t\tdir path")
	fmt.Println("\t-isfile int\n\t\twhich type fild 0 all 1 file  2 dir")
	fmt.Println("\t-size string\n\t\tfile size just like +10m +100 -10M -100")
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

	//显示更多信息
	var showMore bool
	flag.BoolVar(&showMore, "a", false, "show more info")

	flag.Parse()
	fmt.Printf("\n#########we will find %s  from %s#########\n\n", filename, dirpath)
	if len(dirpath) > 0 {
		mf := finder.NewMfinderSimple(dirpath, filename)
		mf.IsOnlyFindType = isOnlyFindType
		//设置大小
		mf.SetFileSize(fileSize)
		//显示更多
		mf.SetShowMore(showMore)
		mf.Run()
	} else {
		Usage()
	}

}
