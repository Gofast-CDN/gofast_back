package models

import "github.com/kamva/mgm/v3"

type Logs struct {
	mgm.DefaultModel `bson:",inline"`
	UserID           string `json:"userId" bson:"userId"`
	Action           string `json:"action" bson:"action"`
	TargetID         string `json:"targetId" bson:"targetId"`
	Timestamp        string `json:"timestamp" bson:"timestamp"`
}
