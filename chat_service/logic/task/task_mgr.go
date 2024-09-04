package task

import (
	hyentity "chat_service/entity/hunyuan_msg"
)

type Task interface {
	Exec(input string, args ...interface{}) (*hyentity.HyChatRsp, error)
	Query(param string) interface{}
}

//func InitTaskList(tasks ...Task) []Task {
//	list := make([]Task, 0)
//	for _, val := range tasks {
//		list = append(list, val)
//	}
//	return list
//}
