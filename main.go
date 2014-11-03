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
	"strconv"
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


// Format unix time int64 to string
func DateInt64(ti int64, format string) string {
	t := time.Unix(int64(ti), 0)
	return timeUtils.DateFormat(t, format)
}
func i18n(key string)interface {}{
	str:= `<i class="i18n" style="">`+key+`</i>`
	return template.HTML(str)
}
func start(){
	web := webUtils.New()
	web.View().FuncMap["DateInt64"] = DateInt64
	web.View().FuncMap["i18n"]=i18n

	web.Get("/admin/:id/",func(context  *webUtils.Context){
			context.Layout("admin/admin")
			context.Render("200",nil)
		})
	web.Run()
}
func main() {
	web := webUtils.New()
	web.View().FuncMap["DateInt64"] = DateInt64

	var data = make(map[string]interface {})
	web.Get("/",func(context *webUtils.Context){
			context.Layout("layout")
			data["Title"] = "首页"
			tasks := model.ListTasks()
			data["Tasks"] = tasks

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
			inputMap := context.Input()

			println(inputMap["StartTime"])

			task := new(model.Task)
			task.Title=inputMap["Title"]
			task.StartTime,_=strconv.ParseInt(inputMap["StartTime"],0,64)
			task.EndTime,_=strconv.ParseInt((inputMap["EndTime"]),0,64)
			task.Score,_=strconv.ParseFloat(inputMap["Score"],64)
			task.User=inputMap["User"]
			model.CreateTask(task)
			context.Redirect("/")
			return
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
