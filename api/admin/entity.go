package admin

import (
	"github.com/birdglove2/nitad-backend/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collectionName = database.COLLECTIONS["ADMIN"]

type Admin struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Username string             `bson:"username" json:"username" form:"username" validate:"required"`
	Password string             `bson:"password" json:"password" form:"password" validate:"required"`
}

type AdminSignup struct {
	Username        string `form:"username" validate:"required"`
	Password        string `form:"password" validate:"required"`
	ConfirmPassword string `form:"confirmPassword" validate:"required"`
}
