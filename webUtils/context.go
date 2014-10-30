package webUtils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)



// Context instance represents a request context.
// All request and response operations are defined in this instance.
type Context struct {
	// raw *http.Request
	Request *http.Request
	// Base url, as http://domain/
	Base string
	// Path url, as http://domain/path
	Url string
	// Request url, as http://domain/path?queryString#fragment
	RequestUrl string
	// Request method, GET,POST, etc
	Method string
	// Client Ip
	Ip string
	// Client user agent
	UserAgent string
	// Last visit refer url
	Referer string
	// Request host
	Host string
	// Request url suffix
	Ext string
	// Is https
	IsSSL bool
	// Is ajax
	IsAjax bool

	// native http.ResponseWriter
	Response http.ResponseWriter
	// Response status
	Status int
	// Response header map
	Header map[string]string
	// Response body bytes
	Body []byte

	routeParams map[string]string
	flashData   map[string]interface{}

	// Response is sent or not
	IsSend bool
	// Response is end or not
	IsEnd bool

	app    *App
	layout string
}

// NewContext creates new context instance by app instance, http request and response.
func NewContext(app *App, res http.ResponseWriter, req *http.Request) *Context {

	// init context fields
	context := new(Context)
	context.flashData = make(map[string]interface{})
	context.app = app
	context.IsSend = false
	context.IsEnd = false

	// context request fields
	context.Request = req
	context.Url = req.URL.Path
	context.RequestUrl = req.RequestURI
	context.Method = req.Method
	context.Ext = path.Ext(req.URL.Path)
	context.Host = req.Host
	context.Ip = strings.Split(req.RemoteAddr, ":")[0]
	context.IsAjax = req.Header.Get("X-Requested-With") == "XMLHttpRequest"
	context.IsSSL = req.TLS != nil
	context.Referer = req.Referer()
	context.UserAgent = req.UserAgent()
	context.Base = "://" + context.Host + "/"
	if context.IsSSL {
		context.Base = "https" + context.Base
	} else {
		context.Base = "http" + context.Base
	}

	// context response fields
	context.Response = res
	context.Status = 200
	context.Header = make(map[string]string)
	context.Header["Content-Type"] = "text/html;charset=UTF-8"

	// parse form automatically
	req.ParseForm()

	return context
}

// Param returns route param by key string which is defined in router pattern string.
func (ctx *Context) Param(key string) string {
	return ctx.routeParams[key]
}

// Flash sets values to this context or gets by key string.
// The flash items are alive in this context only.
func (ctx *Context) Flash(key string, v ...interface{}) interface{} {
	if len(v) > 0 {
		return ctx.flashData[key]
	}
	ctx.flashData[key] = v[0]
	return nil
}



// Input returns all input data map.
func (ctx *Context) Input() map[string]string {
	data := make(map[string]string)
	for key, v := range ctx.Request.Form {
		data[key] = v[0]
	}
	return data
}

// Strings returns string slice of given key.
func (ctx *Context) Strings(key string) []string {
	return ctx.Request.Form[key]
}

// String returns input value of given key.
func (ctx *Context) String(key string) string {
	return ctx.Request.FormValue(key)
}


// Int returns input value of given key.
func (ctx *Context) Int(key string) int {
	str := ctx.String(key)
	i, _ := strconv.Atoi(str)
	return i
}


// Float returns input value of given key.
func (ctx *Context) Float(key string) float64 {
	str := ctx.String(key)
	f, _ := strconv.ParseFloat(str, 64)
	return f
}

// Bool returns input value of given key.
func (ctx *Context) Bool(key string) bool {
	str := ctx.String(key)
	b, _ := strconv.ParseBool(str)
	return b
}

// Cookie gets cookie value by given key when give only string.
// Cookie sets cookie value by given key, value and expire time string.
func (ctx *Context) Cookie(key string, value ...string) string {
	if len(value) < 1 {
		c, e := ctx.Request.Cookie(key)
		if e != nil {
			return ""
		}
		return c.Value
	}
	if len(value) == 2 {
		t := time.Now()
		expire, _ := strconv.Atoi(value[1])
		t = t.Add(time.Duration(expire) * time.Second)
		cookie := &http.Cookie{
			Name:    key,
			Value:   value[0],
			Path:    "/",
			MaxAge:  expire,
			Expires: t,
		}
		http.SetCookie(ctx.Response, cookie)
		return ""
	}
	return ""
}

