package model

import (
	"encoding/json"
	"blog/app/utils"
	"io/ioutil"
	"os"
	"path"
)

var (
	// global data storage instance
	Storage *jsonStorage
)

type jsonStorage struct {
	dir string
}

func (jss *jsonStorage) Init(dir string) {
	jss.dir = dir
}

func (jss *jsonStorage) Has(key string) bool {
	file := path.Join(jss.dir, key+".json")
	_, e := os.Stat(file)
	return e == nil
}
// TimeInc returns time step value devided by d int with time unix stamp.
func (jss *jsonStorage) TimeInc(d int) int {
	return int(utils.Now())%d + 1
}

func (jss *jsonStorage) Get(key string, v interface{}) {
	file := path.Join(jss.dir, key+".json")
	bytes, e := ioutil.ReadFile(file)
	if e != nil {
		println("read storage '" + key + "' error")
		return
	}
	e = json.Unmarshal(bytes, v)
	if e != nil {
		println("json decode '" + key + "' error")
	}
}

func (jss *jsonStorage) Set(key string, v interface{}) {
	locker.Lock()
	defer locker.Unlock()

	bytes, e := json.Marshal(v)
	if e != nil {
		println("json encode '" + key + "' error")
		return
	}
	file := path.Join(jss.dir, key+".json")
	e = ioutil.WriteFile(file, bytes, 0777)
	if e != nil {
		println("write storage '" + key + "' error")
	}
}

func (jss *jsonStorage) Dir(name string) {
	os.MkdirAll(path.Join(jss.dir, name), os.ModePerm)
}
func (jss *jsonStorage)BackUp(){
	Storage.Set("tasks",tasks)
}



// Init does model initialization.
// If first run, write default data.
// v means app.Version number. It's needed for version data.
func init() {
	Storage = new(jsonStorage)
	Storage.Init("data")

	os.Mkdir(Storage.dir, os.ModePerm)
	//writeDefaultData()
	readTasks()
}

