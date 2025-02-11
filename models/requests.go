package models

import "github.com/kamva/mgm/v3"

type Requests struct {
	mgm.DefaultModel `bson:",inline"`
	IP               string `json:"ip" bson:"ip"`
	URL              string `json:"url" bson:"url"`
	Method           string `json:"method" bson:"method"`
	Status           string `json:"status" bson:"status"`
	Timestamp        string `json:"timestamp" bson:"timestamp"`
}
