package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var paydayTimes = make(map[string]time.Time)
var slotTimes = make(map[string]time.Time)

func logSuggest(id, msg string) {
	if _, err := os.Stat("bank/" + id + "/ideas.txt"); os.IsNotExist(err) {
		// path/to/whatever does not exist
		err := ioutil.WriteFile("bank/"+id+"/ideas.txt", []byte("Start of suggestions;\r\n"), 644) //try to write to file.
		if err != nil {
			panic(err)
		}
		b, err := ioutil.ReadFile("bank/" + id + "/ideas.txt") //Get what's there so we don't overwrite
		if err != nil {
			panic(err)
		}
		s := string(b)
		s += msg
		err1 := ioutil.WriteFile("bank/"+id+"/ideas.txt", []byte(s+"\r\n"), 644) //try to write to file.
		if err1 != nil {
			panic(err1)
		}
	} else {
		b, err := ioutil.ReadFile("bank/" + id + "/ideas.txt") //Get what's there so we don't overwrite
		if err != nil {
			panic(err)
		}
		s := string(b)
		s += msg
		err1 := ioutil.WriteFile("bank/"+id+"/ideas.txt", []byte(s+"\r\n"), 644) ///try to write to file.
		if err1 != nil {
			panic(err1)
		}
	}

}

func makeBank(id string) {
	os.Mkdir("bank/"+id, 644)                                                        //make their bank directory
	err := ioutil.WriteFile("bank/"+id+"/money.txt", []byte(strconv.Itoa(100)), 644) //try to deposit money into their account
	if err != nil {
		panic(err)
	}
}

func hasBank(id string) bool {
	if _, err := os.Stat("bank/" + id); os.IsNotExist(err) {
		// path/to/whatever does not exist
		return false
	} else {
		return true
	}
}

func bankX(id, target string, dolla int) string {
	if dolla > 0 {
		if dolla > getMoney(id) {
			return fmt.Sprintf("Stop being so poor, you don't have that kind of money.")
		} else {
			takeMoney(id, dolla)
			targetid := strings.TrimPrefix(target, "<@")
			targetid = strings.TrimSuffix(targetid, ">")
			putMoney(targetid, dolla)
			return fmt.Sprintf("%d dolla transfered from <@%s>(new balance: %d) to <@%s>(new balance: %d)", dolla, id, getMoney(id), targetid, getMoney(targetid))
		}

	} else {
		return fmt.Sprintf("There's no theivery allowed here, that's for staff only.")
	}
}
func getMoney(id string) int {
	b, err := ioutil.ReadFile("bank/" + id + "/money.txt") //try to see how much money they have
	if err != nil {
		panic(err)
	}
	s := string(b)
	bumhole, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("Something went wrong, couldn't find what I was looking for. Try it again maybe? Error: %s", err)
	}
	return bumhole
}

func putMoney(id string, dolla int) {
	startmoney := getMoney(id)
	startmoney += dolla

	err := ioutil.WriteFile("bank/"+id+"/money.txt", []byte(strconv.Itoa(startmoney)), 644) //try to deposit money into their account
	if err != nil {
		panic(err)
	}
}
func takeMoney(id string, dolla int) {
	startmoney := getMoney(id)
	startmoney -= dolla

	err := ioutil.WriteFile("bank/"+id+"/money.txt", []byte(strconv.Itoa(startmoney)), 644) //try to take money from their account
	if err != nil {
		panic(err)
	}
}
func givepayday(id string) string {

	lastTime, ok := paydayTimes[id]
	if !ok || time.Since(lastTime).Minutes() >= 5 {
		startdolla := getMoney(id)
		putMoney(id, 250)
		paydayTimes[id] = time.Now()
		return fmt.Sprintf("You had: %d Dolla bill. Now you have: %d", startdolla, getMoney(id))
	} else {
		fullLeft := 300 - time.Since(lastTime).Seconds()

		seconds := int(fullLeft) % 60
		minutes := math.Floor(fullLeft / 60)

		return fmt.Sprintf("You need to wait %d Mins and %d Seconds", int(minutes), seconds)
	}
}
func sendm(chn, msg string) {
	chnp, _ := ds.Channel(chn) //We might need this error, but I can't think why atm. Just be warned.
	switch chn {
	case "vegas":
		ds.ChannelMessageSend("160762694549372929", msg)
	case "160762694549372929":
		ds.ChannelMessageSend("160762694549372929", msg)
	default:
		if chnp.IsPrivate == true {
			ds.ChannelMessageSend(chn, msg)
		} else {
			// Need to figure a better way of handling this.
			// Basically, they're calling the bot from a random
			// channel but we don't have access to their ID or
			// anything in this function just yet and I don't
			// want to leap to anything right now otherwise
			// it might turn ugly. Very very very very ugly.
		}

	}

}

