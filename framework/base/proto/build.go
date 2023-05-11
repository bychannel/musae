package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

var srcDir string
var outDir string
var protoc string

func ListDir(folder string) {
	fmt.Println(folder)
	files, errDir := ioutil.ReadDir(folder)
	if errDir != nil {
		fmt.Println(errDir)
	}

	for _, file := range files {
		if file.IsDir() {
			ListDir(folder + "/" + file.Name())
		} else {
			file.Name()
			fullFilename := file.Name()
			var filenameWithSuffix string
			filenameWithSuffix = path.Base(fullFilename)
			var fileSuffix string
			fileSuffix = path.Ext(filenameWithSuffix)
			if fileSuffix != ".proto" {
				continue
			}

			//filename := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
			/*
				err := os.Mkdir(outDir+"/"+filename, 0777)
				if err != nil {
					log.Fatal(err)
					continue
				}
			*/

			fmt.Println(protoc, "--proto_path="+srcDir, "--go_out="+outDir, "--go-grpc_out=./", "--go_opt=paths=source_relative", filenameWithSuffix)
			cmd := exec.Command(protoc, "--proto_path="+srcDir, "--go_out="+outDir, "--go-grpc_out=./", "--go_opt=paths=source_relative", filenameWithSuffix)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("%+v", err.Error())
			}
			fmt.Println(string(output))

		}
	}

}

var src = flag.String("srcDir", "", "数据源")
var out = flag.String("outDir", "", "输出目录")
var proto = flag.String("protoc", "", "输出目录")

func main() {

	flag.Parse()
	srcDir = *src
	outDir = *out
	protoc = *proto
	if srcDir == "" || outDir == "" || protoc == "" {
		fmt.Println("运行参数错误，请指定输入、输出路径、protoc：")
		fmt.Println("-srcDir=./src -srcDir=./outDir")
		return
	}
	fmt.Printf("运行行参数: srcDir=%v, outDir=%v, protoc=%v", srcDir, outDir, protoc)

	//只打印某个参数
	//outDir = os.Args[1]
	//遍历参数并打印
	//for k, v := range os.Args {
	//	fmt.Printf("args[%v]=[%v]\n", k, v)
	//}

	//删除多级目录
	_ = os.RemoveAll(outDir)
	_ = os.Mkdir(outDir, 0777)
	ListDir(srcDir)
}
