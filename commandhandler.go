package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
	"time"
)

var prefix = "-="

type Command interface {
	Base() string                                                                           // Returns the base for this command
	Run(session *discordgo.Session, msg *discordgo.Message, split []string, isPrivate bool) // Runs this command
	Help(specific bool) string                                                              // Called when help, specific is true when we called help on this specific command "!help somecommand"
}

// One command
type CommandHelp struct{}

func (c *CommandHelp) Base() string { return "help" }
func (c *CommandHelp) Run(s *discordgo.Session, m *discordgo.Message, split []string, isPrivate bool) {
	// run this help command
	// You can loop over all the commands in the commands slice and call "Help()" on each of them
	// To list help on all commands here for example, or a specific one if we provided one
	helpmsg := "We have these commands; \n"
	for _, cmd := range commands {
		helpmsg += fmt.Sprintf("%v | %10v\n", cmd.Base(), cmd.Help(false))

	}
	s.ChannelMessageSend(m.ChannelID, helpmsg)
}
func (c *CommandHelp) Help(specific bool) string {
	return "You stupid i smart"
}

// Another command
type CommandEcho struct{}

func (c *CommandEcho) Base() string { return "echo" }
func (c *CommandEcho) Run(s *discordgo.Session, m *discordgo.Message, split []string, isPrivate bool) {
	// Echo what was sent to us
	s.ChannelMessageSend(m.ChannelID, m.ContentWithMentionsReplaced())
}

func (c *CommandEcho) Help(specific bool) string {
	return "Make me say stuff ;)"
}

//slot command

type CommandSlot struct{}

func (c *CommandSlot) Base() string { return "slot" }
func (c *CommandSlot) Run(s *discordgo.Session, m *discordgo.Message, split []string, isPrivate bool) {
	//A command call has been detected and we need to check if it's valid
	switch len(split) {
	case 1:
		//-=example
		sendm(m.ChannelID, "Try using the help command if you're not sure how to play.")
	case 2:
		//-=example command
		bet, err := strconv.Atoi(split[1])
		if err != nil {
			sendm(m.ChannelID, "We couldn't convert that into a string.")
		}
		isWinner, pot, msg := slots(bet, m.Author.ID)
		if isPrivate == true { //Check if it's a private message, if so we'll respond differently.
			sendm(m.ChannelID, msg)
		} else {
			if isWinner == false {
				sendm("vegas", fmt.Sprintf("<@%v> just lost a %v bet.", m.Author.ID, pot))
			} else {
				sendm("vegas", fmt.Sprintf("<@%v> just WON %v", m.Author.ID, pot))
			}
		}
	default:
		//More than we accept.
		sendm(m.ChannelID, "Invalid request. Looking for; -=slot X |where X is an int above 5.")
	}
}
func (c *CommandSlot) Help(specific bool) string {
	switch specific {
	case true:
		return "The specific bool is true so we return a much lengthier string. This will be intended to list and explain all subcommands it may have"
	default:
		return "Wasn't called specifically so we return a breif description."
	}

}

// Add all commands to a slice here
var commands = []Command{
	&CommandHelp{},
	&CommandEcho{},
	&CommandSlot{},
}

func HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Printf("%20s %20s %20s > %s |%v|%v|\n", m.Message.ChannelID, time.Now().Format(time.Stamp), m.Author.ID, m.Content, len(m.Content), len(strings.Split(m.Content, " ")))
	if strings.HasPrefix(m.Content, prefix) {
		presplit := strings.TrimPrefix(m.Content, "-=")
		split := strings.Split(presplit, " ")
		chn, err := s.Channel(m.ChannelID)
		if err != nil {
			panic(err)
		}
		for _, cmd := range commands {
			if split[0] == cmd.Base() {
				cmd.Run(s, m.Message, split, chn.IsPrivate)
			}
		}
	}
}
