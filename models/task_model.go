package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	TaskName    string             `json:"taskname,omitempty"`
	Assignee    string             `json:"assignee,omitempty"`
	TaskDone    bool               `json:"taskdone"`
	CreatedDate time.Time          `json:"createddate"`
	DeadLine    string             `json:"deadline"`
}
