package models

import (
	"context"
	"fmt"
	"time"

	"github.com/ivinayakg/hypergroai_assignment/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func (u *User) Save() error {
	createdAt := time.Now()

	u.CreatedAt = createdAt

	ctx := context.TODO()

	res, err := helpers.CurrentDb.User.InsertOne(ctx, u)
	if err != nil {
		fmt.Println(err)
		return err
	}

	u.ID = res.InsertedID.(primitive.ObjectID)
	fmt.Printf("User created with id %v\n", u.ID)

	return nil
}

func GetUser(email string) (*User, error) {
	user := new(User)

	ctx := context.TODO()
	userFilter := bson.M{"email": email}

	err := helpers.CurrentDb.User.FindOne(ctx, userFilter).Decode(user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Printf("User found with id %v\n", user.ID)
	return user, nil
}
