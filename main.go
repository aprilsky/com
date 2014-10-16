package main

import (
	"fmt"
	"github.com/aprilsky/goutils/convertor"
	"github.com/aprilsky/goutils/timeutils"
	"time"
	"github.com/aprilsky/goutils/filetool"
)
func testTimeUtils(){
	fmt.Println(timeutils.DateFormat(time.Now(), "YYYY-MM-DD"))
}
func testConvertorUtils(){
	b := convertor.Int64ToBytes(123)
	i :=convertor.BytesToInt64(b)
	fmt.Println(b)
	fmt.Println(i)
}
func testFileUtils(){
	i64,err:=filetool.FileToInt64("README.md")
	fmt.Println(i64,err)
}
func main() {
	testFileUtils();
}
