package models

import (
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Email            string              `json:"email" bson:"email"`
	Password         string              `json:"password" bson:"password"`
	Role             string              `json:"role" bson:"role"`
	RootContainerID  *primitive.ObjectID `json:"rootContainerID,omitempty" bson:"rootContainerID,omitempty"`
	DeletedAt        *time.Time          `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}
