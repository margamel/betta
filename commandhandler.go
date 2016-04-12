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
	switch len(split) {
	case 2: //-=help CMD?
		for _, cmd := range commands {
			if split[1] == cmd.Base() {
				helpmsg := cmd.Help(true)
				sendm(m.ChannelID, helpmsg)
			}
		}
	default:
		helpmsg := "We have these commands; \n"
		for _, cmd := range commands {
			helpmsg += fmt.Sprintf("%v | %10v\n", cmd.Base(), cmd.Help(false))

		}
		sendm(m.ChannelID, helpmsg)
	}
}
func (c *CommandHelp) Help(specific bool) string {
	return "You stupid i smart"
}

// Another command
type CommandEcho struct{}

func (c *CommandEcho) Base() string { return "echo" }
func (c *CommandEcho) Run(s *discordgo.Session, m *discordgo.Message, split []string, isPrivate bool) {
	// Echo what was sent to us
	msg := strings.TrimPrefix(m.Content, "-=echo ")
	switch m.Author.ID {
	case "160718675689603072":
		sendm(m.ChannelID, "C'mon now, I serve you, not myself. I'm not worthy of using this command.")
	default:
		sendm(m.ChannelID, msg)
	}

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
			sendm(m.ChannelID, "We couldn't convert that into an int.")
		}
		if bet > 5 && bet <= getMoney(m.Author.ID) {
			isWinner, pot, msg := slots(bet, m.Author.ID)
			if isPrivate == true { //Check if it's a private message, if so we'll respond differently.
				sendm(m.ChannelID, msg)
				if isWinner == true {
					sendm("vegas", fmt.Sprintf("<@%v> just WON a %v!", m.Author.ID, pot))
				}
			} else {
				if isWinner == false {
					sendm("vegas", fmt.Sprintf("<@%v> just lost a %v bet.", m.Author.ID, pot))
				} else {
					sendm("vegas", fmt.Sprintf("<@%v> just WON %v!", m.Author.ID, pot))
				}
			}
		} else {
			sendm(m.ChannelID, "Your bet needs to be between 10 and your wealth.")
		}

	default:
		//More than we accept.
		sendm(m.ChannelID, "Invalid request. Looking for; -=slot X |where X is an int above 10.")
	}
}
func (c *CommandSlot) Help(specific bool) string {
	switch specific {
	case true:
		return "Risk your wealth on the slots, -=slot 500 if you're feeling brave enough."
	default:
		return "Play the slots."
	}

}

//bank command

type CommandBank struct{}

func (c *CommandBank) Base() string { return "bank" }
func (c *CommandBank) Run(s *discordgo.Session, m *discordgo.Message, split []string, isPrivate bool) {
	//A command call has been detected and we need to check if it's valid
	switch len(split) {
	case 1:
		//-=bank
		sendm(m.ChannelID, "Try using the help command if you're not sure how to .")
	case 2:
		//-=bank register
		switch split[1] {
		case "balance":
			sendm(m.ChannelID, fmt.Sprintf("You current account standing is: %v", getMoney(m.Author.ID)))

		default:
			sendm(m.ChannelID, "Hmm, we didn't quite understand that one. Try something else?")
		}
	case 4:
		//-=bank xfer target money
		switch split[1] {
		case "xfer":
			//lets check if the 3rd argument can be casted to an int.
			//Yes- they either put in random numbers, or they put in the id number. We need to check for both and handle it.
			//No - they either put a random word, or they @mentioned them. We need to check for both and handle it.

			switch m.Author.ID { // Fucking ashe making me protect against stupid shit all the time.
			case "160718675689603072": //But yeah, this means they can still play with the bots money, but not touch it.
				sendm(m.ChannelID, "You can make me money, but you can't make me give it to you ;)")
			default:
				_, err := strconv.Atoi(split[2])
				if err != nil { //We couldn't cast to int so we're going to see if they @mentioned them.
					target := strings.TrimPrefix(split[2], "<@")
					target = strings.TrimSuffix(target, ">")
					dolla, err := strconv.Atoi(split[3])
					if err != nil {
						//no fucking clue mate.
						panic(err)
					}

					_, err1 := strconv.Atoi(target)
					if err1 != nil { //we still can't cast to int. This means it was just a random word all along.
						sendm(m.ChannelID, "We were expecting a target for the transfer. Either their ID or @mention.")
					} else {
						//We can cast to int now, where before we couldn't. This means that it was a format mention all along.
						sendm(m.ChannelID, bankX(m.Author.ID, target, dolla))
					}

				} else {
					//we could cast to int so it was either random numbers, or the ID
					switch len(split[3]) {
					case 18: //95% sure this is going to be the ID so we can cast to string and call the transfer.
						target := strings.TrimPrefix(split[2], "<@")
						target = strings.TrimSuffix(target, ">")
						dolla, err := strconv.Atoi(split[3])
						if err != nil {
							//no fucking clue mate.
							panic(err)
						}
						sendm(m.ChannelID, bankX(m.Author.ID, target, dolla))
					default:
					}
				}
			}

		default:
		}
	default:
		//More than we accept.
		sendm(m.ChannelID, "Invalid request. Feel free to ask us for help should you need it.")
	}
}
func (c *CommandBank) Help(specific bool) string {
	switch specific {
	case true:
		return "Welcome to the bank. We have a few subcommands for you to use. xfer and balance come to mind. \n '-=bank xfer @Runi 100' gives my creator a tip of 100 dosh."
	default:
		return "Use the bank to hold, check, and transfer your money."
	}
}

