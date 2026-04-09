package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wallet struct {
	USD         float64 `bson:"usd" json:"usd"`
	TotalTrades int     `bson:"total_trades" json:"total_trades"`
}

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email" json:"email" binding:"required,email"`
	PasswordHash string             `bson:"password_hash" json:"-"` //omit the feild for json beause when gin serialize this struct to an http response it will ignore this entry because in response password should not be returned to the client/frontend.
	Role         string             `bson:"role" json:"role"`
	Wallet       Wallet             `bson:"wallet" json:"wallet"`
	Created      time.Time          `bson:"created_at" json:"created_at"`
	Updated      time.Time          `bson:"updated_at" json:"updated_at"`
}
