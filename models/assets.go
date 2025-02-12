package models

import (
	"github.com/kamva/mgm/v3"
)

type Assets struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string   `json:"name" bson:"name"`
	Type             string   `json:"type" bson:"type"`
	ParentID         string   `json:"parentId" bson:"parentId,omitempty"`
	OwnerID          string   `json:"ownerId" bson:"ownerId"`
	Childs           []string `json:"childs" bson:"childs"`
	Size             int64    `json:"size" bson:"size"`
	URL              string   `json:"url" bson:"url"`
	Path             string   `json:"path" bson:"path"`
	DeletedAt        string   `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}
