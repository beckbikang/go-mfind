package finder

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type Mfinder struct {
	DirPath        string   //查找的路径
	FindName       string   //查找的文件名中含有的字段
	MaxFileSize    int64    //文件的最大
	MinFileSize    int64    //文件的最小
	listFiles      []string //文件列表
	IsOnlyFindType int      //查找 1目录 2文件
	showMore       bool
}

var KbToByte int64 = 1024
var MbToByte int64 = 1048576

var PthSep = string(os.PathSeparator)

//初始化
func NewMfinderSimple(dirpath, filename string) *Mfinder {
	//忽略大小写
	filename = strings.ToLower(filename)

	return &Mfinder{showMore: false, MaxFileSize: 0, MinFileSize: 0, DirPath: dirpath, FindName: filename}
}

func (mf *Mfinder) SetShowMore(showMore bool) {
	mf.showMore = showMore
}

//获取列表
func (mf *Mfinder) GetListFiles() []string {
	return mf.listFiles
}

func (mf *Mfinder) SetFileSize(fileSize string) {
	//has k
	var hasMore int8

	fileSize = strings.ToLower(fileSize)

	if strings.Contains(fileSize, "k") {
		hasMore = 1
		fileSize = strings.Replace(fileSize, "k", "", -1)
	}
	if strings.Contains(fileSize, "m") {
		hasMore = 2
		fileSize = strings.Replace(fileSize, "m", "", -1)
	}

	var isMax bool = false
	if strings.Contains(fileSize, "+") {
		fileSize = strings.Replace(fileSize, "+", "", -1)
	} else if strings.Contains(fileSize, "-") {
		fileSize = strings.Replace(fileSize, "-", "", -1)
		isMax = true
	}
	fileSizeInt64, _ := strconv.ParseInt(fileSize, 10, 64)
	switch hasMore {
	case 1:
		fileSizeInt64 = fileSizeInt64 * KbToByte
	case 2:
		fileSizeInt64 = fileSizeInt64 * MbToByte
	default:
	}
	if isMax {
		mf.MaxFileSize = fileSizeInt64

	} else {
		mf.MinFileSize = fileSizeInt64
	}
}

//获取文件列表
func (mf *Mfinder) ListDir(dirPth string) (err error) {
	dirSlice, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return err
	}
	var sizeIsOk bool
	for _, filename := range dirSlice {
		curFile := dirPth + PthSep + filename.Name()
		sizeIsOk = true
		//检查尺寸
		if mf.MaxFileSize > 0 {
			if filename.Size() > mf.MaxFileSize {
				sizeIsOk = false
			}
		}
		if mf.MinFileSize > 0 {
			if filename.Size() < mf.MinFileSize {
				sizeIsOk = false
			}
		}
		if filename.IsDir() {
			//全部查找或者只查找目录才放进去
			if (mf.IsOnlyFindType == 0) || (mf.IsOnlyFindType) == 2 {
				if sizeIsOk {
					mf.listFiles = append(mf.listFiles, curFile)
				}
			}
			mf.ListDir(curFile)
		} else {
			//全部查找或者只查找文件才放进去
			if (mf.IsOnlyFindType == 0) || (mf.IsOnlyFindType) == 1 {
				if sizeIsOk {
					mf.listFiles = append(mf.listFiles, curFile)
				}
			}
		}
	}
	return nil
}

//查找符合要求的文件
func (mf *Mfinder) GetRet(list []string, aChan chan []string) {
	aSlice := make([]string, 1, 100)
	for _, v := range list {
		if mf.CheckFindRet(v) == true {
			aSlice = append(aSlice, v)
		}
	}
	aChan <- aSlice
}

//是否含有子串
func (mf *Mfinder) CheckFindRet(path string) bool {
	flag := false

	if len(mf.FindName) == 0 {
		return true
	}

	path = strings.ToLower(path)
	if mf.IsOnlyFindType == 0 || mf.IsOnlyFindType == 2 {
		flag = strings.Contains(path, mf.FindName)
	} else {
		//只查找文件名
		AIndex := strings.LastIndex(path, PthSep)
		fileName := path[(AIndex + 1):]
		flag = strings.Contains(fileName, mf.FindName)
	}
	return flag
}

//获取
func (mf *Mfinder) Run() {
	//利用cpu的多核心
	CPU_NUMS := runtime.NumCPU()
	runtime.GOMAXPROCS(CPU_NUMS)

	//获取文件列表

	mf.ListDir(mf.DirPath)
	listData := mf.GetListFiles()

	//文件查找
	aListLen := len(listData)

	//调整协程数目
	var threadCount int = CPU_NUMS
	if aListLen < threadCount {
		aListLen = threadCount
	}
	step := (int)(aListLen / threadCount)

	//分配新的slice
	newSlice := make([][]string, threadCount)

	//创建携程
	chanArr := make([]chan []string, threadCount)

	//运行代码
	for i := 0; i < threadCount; i++ {
		chanArr[i] = make(chan []string)
		newSlice[i] = listData[(i * step):(i*step + step)]
		go mf.GetRet(newSlice[i], chanArr[i])
	}
	//打印找到的文件
	strList := make([]string, 1, 100)
	fmt.Println("\t#####start show result ####")
	for _, v := range chanArr {
		values := <-v
		for _, v1 := range values {
			if len(v1) == 0 {
				continue
			}
			fmt.Print("\t\t", v1)
			if mf.showMore {
				fileInfo, err := os.Stat(v1)
				if err == nil {
					fmt.Print("\t", fileInfo.Size())
					if fileInfo.IsDir() {
						fmt.Print("\t", "dir")
					} else {
						fmt.Print("\t", "not dir")
					}
				}
			}
			fmt.Println()
			strList = append(strList, v1)
		}
	}
	fmt.Println("\t#####end show result ####")
	fmt.Printf("\n#########wtotal find %d result\n", len(strList)-1)
}
