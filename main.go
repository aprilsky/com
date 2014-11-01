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
	"github.com/aprilsky/goutils/model"
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
	var data = make(map[string]interface {})
	web.Get("/",func(context *webUtils.Context){
			context.Layout("layout")
			data["Title"] = "首页"
			tasks := model.ListTasks()
			data["Tasks"] = tasks
			for _,task := range tasks{
				println(task.Title)
			}
			context.Render("index",data)
		})
	web.Route("POST,GET","/list/",func(context *webUtils.Context){
			data["Title"] = "列表"
			context.Layout("layout")
			context.Render("list",data)
		})
	web.Route("GET","/add/",func(context *webUtils.Context){
			context.Layout("layout")
			context.Render("add",nil)
		})
	web.Route("POST","/add/",func(context *webUtils.Context){
			context.Json("")
		})
	web.Route("PUT","/update/",func(context *webUtils.Context){
			context.Json("")
		})
	web.Route("DELETE","/update/",func(context *webUtils.Context){
			context.Json("")
		})
	web.Get("/backup/",func(context *webUtils.Context){
			model.Storage.BackUp()
			context.Json("")
		})
	web.Run()
}
