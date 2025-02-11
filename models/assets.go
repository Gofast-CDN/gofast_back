package models

import (
	"github.com/kamva/mgm/v3"
)

type Assets struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name" bson:"name"`
	Type             string `json:"type" bson:"type"`
	ParentID         string `json:"parentId" bson:"parentId,omitempty"`
	OwnerID          string `json:"ownerId" bson:"ownerId"`
	Size             int64  `json:"size" bson:"size"`
	URL              string `json:"url" bson:"url"`
	Container        string `json:"container" bson:"container"`
	DeletedAt        string `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}
