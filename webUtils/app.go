package webUtils

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
)


type App struct {
	router  *Router
	view    *View
	middle  []Handler
	config  *Config
	inter   map[string]Handler
}


func New() *App {
	a := new(App)
	a.router = NewRouter()
	a.middle = make([]Handler, 0)
	a.config, _ = NewConfig("config.json")
	a.view = NewView(a.config.StringOr("app.view_dir", "view"))
	a.inter = make(map[string]Handler)

	notFoundHandle := func(context *Context) {
		context.Body = append([]byte("<pre>"), context.Body...)
		context.Body = append(context.Body, []byte("the page you find is 404")...)
		context.Body = append(context.Body, []byte("</pre>")...)
		context.Send()
	}
	a.NotFound(notFoundHandle)


	return a
}

// Use adds middleware handlers.
// Middleware handlers invoke before route handler in the order that they are added.
func (app *App) Use(h ...Handler) {
	app.middle = append(app.middle, h...)
}

// Config returns global *Config instance.
func (app *App) Config() *Config {
	return app.config
}

// View returns global *View instance.
func (app *App) View() *View {
	return app.view
}

func (app *App) handler(res http.ResponseWriter, req *http.Request) {
	context := NewContext(app, res, req)

	defer func() {
		e := recover()
		if e == nil {
			context = nil
			return
		}
		context.Body = []byte(fmt.Sprint(e))
		context.Status = 503
		debug.PrintStack()
		if _, ok := app.inter["recover"]; ok {
			app.inter["recover"](context)
		}
		if !context.IsEnd {
			context.End()
		}
		context = nil
	}()

	if _, ok := app.inter["static"]; ok {
		app.inter["static"](context)
		if context.IsEnd {
			return
		}
	}

	if len(app.middle) > 0 {
		for _, h := range app.middle {
			h(context)
			if context.IsEnd {
				break
			}
		}
	}

	if context.IsSend {
		return
	}
	var (
		params map[string]string
		fn     []Handler
		url    = req.URL.Path
	)

	params, fn = app.router.Find(url, req.Method)

	if params != nil && fn != nil {
		context.routeParams = params

		for _, f := range fn {
			f(context)
			if context.IsEnd {
				break
			}
		}
		if !context.IsEnd {
			context.End()
		}
	} else {
		println("router is missing at " + req.URL.Path)
		context.Status = 404
		if notFoundHandler, ok := app.inter["notfound"]; ok {
			notFoundHandler(context)
			if !context.IsEnd {
				context.End()
			}
		} else {
			context.Throw(404)
		}
	}

	context = nil
}

// ServeHTTP is HTTP server implement method. It makes App compatible to native http handler.
func (app *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	app.handler(res, req)
}

// Run http server and listen on config value or 9001 by default.
func (app *App) Run() {
	addr := app.config.StringOr("app.server", "localhost:9001")
	println("http server run at " + addr)
	e := http.ListenAndServe(addr, app)
	panic(e)
}

// Set app config value.
func (app *App) Set(key string, v interface{}) {
	app.config.Set("app."+key, v)
}

// Get app config value if only key string given, return string value.
// If fn slice given, register GET handlers to router with pattern string.
func (app *App) Get(key string, fn ...Handler) string {
	if len(fn) > 0 {
		app.router.Get(key, fn...)
		return ""
	}
	return app.config.String("app." + key)
}

// Register POST handlers to router.
func (app *App) Post(key string, fn ...Handler) {
	app.router.Post(key, fn...)
}

// Register PUT handlers to router.
func (app *App) Put(key string, fn ...Handler) {
	app.router.Put(key, fn...)
}

// Register DELETE handlers to router.
func (app *App) Delete(key string, fn ...Handler) {
	app.router.Delete(key, fn...)
}

// Register handlers to router with custom methods and pattern string.
// Support GET,POST,PUT and DELETE methods.
// Usage:
//     app.Route("GET,POST","/test",handler)
//
func (app *App) Route(method string, key string, fn ...Handler) {
	methods := strings.Split(method, ",")
	for _, m := range methods {
		switch m {
		case "GET":
			app.Get(key, fn...)
		case "POST":
			app.Post(key, fn...)
		case "PUT":
			app.Put(key, fn...)
		case "DELETE":
			app.Delete(key, fn...)
		default:
			println("unknow route method " + m)
		}
	}
}

// Register static file handler.
// It's invoked before route handler after middleware handler.
func (app *App) Static(h Handler) {
	app.inter["static"] = h
}

// Register panic recover handler.
// It's invoked when panic error in middleware and route handlers.
func (app *App) Recover(h Handler) {
	app.inter["recover"] = h
}

// Register NotFound handler.
// It's invoked after calling route handler but not matched.
func (app *App) NotFound(h Handler) {
	app.inter["notfound"] = h
}