// GetHeader returns header string by given key.
func (ctx *Context) GetHeader(key string) string {
	return ctx.Request.Header.Get(key)
}

// Redirect does redirection response to url string and status int optional.
func (ctx *Context) Redirect(url string, status ...int) {
	ctx.Header["Location"] = url
	if len(status) > 0 {
		ctx.Status = status[0]
		return
	}
	ctx.Status = 302
}

// ContentType sets content-type string.
func (ctx *Context) ContentType(contentType string) {
	ctx.Header["Content-Type"] = contentType
}

// Json set json response with data and proper header.
func (ctx *Context) Json(data interface{}) {
	bytes, e := json.MarshalIndent(data, "", "    ")
	if e != nil {
		panic(e)
	}
	ctx.ContentType("application/json;charset=UTF-8")
	ctx.Body = bytes
}

// Send does response sending.
// If response is sent, do not sent again.
func (ctx *Context) Send() {
	if ctx.IsSend {
		return
	}
	for name, value := range ctx.Header {
		ctx.Response.Header().Set(name, value)
	}
	ctx.Response.WriteHeader(ctx.Status)
	ctx.Response.Write(ctx.Body)
	ctx.IsSend = true
}

// End does end for this context.
// If context is end, handlers are stopped.
// If context response is not sent, send response.
func (ctx *Context) End() {
	if ctx.IsEnd {
		return
	}
	if !ctx.IsSend {
		ctx.Send()
	}
	ctx.IsEnd = true
}

// Throw throws http status error and error message.
// It call event named as status.
// The context will be end.
func (ctx *Context) Throw(status int, message ...interface{}) {
	//_ := strconv.Itoa(status)
	ctx.Status = status
	ctx.End()
}

// Layout sets layout string.
func (ctx *Context) Layout(str string) {
	ctx.layout = str
}

// Tpl returns string of rendering template with data.
// If error, panic.
func (ctx *Context) Tpl(tpl string, data map[string]interface{}) string {
	b, e := ctx.app.view.Render(tpl+".html", data)
	if e != nil {
		panic(e)
	}
	return string(b)
}

// Render does template and layout rendering with data.
// The result bytes are assigned to context.Body.
// If error, panic.
func (ctx *Context) Render(tpl string, data map[string]interface{}) {
	b, e := ctx.app.view.Render(tpl+".html", data)
	if e != nil {
		panic(e)
	}
	if ctx.layout != "" {
		l, e := ctx.app.view.Render(ctx.layout+".layout", data)
		//l, e := ctx.app.view.Render(ctx.layout+".layout", data)
		if e != nil {
			panic(e)
		}
		b = bytes.Replace(l, []byte("{@Content}"), b, -1)
	}
	ctx.Body = b
}

// Func adds template function to view.
// It will affect global *View instance.
func (ctx *Context) Func(name string, fn interface{}) {
	ctx.app.view.FuncMap[name] = fn
}

// App returns *App instance in this context.
func (ctx *Context) App() *App {
	return ctx.app
}

// Download sends file download response by file path.
func (ctx *Context) Download(file string) {
	f, e := os.Stat(file)
	if e != nil {
		ctx.Status = 404
		return
	}
	if f.IsDir() {
		ctx.Status = 403
		return
	}
	output := ctx.Response.Header()
	output.Set("Content-Type", "application/octet-stream")
	output.Set("Content-Disposition", "attachment; filename="+path.Base(file))
	output.Set("Content-Transfer-Encoding", "binary")
	output.Set("Expires", "0")
	output.Set("Cache-Control", "must-revalidate")
	output.Set("Pragma", "public")
	http.ServeFile(ctx.Response, ctx.Request, file)
	ctx.IsSend = true
}
