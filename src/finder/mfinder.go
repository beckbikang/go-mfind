package finder

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

type Mfinder struct {
	DirPath     string   //查找的路径
	FindName    string   //查找的文件名中含有的字段
	MaxFileSize uint64   //文件的最大
	MinFileSize uint64   //文件的最小
	listFiles   []string //文件列表
	HasDir      bool     //是否包含dir
}

//初始化
func NewMfinderSimple(dirpath, filename string) *Mfinder {
	return &Mfinder{DirPath: dirpath, FindName: filename, HasDir: true}
}

//获取列表
func (mf *Mfinder) GetListFiles() []string {
	return mf.listFiles
}

//获取文件列表
func (mf *Mfinder) ListDir(dirPth string) (err error) {
	dirSlice, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return err
	}
	PthSep := string(os.PathSeparator)

	for _, filename := range dirSlice {
		curFile := dirPth + PthSep + filename.Name()
		if filename.IsDir() {
			if mf.HasDir {
				mf.listFiles = append(mf.listFiles, curFile)
			}
			mf.ListDir(curFile)
		} else {
			mf.listFiles = append(mf.listFiles)
		}
	}
	return nil
}

//查找符合要求的文件
func (mf *Mfinder) GetRet(list []string, aChan chan []string) {
	aSlice := make([]string, 100)
	for _, v := range list {
		if mf.checkFindRet(v) {
			aSlice = append(aSlice, v)
		}
	}
	aChan <- aSlice
}

//是否含有子串
func (mf *Mfinder) checkFindRet(path string) bool {
	flag := false
	flag = strings.Contains(path, mf.FindName)

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
	step := (int)(aListLen / CPU_NUMS)

	//分配新的slice
	newSlice := make([][]string, CPU_NUMS)

	//创建携程
	chanArr := make([]chan []string, CPU_NUMS)

	//运行代码
	for i := 0; i < CPU_NUMS; i++ {
		chanArr[i] = make(chan []string)
		newSlice[i] = listData[(i * step):(i*step + step)]
		go mf.GetRet(newSlice[i], chanArr[i])
	}
	//打印找到的文件
	strList := make([]string, 30)
	for _, v := range chanArr {
		values := <-v
		for _, v1 := range values {
			if len(v1) == 0 {
				continue
			}
			fmt.Println(v1)
			strList = append(strList, v1)
		}
	}
}
