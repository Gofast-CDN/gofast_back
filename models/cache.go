package models

import "github.com/kamva/mgm/v3"

type Cache struct {
	mgm.DefaultModel `bson:",inline"`
	Content          []byte `json:"content" bson:"content"`
	ExpiresAt        string `json:"expiresAt" bson:"expiresAt"`
}
