package gitlab

import (
	"time"
)

type PendingTaskList []PendingTask

type PendingTask struct {
	ID                      string
	ProjectID               string
	IID                     string
	WebURL                  string
	AuthorUsername          string
	LastCommentatorUsername string
	LastUpdated             time.Time
	UserNotesCount          uint
}

func (task *PendingTask) DiffersFrom(anotherTask *PendingTask) bool {
	if task == nil && anotherTask == nil {
		return false
	}
	if task == nil || anotherTask == nil {
		return true
	}
	if task.ID != anotherTask.ID {
		return true
	}
	if task.ProjectID != anotherTask.ProjectID {
		return true
	}
	if task.IID != anotherTask.IID {
		return true
	}
	if task.WebURL != anotherTask.WebURL {
		return true
	}
	if task.LastUpdated != anotherTask.LastUpdated {
		return true
	}
	if task.UserNotesCount != anotherTask.UserNotesCount {
		return true
	}
	if task.AuthorUsername != anotherTask.AuthorUsername {
		return true
	}
	if task.LastCommentatorUsername != anotherTask.LastCommentatorUsername {
		return true
	}
	return false
}

func (list PendingTaskList) DiffersFrom(anotherList []PendingTask) bool {
	if list == nil && anotherList == nil {
		return false
	}
	if len(list) != len(anotherList) {
		return true
	}
	for i, task := range list {
		if task.DiffersFrom(&anotherList[i]) {
			return true
		}
	}
	return false
}
