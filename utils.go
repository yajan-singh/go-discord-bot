package main

import (
	"context"
	"fmt"
	"log"
	"net/mail"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID            string
	Email         string
	Username      string
	Discriminator string
	exploits      int
}

func idEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func emailExists(email string, id string) bool {
	coll := mongoClient.Database("ethical").Collection("members")
	filter := bson.D{{Key: "email", Value: email}}
	var result bson.M
	err := coll.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if result["email"] == email && result["id"] != id {
		fmt.Println("Email already exists")
		return true
	}
	return false
}

func getEmail(id string) string {
	coll := mongoClient.Database("ethical").Collection("members")
	filter := bson.D{{Key: "id", Value: id}}
	var result bson.M
	err := coll.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return result["email"].(string)
}

func insertToDb(id string, email string, username string, discriminator string) {
	coll := mongoClient.Database("ethical").Collection("members")
	filter := bson.D{{Key: "id", Value: id}}
	opts := options.Delete().SetHint(bson.D{{Key: "_id", Value: 1}})
	result, err := coll.DeleteMany(context.TODO(), filter, opts)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Number of documents deleted: %d\n", result.DeletedCount)

	usr := User{
		ID:            id,
		Email:         email,
		Username:      username,
		Discriminator: discriminator,
		exploits:      0,
	}

	results, err := coll.InsertOne(context.TODO(), usr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated document with _id: %v\n", results.InsertedID)
}

func makeMember(id string) bool {
	_ = sess.GuildMemberRoleRemove(cfg.Discord.ServerID, id, cfg.Discord.LimitedRole)
	err := sess.GuildMemberRoleAdd(cfg.Discord.ServerID, id, cfg.Discord.MemberRole)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
