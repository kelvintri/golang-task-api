package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	TaskName    string             `json:"taskname,omitempty" validate:"required"`
	Assignee    string             `json:"assignee,omitempty" validate:"required"`
	TaskDone    bool               `json:"taskdone"`
	CreatedDate time.Time          `json:"createddate"`
}
