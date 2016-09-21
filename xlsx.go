package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"path"
)

type selfxlsx struct {
	file     *xlsx.File
	sheet    *xlsx.Sheet
	row      *xlsx.Row
	cell     *xlsx.Cell
	oname    string
	odir     string
	maxlimit int
	counter  int
	first    bool
	keys     []string
	index    int
	zfiles   []string
	zOutFile string
	count    int
}

func (x *selfxlsx) getKeys(one Record) (arr []string) {
	for key, _ := range one {
		arr = append(arr, key)
	}
	return
}

func (x *selfxlsx) getName() string {
	name := fmt.Sprintf("%v_%v", x.index, x.oname)
	x.index++
	file := path.Join(x.odir, name)
	x.zfiles = append(x.zfiles, file)
	return file
}

func (x *selfxlsx) preDo() {
	if x.file == nil {
		x.reset()
	}
	x.count++
}

func (x *selfxlsx) reset() {
	x.counter = 0
	x.file = xlsx.NewFile()
	x.sheet, _ = x.file.AddSheet("k_sheet")
}

func (x *selfxlsx) addColumnName() {
	x.row = x.sheet.AddRow()
	for _, v := range x.keys {
		x.cell = x.row.AddCell()
		x.cell.SetValue(v)
	}
}

func (x *selfxlsx) Save(one Record) {
	x.preDo()
	if x.first {
		x.first = false
		x.keys = x.getKeys(one)
		x.addColumnName()
	}

	x.row = x.sheet.AddRow()
	for i := 0; i < len(x.keys); i++ {
		if _, ok := one[x.keys[i]]; !ok {
			one[x.keys[i]] = ""
		}
		x.cell = x.row.AddCell()
		x.cell.SetValue(one[x.keys[i]])
	}
	x.counter++
	if x.counter >= x.maxlimit {
		x.file.Save(x.getName())
		x.reset()
		x.addColumnName()
	}
}

func (x *selfxlsx) End() {
	if x.counter > 0 {
		x.file.Save(x.getName())
	}
	if x.count == 0 {
		x.file.Save(x.getName())
	}
	if len(x.zfiles) > 0 {
		ZipFiles(x.zfiles, x.zOutFile)
	}
}

func NewXlsx(odir string, oname string, ml int) *selfxlsx {
	base := path.Base(oname)
	zOutFile := path.Join(odir, path.Base(oname)[0:len(base)-len(path.Ext(base))]+".zip")
	return &selfxlsx{
		odir:     odir,
		oname:    oname,
		maxlimit: ml,
		counter:  0,
		first:    true,
		index:    0,
		zOutFile: zOutFile,
		count:    0,
	}
}
