package main

import (
	"fmt"
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

	err := ioutil.WriteFile("bank/"+id+"/money.txt", []byte(strconv.Itoa(startmoney)), 644) //try to deposit money into their account
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
