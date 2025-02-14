package models

import (
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Assets struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string              `json:"name" bson:"name"`
	Type             string              `json:"type" bson:"type"`
	ParentID         *primitive.ObjectID `json:"parentId,omitempty" bson:"parentId,omitempty"`
	OwnerID          primitive.ObjectID  `json:"ownerId" bson:"ownerId"`
	Childs           []Assets            `json:"childs,omitempty" bson:"childs,omitempty"`
	Size             int64               `json:"size" bson:"size"`
	Depth            int64               `json:"depth" bson:"depth"`
	URL              string              `json:"url" bson:"url"`
	Path             string              `json:"path" bson:"path"`
	DeletedAt        *time.Time          `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}
