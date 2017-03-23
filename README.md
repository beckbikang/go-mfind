
a find tool write in golang we can use it find some file

一个简单的文本查找工具，多线程的查找文件


how to use it (如何使用它)


just use go to build it ,it don't have other dependence

you can use it in linux, i am not test it in windows 

	1 build it to mfind

	2 cp it to their path /usr/bin/mfind


for example

	mfind -h
		Usage of mfind:
		  -filename string
		        file name
		  -isfile int
		        which type fild 0 all 1 file  2 dir 
		  -path string
		        dir path


	mfind -filename=mfind -path=/Users/kang/Documents/GoProject/go-go-go/github_project/go-mfind/ -isfile=1








DONE 

	1 find filename from a dirpath
	2 Ignore case
	3 just find filename
	4 just find dirname
	
TODO 

	5 find file by size
	6 find file by type
	7 list more file info
	8 write unittest for every function