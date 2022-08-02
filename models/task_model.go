package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	TaskName string             `json:"taskname,omitempty" validate:"required"`
	Assignee string             `json:"assignee,omitempty" validate:"required"`
	TaskDone bool               `json:"taskdone,omitempty"`
	Date     primitive.DateTime `json:"date,omitempty"`
}
