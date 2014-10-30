package webUtils

import (
	"bytes"
	"html/template"
	"path"
	"strings"
)

// View instance provides simple template render.
type View struct {
	// template directory
	Dir string
	// view functions map
	FuncMap template.FuncMap
}

func (v *View) getTemplateInstance(tpl []string) (*template.Template, error) {
	var (
		t    *template.Template
		e    error
		file []string = make([]string, len(tpl))
	)
	for i, tp := range tpl {
		file[i] = path.Join(v.Dir, tp)
	}
	t = template.New(path.Base(tpl[0]))
	t.Funcs(v.FuncMap)
	t, e = t.ParseFiles(file...)
	if e != nil {
		return nil, e
	}
	return t, nil
}

// Render renders template with data.
// Tpl is the file names under template directory, like tpl1,tpl2,tpl3.
func (v *View) Render(tpl string, data map[string]interface{}) ([]byte, error) {
	t, e := v.getTemplateInstance(strings.Split(tpl, ","))
	if e != nil {
		return nil, e
	}
	var buf bytes.Buffer
	e = t.Execute(&buf, data)
	if e != nil {
		return nil, e
	}

	/*t, err := template.ParseFiles("template/html/admin/index.html")
	if (err != nil) {
		println(err)
	}
	var buf bytes.Buffer
	t.Execute(w, &User{user})*/
	return buf.Bytes(), nil
}


// NewView returns view instance with directory.
// It contains bundle template function HTML(convert string to template.HTML).
func NewView(dir string) *View {
	v := new(View)
	v.Dir = dir
	v.FuncMap = make(template.FuncMap)
	v.FuncMap["Html"] = func(str string) template.HTML {
		return template.HTML(str)
	}
	return v
}
