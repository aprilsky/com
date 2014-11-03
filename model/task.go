package model

type Task struct {
	Id    int
	Title string
	StartTime int64
	EndTime int64
	Status string
	User string
	Score float64
}
var(
	tasks  []*Task
)
//供storage调用
func readTasks(){
	tasks = make([]*Task, 0)
	Storage.Get("tasks", &tasks)
}

func ListTasks()[]*Task{
	return tasks;
}

func CreateTask(task *Task){
	tasks = append(tasks,task)
}
