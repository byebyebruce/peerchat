package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/manishmeganathan/peerchat/src"
	"github.com/sirupsen/logrus"
)

const figlet = `

W E L C O M E  T O
					     db                  db   
					     88                  88   
.8d888b. .d8888b. .d8888b. .d8888b. .d8888b. 88d888b. .d8888b. d8888P 
88'  '88 88ooood8 88ooood8 88'  '88 88'      88'  '88 88'  '88   88   
88.  .88 88.      88.      88       88.      88    88 88.  .88   88   
888888P' '88888P' '88888P' db       '88888P' db    db '8888888   '88P   
88                                                                    
dP     

`

func init() {
	// Log as Text with color
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})

	// Log to stdout
	logrus.SetOutput(os.Stdout)

	// Log Info logs and above
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	// Define input flags
	username := flag.String("user", "", "username to use in the chatroom.")
	chatroom := flag.String("room", "", "chatroom to join.")
	// Parse input flags
	flag.Parse()

	// Set background context
	ctx := context.Background()

	// Create a new P2PHost
	p2phost := src.NewP2PHost(ctx)
	// Bootstrap the DHT
	p2phost.Bootstrap()
	// Announce Service CID
	p2phost.Announce()
	// Connect to fellow Service CID providers
	p2phost.Connect()

	// Join the chat room
	chatapp, _ := src.JoinChatRoom(p2phost, *username, *chatroom)

	// Display the welcome figlet
	fmt.Print(figlet)
	time.Sleep(time.Second * 3)

	// Create the Chat UI
	ui := src.NewUI(chatapp)
	// Start the UI system
	ui.Run()
}