func slots(bet int, id string) (bool, int, string) {

	choices := []string{":cherries:", ":8ball:", ":keycap_ten:", ":four_leaf_clover:", ":fire:", ":octopus:", ":zero:", ":tada:", ":heart:", ":snowflake:"}
	slots := [][]string{
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
	}
	pot := bet

	slots[0][0] = choices[rand.Intn(len(choices))]
	slots[0][1] = choices[rand.Intn(len(choices))]
	slots[0][2] = choices[rand.Intn(len(choices))]

	slots[1][0] = choices[rand.Intn(len(choices))]
	slots[1][1] = choices[rand.Intn(len(choices))]
	slots[1][2] = choices[rand.Intn(len(choices))]

	slots[2][0] = choices[rand.Intn(len(choices))]
	slots[2][1] = choices[rand.Intn(len(choices))]
	slots[2][2] = choices[rand.Intn(len(choices))]

	center := slots[1]

	msg := ""
	won := false
	for i := 0; i < len(slots); i++ {
		if i == 1 {
			msg += "> " + strings.Join(slots[i], " ") + "\n"
		} else {
			msg += "   " + strings.Join(slots[i], " ") + "\n"
		}
	}
	result := ""
	//sm := getMoney(id)
	if center[0] == center[1] && center[1] == center[2] && center[2] == center[0] {
		pot *= 6
		putMoney(id, pot)
		result = fmt.Sprintf("3 of a kind! *6! You get: %d. <@%s>, you have %d left.", pot, id, getMoney(id))
	} else if center[0] == ":four_leaf_clover:" && center[1] == center[0] && center[2] == center[1] {
		pot *= 1200
		putMoney(id, pot)
		result = fmt.Sprintf("Twelve leaf clover! *1200! You get: %d. <@%s>, you have %d left.", pot, id, getMoney(id))
	} else if center[0] == ":keycap_ten:" && center[1] == center[2] && center[2] == ":zero:" {
		pot += 1000
		putMoney(id, pot)
		result = fmt.Sprintf("1000 get! +1000! You get: %d. <@%s>, you have %d left.", pot, id, getMoney(id))
	} else if center[0] == ":8ball:" && center[1] == ":8ball:" || center[1] == ":8ball:" && center[2] == ":8ball:" {
		pot *= 4
		putMoney(id, pot)
		result = fmt.Sprintf("Nice balls! *4! You get: %d. <@%s>, you have %d left.", pot, id, getMoney(id))
	} else if center[0] == ":fire:" && center[1] == ":fire:" || center[1] == ":fire:" && center[2] == ":fire:" {
		pot += 500
		putMoney(id, pot)
		result = fmt.Sprintf("Two copies of my mixtape! +500! You get: %d. <@%s>, you have %d left.", pot, id, getMoney(id))
	} else if center[0] == ":tada:" && center[1] == ":tada:" || center[1] == ":tada:" && center[2] == ":tada:" {
		pot += 100
		putMoney(id, pot)
		result = fmt.Sprintf("TADAAAA! +100! You get: %d. <@%s>, you have %d left.", pot, id, getMoney(id))
	} else if center[0] == ":squirrel:" && center[1] == ":8ball:" || center[1] == ":squirrel:" && center[2] == ":8ball:" {
		pot *= 3
		putMoney(id, pot)
		result = fmt.Sprintf("Squirrelly goodness! *3! You get: %d. <@%s>, you have %d left.", pot, id, getMoney(id))
	} else if center[0] == center[1] || center[1] == center[2] || center[2] == center[0] {
		pot *= 2
		putMoney(id, pot)
		result = fmt.Sprintf("Two of a kind! *2! You get: %d. <@%s>, you have %d left.", pot, id, getMoney(id))
	} else {
		takeMoney(id, pot)
		//addToLotteryPot(id, bet)
		result = fmt.Sprintf("Say goodbye to your %d dolla bet. <@%s>! hahaaha. You have %d dollas left.\n", pot, id, getMoney(id))
	}
	if pot > bet { //If the pot is bigger than your bet, then you must have won something.
		won = true
	}
	return won, pot, fmt.Sprintf("\n%s\n%s", msg, result)
}

