package models

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Email            string     `json:"email" bson:"email"`
	Password         string     `json:"password" bson:"password"`
	Role             string     `json:"role" bson:"role"`
	DeletedAt        *time.Time `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}
