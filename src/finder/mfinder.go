package finder

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

type Mfinder struct {
	DirPath        string   //查找的路径
	FindName       string   //查找的文件名中含有的字段
	MaxFileSize    uint64   //文件的最大
	MinFileSize    uint64   //文件的最小
	listFiles      []string //文件列表
	IsOnlyFindType int      //查找 1目录 2文件

}

var PthSep = string(os.PathSeparator)

//初始化
func NewMfinderSimple(dirpath, filename string) *Mfinder {
	//忽略大小写
	filename = strings.ToLower(filename)
	return &Mfinder{DirPath: dirpath, FindName: filename}
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
	for _, filename := range dirSlice {
		curFile := dirPth + PthSep + filename.Name()
		if filename.IsDir() {
			//全部查找或者只查找目录才放进去
			if (mf.IsOnlyFindType == 0) || (mf.IsOnlyFindType) == 2 {
				mf.listFiles = append(mf.listFiles, curFile)
			}
			mf.ListDir(curFile)
		} else {
			//全部查找或者只查找文件才放进去
			if (mf.IsOnlyFindType == 0) || (mf.IsOnlyFindType) == 1 {
				mf.listFiles = append(mf.listFiles, curFile)
			}
		}
	}
	return nil
}

//查找符合要求的文件
func (mf *Mfinder) GetRet(list []string, aChan chan []string) {
	aSlice := make([]string, 1, 100)
	for _, v := range list {
		if mf.checkFindRet(v) == true {
			aSlice = append(aSlice, v)
		}
	}
	aChan <- aSlice
}

//是否含有子串
func (mf *Mfinder) checkFindRet(path string) bool {
	flag := false

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
	strList := make([]string, 1, 100)
	fmt.Println("\t#####start show result ####")
	for _, v := range chanArr {
		values := <-v
		for _, v1 := range values {
			if len(v1) == 0 {
				continue
			}
			fmt.Println("\t\t", v1)
			strList = append(strList, v1)
		}
	}
	fmt.Println("\t#####end show result ####")
	fmt.Printf("\n#########wtotal find %d result\n", len(strList)-1)
}