//func newslots(bet int, id string) (bool, int, string) {
//	pot := bet
//	msg := ""
//	won := false
//	lreel := []string{":cherries:", ":four_leaf_clover:", ":octopus:", ":fire:", ":8ball:", ":keycap_ten:", ":zero:", ":snowflake:", ":tada:", ":heart:"}
//	mreel := []string{":fire:", ":8ball:", ":keycap_ten:", ":zero:", ":snowflake:", ":tada:", ":heart:", ":cherries:", ":four_leaf_clover:", ":octopus:"}
//	rreel := []string{":snowflake:", ":tada:", ":fire:", ":8ball:", ":keycap_ten:", ":zero:", ":cherries:", ":octopus:", ":heart:", ":four_leaf_clover:"}
//
//	midnum := []int{rand.Intn(len(lreel)), rand.Intn(len(mreel)), rand.Intn(len(rreel))}
//
//	toprow := []string{lreel[midnum[0]-1%len(lreel)], mreel[midnum[1]-1%len(mreel)], rreel[midnum[2]-1%len(rreel)]} //This is going to be out of range :c
//	midrow := []string{lreel[midnum[0]], mreel[midnum[1]], rreel[midnum[2]]}
//	botrow =: []string{lreel[midnum[0]+1%len(lreel)], mreel[midnum[1]+1%len(mreel)], rreel[midnum[2]+1%len(rreel)]} // This too. How do I wrap around? or is it possible?
//
//	msg += toprow[0] + " " + toprow[2] + " " + toprow[2] + "\n"
//	msg += midrow[0] + " " + midrow[2] + " " + midrow[2] + "\n"
//	msg += botrow[0] + " " + botrow[2] + " " + botrow[2] + "\n"
//	msg += ""
//
//	if pot > bet { //If the pot is bigger than your bet, then you must have won something.
//		won = true
//	}
//	return won, pot, fmt.Sprintf("\n%s\n Hopefully this worked.", msg)
//}

