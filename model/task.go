package model

type Task struct {
	Id    int
	Title string
	CreateTime int64
	EndTime int64
	Status string
	User string
	Score float64
}
var(
	tasks  []*Task
)
func readTasks(){
	tasks = make([]*Task, 0)
	Storage.Get("tasks", &tasks)
}

func ListTasks()[]*Task{
	return tasks;
}
