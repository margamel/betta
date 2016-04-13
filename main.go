// This file provides a basic "quick start" example of using the Discordgo
// package to connect to Discord using the New() helper function.
package main

import (
	"fmt"
	"os"
	//"strconv"
	//"strings"
	//"time"

	"github.com/bwmarrin/discordgo"
)

var ds *discordgo.Session

func main() {

	// Check for Username and Password CLI arguments.
	if len(os.Args) != 3 {
		fmt.Println("You must provide username and password as arguments. See below example.")
		fmt.Println(os.Args[0], " [username] [password]")
		return
	}
	if _, err := os.Stat("bank/"); os.IsNotExist(err) {
		// path/to/whatever does not exist
		os.Mkdir("bank/", 644)
		fmt.Println("Making bank directory for you.")
	}

	// Call the helper function New() passing username and password command
	// line arguments. This returns a new Discord session, authenticates,
	// connects to the Discord data websocket, and listens for events.
	dg, err := discordgo.New(os.Args[1], os.Args[2])
	ds = dg
	if err != nil {
		fmt.Println(err)
		return
	}

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(HandleMessage)

	// Open the websocket and begin listening.
	dg.Open()

	// Simple way to keep program running until any key press.
	var input string
	fmt.Scanln(&input)
	return
}
