package main

import (
	"fmt"
	"github.com/aprilsky/goutils/convertor"
	"github.com/aprilsky/goutils/timeUtils"
	"time"
	"github.com/aprilsky/goutils/fileUtils"
	"github.com/aprilsky/goutils/stringUtils"
	webUtils "github.com/aprilsky/goutils/webUtils"
	"html/template"
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
	fullpath,err:=fileUtils.SearchFile("error.html","/Users/apple/Downloads/GoBlog-master/view/saber/error")
	fmt.Println(fullpath,err)
}

func testStringUtils(){
	str := stringUtils.Md5("000000")
	fmt.Println(str)

	b :=stringUtils.Base64Encode("hello_____world")
	str ,_=stringUtils.Base64Decode([]byte(b))
	fmt.Println(b)
	fmt.Println(str)
}
func getString(key string) string {
	return "this is a "+ key
}
func i18n(key string)interface {}{
	str:= `<i class="i18n" style="">`+key+`</i>`
	return template.HTML(str)
}
func start(){
	web := webUtils.New()

	web.View().FuncMap["getString"] = getString
	web.View().FuncMap["i18n"]=i18n

	web.Get("/admin/:id/",func(context  *webUtils.Context){
			context.Layout("admin/admin")
			context.Render("200",nil)
		})
	web.Run()
}
func main() {
	web := webUtils.New()
	web.Get("/",func(context *webUtils.Context){
			context.Layout("layout")
			context.Render("index",nil)
		})
}
