package main

import (
	"context"
	"encoding/json"
	"fmt"

	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cfg Config

var mongoClient *mongo.Client
var sess *discordgo.Session

type Config struct {
	Variables struct {
		MonthlyPrice string `json:"monthly_price"`
		YearlyPrice  string `json:"yearly_price"`
	} `json:"variables"`
	Database struct {
		ConnectionString string `json:"connection_string"`
	} `json:"database"`
	Discord struct {
		Token               string `json:"token"`
		ServerID            string `json:"server_id"`
		MembershipChannelID string `json:"membership_channel_id"`
		MemberRole          string `json:"member_role"`
		LimitedRole         string `json:"limited_member_role"`
	} `json:"discord"`
	Stripe struct {
		APIKey string `json:"api_key"`
	} `json:"stripe"`
}

func main() {
	var err error

	// ---------Get Config---------
	file, _ := os.ReadFile("config.json")
	_ = json.Unmarshal([]byte(file), &cfg)
	//
	// ---------Connect to MongoDB---------
	//
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(cfg.Database.ConnectionString).SetServerAPIOptions(serverAPI)
	//
	// Create a new client and connect to the server
	mongoClient, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	// Defer closing the connection
	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	//
	// ---------Connect to Discord---------
	sess, err = discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		panic(err)
	}
	// clear old messages from membership channel
	msg, err := sess.ChannelMessages(cfg.Discord.MembershipChannelID, 100, "", "", "")
	if err == nil && len(msg) > 0 {
		sess.ChannelMessageDelete(cfg.Discord.MembershipChannelID, msg[0].ID)
	}
	// send new message
	message := getMembershipPrompt()
	sess.ChannelMessageSendComplex(cfg.Discord.MembershipChannelID, message)
	// add handler
	sess.AddHandler(handler)
	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	// open connection
	err = sess.Open()
	if err != nil {
		panic(err)
	}
	// defer closing the connection
	defer func() {
		msg, err := sess.ChannelMessages(cfg.Discord.MembershipChannelID, 100, "", "", "")
		if err == nil {
			sess.ChannelMessageDelete(cfg.Discord.MembershipChannelID, msg[0].ID)
		}

		sess.Close()
	}()
	// log that the bot is running
	fmt.Println("Bot is running!")

	// ---------Wait for SIGINT---------
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
