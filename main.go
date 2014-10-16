package main

import (
	"fmt"
	"github.com/aprilsky/goutils/convertor"
	"github.com/aprilsky/goutils/timeUtils"
	"time"
	"github.com/aprilsky/goutils/filetool"
	"github.com/aprilsky/goutils/stringUtils"
)
func testTimeUtils(){
	fmt.Println(timeUtils.DateFormat(time.Now(), "YYYY-MM-DD"))
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
func testStringUtils(){
	str := stringUtils.Md5("000000")
	fmt.Println(str)

	b :=stringUtils.Base64Encode("hello_____world")
	str ,_=stringUtils.Base64Decode([]byte(b))
	fmt.Println(b)
	fmt.Println(str)
}
func main() {
	testStringUtils();
}