func top10bad() string {
	one, two, three, four, five, six, seven, eight, nine, ten := "", "", "", "", "", "", "", "", "", ""
	onev, twov, threev, fourv, fivev, sixv, sevenv, eightv, ninev, tenv := 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
	curr := ""
	curv := 0

	files, _ := ioutil.ReadDir("bank/")
	for _, f := range files {
		curr = f.Name()
		curv = getMoney(f.Name())
		if curv >= onev {
			ten, tenv = nine, ninev
			nine, ninev = eight, eightv
			eight, eightv = seven, sevenv
			seven, sevenv = six, sixv
			six, sixv = five, fivev
			five, fivev = four, fourv
			four, fourv = three, threev
			three, threev = two, twov
			two, twov = one, onev
			one, onev = curr, curv
		} else if curv >= twov && curv < onev {
			ten, tenv = nine, ninev
			nine, ninev = eight, eightv
			eight, eightv = seven, sevenv
			seven, sevenv = six, sixv
			six, sixv = five, fivev
			five, fivev = four, fourv
			four, fourv = three, threev
			three, threev = two, twov
			two, twov = curr, curv
		} else if curv >= threev && curv < twov {
			ten, tenv = nine, ninev
			nine, ninev = eight, eightv
			eight, eightv = seven, sevenv
			seven, sevenv = six, sixv
			six, sixv = five, fivev
			five, fivev = four, fourv
			four, fourv = three, threev
			three, threev = curr, curv
		} else if curv >= fourv && curv < threev {
			ten, tenv = nine, ninev
			nine, ninev = eight, eightv
			eight, eightv = seven, sevenv
			seven, sevenv = six, sixv
			six, sixv = five, fivev
			five, fivev = four, fourv
			four, fourv = curr, curv
		} else if curv >= fivev && curv < fourv {
			ten, tenv = nine, ninev
			nine, ninev = eight, eightv
			eight, eightv = seven, sevenv
			seven, sevenv = six, sixv
			six, sixv = five, fivev
			five, fivev = curr, curv
		} else if curv >= sixv && curv < fivev {
			ten, tenv = nine, ninev
			nine, ninev = eight, eightv
			eight, eightv = seven, sevenv
			seven, sevenv = six, sixv
			six, sixv = curr, curv
		} else if curv >= sevenv && curv < sixv {
			ten, tenv = nine, ninev
			nine, ninev = eight, eightv
			eight, eightv = seven, sevenv
			seven, sevenv = curr, curv
		} else if curv >= eightv && curv < sevenv {
			ten, tenv = nine, ninev
			nine, ninev = eight, eightv
			eight, eightv = curr, curv
		} else if curv >= ninev && curv < eightv {
			ten, tenv = nine, ninev
			nine, ninev = curr, curv
		} else if curv >= tenv && curv < ninev {
			ten, tenv = curr, curv
		}
	}
	return fmt.Sprintf("```1:%v-%v\n2:%v-%v\n3:%v-%v\n4:%v-%v\n5:%v-%v\n6:%v-%v\n7:%v-%v\n8:%v-%v\n9:%v-%v\n10:%v-%v\n```", one, onev, two, twov, three, threev, four, fourv, five, fivev, six, sixv, seven, sevenv, eight, eightv, nine, ninev, ten, tenv)
}

type Player struct {
	Id    string
	Money int
}

func top10(s *discordgo.Session, m *discordgo.Message) string {
	top := make([]*Player, 10)
	files, _ := ioutil.ReadDir("bank/")
	for _, f := range files {
		name := f.Name()
		money := getMoney(f.Name())
		// Iterate over every file
		for i := 0; i < len(top); i++ {
			// We insert it here then
			if top[i] == nil || money > top[i].Money {
				// Move everything back once
				// [100, 99, 55] turns into
				// [x, 100, 99] where x is what we inserted (which has to be higher than 100 in this example)
				for j := len(top) - 1; j > i; j-- {
					top[j] = top[j-1]
				}
				// Assign it
				top[i] = &Player{
					Id:    name,
					Money: money,
				}
				break
			}
		}
	}
	str := "```"
	str += top10Str(top, s, m)
	str += "```"
	return str
	//return fmt.Sprintf("1; %v | %v\n2; %v | %v\n3;%v | %v", top[0].Id, top[0].Money, top[1].Id, top[1].Money, top[2].Id, top[2].Money)
}
func top10Str(top10 []*Player, s *discordgo.Session, m *discordgo.Message) string {
	str := ""
	for k, v := range top10 {
		if v == nil {
			break
		} //v.Id
		channel, _ := s.Channel(m.ChannelID)

		name, err := s.State.Member(channel.GuildID, v.Id)
		namev := "OFFLINE:" + v.Id
		if err == nil {
			namev = name.User.Username
		}
		//namev := name.User.Username

		str += fmt.Sprintf("%2d. %27s: %d\n", k+1, namev, v.Money)
	}
	return str
}

func setmoney(monies int, id string) {
	err := ioutil.WriteFile("bank/"+id+"/money.txt", []byte(strconv.Itoa(monies)), 644) //try to deposit money into their account
	if err != nil {
		panic(err)
	}

}
