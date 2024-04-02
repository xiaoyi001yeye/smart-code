package main

type TaskStatus string

const (
	TaskStatusNew    TaskStatus = "New"    // 新建任务
	TaskStatusDoing  TaskStatus = "Doing"  // 新进行中
	TaskStatusDone   TaskStatus = "Done"   // 任务完成
	TaskStatusFailed TaskStatus = "Failed" // 任务失败
)

func StringToTaskStatus(statusStr string) TaskStatus {
	switch statusStr {
	case "New":
		return TaskStatusNew
	case "Doing":
		return TaskStatusDoing
	case "Done":
		return TaskStatusDone
	case "Failed":
		return TaskStatusFailed
	default:
		return TaskStatusNew // 默认返回新建状态
	}
}