// Another command
type CommandSuggest struct{}

func (c *CommandSuggest) Base() string { return "suggest" }
func (c *CommandSuggest) Run(s *discordgo.Session, m *discordgo.Message, split []string, isPrivate bool) {
	// Echo what was sent to us
	msg := strings.TrimPrefix(m.Content, "-=suggest ")
	switch m.Author.ID {
	case "160718675689603072":
		sendm(m.ChannelID, "C'mon now, I serve you, not myself. I'm not worthy of using this command.")
	default:
		logSuggest(m.Author.ID, msg)
	}

}

func (c *CommandSuggest) Help(specific bool) string {
	switch specific {
	case true:
		return "Suggest games, commands, how you think the bot should work, what it should say. Anything you want that is bot related, suggest it."
	default:
		return "Suggest features that you would like to see implemented in the future."
	}
}

type CommandMexicanwave struct{}

func (c *CommandMexicanwave) Base() string { return "mexicanwave" }
func (c *CommandMexicanwave) Run(s *discordgo.Session, m *discordgo.Message, split []string, isPrivate bool) {
	switch m.Author.ID {
	case "105661408393302016":
		//¯\_(ツ)_/¯¯\_(ツ)_/¯¯\_(ツ)_/¯¯\_(ツ)_/¯¯\_(ツ)_/¯¯\_(ツ)_/¯¯\_(ツ)_/¯
		frame := 1
		frames := 0
		wave := "¯\\(ツ)_/¯"
		rest := "-__(ツ)__-"
		outwave := rest + rest + rest + rest + wave
		frame1 := wave + rest + rest + rest + rest
		frame2 := rest + wave + rest + rest + rest
		frame3 := rest + rest + wave + rest + rest
		frame4 := rest + rest + rest + wave + rest
		frame5 := rest + rest + rest + rest + wave

		msg, _ := s.ChannelMessageSend("160762694549372929", "```"+outwave+"```")
		ticker := time.NewTicker(time.Second * 2)
		for {
			select {
			case <-ticker.C:
				if frames < 5 {
					frame++
					frames++
					if frame > 4 {
						frame = 1
					}
					switch frame {
					case 1:
						outwave = frame1
					case 2:
						outwave = frame2
					case 3:
						outwave = frame3
					case 4:
						outwave = frame4
					case 5:
						outwave = frame5
					default:
						outwave = "Error, lol."
					}
					msg, _ = s.ChannelMessageEdit(msg.ChannelID, msg.ID, "```"+outwave+"```")
				} else {
					return
				}
			}
		}
	default:
		s.ChannelMessageSend("160762694549372929", "```NO!```")

	}

}

func (c *CommandMexicanwave) Help(specific bool) string {
	switch specific {
	case true:
		return "I made this on a whim. It's not too good though since it rate limits faaaaaast."
	default:
		return "Do tha mexican waaaaAAAAAAaaaaVVVVVEeeeeee."
	}
}

type CommandPayday struct{}

func (c *CommandPayday) Base() string { return "payday" }
func (c *CommandPayday) Run(s *discordgo.Session, m *discordgo.Message, split []string, isPrivate bool) {
	switch isPrivate {
	case true:
		sendm(m.ChannelID, givepayday(m.Author.ID))
	default:
		givepayday(m.Author.ID)
	}

}

func (c *CommandPayday) Help(specific bool) string {
	switch specific {
	case true:
		return "Once every 5 mins you can get some free money. Try to gamble responsibly... LOL!"
	default:
		return "Free money every once in a while."
	}
}

// Add all commands to a slice here
var commands = []Command{
	&CommandHelp{},
	&CommandEcho{},
	&CommandSlot{},
	&CommandBank{},
	&CommandSuggest{},
	&CommandMexicanwave{},
	&CommandPayday{},
}

func HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	chn, err := s.Channel(m.ChannelID)
	if err != nil {
		panic(err)
	}
	if m.Author.ID == "" { //sometimes the ID will be nil and then it just dies an inglorious death. Not sure why or what to really do about it.

	}
	fmt.Printf("%20s %20s %20s > %s |%v|%v|\n", chn.Name, time.Now().Format(time.Stamp), m.Author.ID, m.Content, len(m.Content), len(strings.Split(m.Content, " ")))
	if strings.HasPrefix(m.Content, prefix) {
		if hasBank(m.Author.ID) == false { //Just go ahead and create accounts for everyone. I mean, why not. Saves the heartache of bankchecking later on.
			makeBank(m.Author.ID)
		}

		presplit := strings.TrimPrefix(m.Content, "-=") //clear the prefix before we split the string
		split := strings.Split(presplit, " ")           //split the sting by the spaces.

		for _, cmd := range commands { //Loop through the commands and run whichever command matches.
			if split[0] == cmd.Base() {
				cmd.Run(s, m.Message, split, chn.IsPrivate)
			}
		}
	}
}
