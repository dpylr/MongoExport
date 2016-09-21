// MongoExport project main.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	//	"path"
	//	"strings"
	"path/filepath"
)

type Record map[string]interface{}

var mgourl = flag.String("mgourl", "localhost", "-mgourl=localhost mongo 地址")
var mgodb = flag.String("mgodb", "storeServer", "-mgourl=storeServer 要访问的 db name")
var mgoc = flag.String("mgoc", "userrecords", "-mgoc=userrecords 要访问的 collections name")

//var mgom = flag.String("mgom", "Find", "-mgom=Find 要调用的方法, 目前只支持 Find")

//var intype = flag.String("intype", "File", "-intype=File|Pipe Pipe 表示使用管道方式读取查询条件,File 表示使用文件,如果使用文件,必须使用参数 -afile")
var afile = flag.String("afile", "in.txt", "-afile=in.txt 参数输入源文件")
var dfile = flag.String("dfile", "out.xlsx", "-dfile=out.xlsx 格式化后的输出文件")
var dcount = flag.Int("dcount", 500000, "-dcount=500000 假如查询出的记录大于 50W,那么就以50W分一个文件进行保存")

//var fkey = flag.String("fkey", "data", "-fkey=data 需要的关键字 默认只会将mongo数据的data字段内容进行导出,多个字段使用 | 分割")
var counter int = 0

func readCond() bson.M {
	file, err := os.Open(*afile)
	if err != nil {
		return bson.M{}
	}
	defer file.Close()
	d := json.NewDecoder(file)
	var arg bson.M
	err = d.Decode(&arg)
	if err != nil {
		return bson.M{}
	}
	return arg
}

func export() error {
	s, err := mgo.Dial(*mgourl)
	if err != nil {
		return err
	}
	s.SetMode(mgo.Monotonic, true)
	defer s.Close()
	pipe := s.DB(*mgodb).C(*mgoc).Pipe([]bson.M{{"$match": readCond()}})
	iter := pipe.Iter()
	defer iter.Close()
	var one Record
	dir := filepath.Dir(filepath.Clean(*dfile))
	name := filepath.Base(filepath.Clean(*dfile))
	xlsx := NewXlsx(dir, name, *dcount)
	for iter.Next(&one) {
		one, _ := one["data"].(Record)
		xlsx.Save(one)
	}
	xlsx.End()
	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()
	fmt.Println(os.Args)
	err := export()
	if err != nil {
		fmt.Errorf("%v\n", err.Error())
		panic(err)
	}
}
