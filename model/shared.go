package app

// Task Queue
const MyTaskQueue1 = "MY_TASK_QUEUE_1"
const MyTaskQueue2 = "MY_TASK_QUEUE_2"

// response
type Response struct {
	Status int    `json:"status"`
	Data   string `json:"data"`
}
