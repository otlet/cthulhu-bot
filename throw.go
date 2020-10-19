package main

import (
	"fmt"
	"math/rand"
	"strings"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func throw(s *discordgo.Session, m *discordgo.MessageCreate) {
	command := strings.Split(m.Content, " ")
	if len(command) < 2 {
		s.ChannelMessageSend(m.ChannelID, "A podasz mi kostkę? Do wyboru: K4 K6 K8 K10 K12 K20 K100")
		return
	}

	howManyThrow := 1
	if len(command) >= 3 {
		i, err := strconv.Atoi(command[2])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "A pałą przez łeb dawno nie dostałeś? ILE RAZY MAM RZUCIĆ KOSTKĄ TO CYFRA, a nie...")
		return
		}
		howManyThrow = i
	}

	for i := 0; i < howManyThrow; i++ {
		randomNumber := 0
		fmt.Println(command)
		switch strings.ToLower(command[1]) {
		case "k4":
			randomNumber = rand.Intn(4) + 1
		case "k6":
			randomNumber = rand.Intn(6) + 1
		case "k8":
			randomNumber = rand.Intn(8) + 1
		case "k10":
			randomNumber = rand.Intn(10) + 1
		case "k12":
			randomNumber = rand.Intn(12) + 1
		case "k20":
			randomNumber = rand.Intn(20) + 1
		case "k100":
			randomNumber = (rand.Intn(10) + 1) * 10
		default:
			s.ChannelMessageSend(m.ChannelID, "Nie znam takiej kostki...")
			return
		}
	
		embed := &discordgo.MessageEmbed{
			Title: fmt.Sprintf("%v rzucił kostką %v i wypadło: %v", m.Author.Username, strings.ToUpper(command[1]), randomNumber),
			//Description: fmt.Sprintf("Wynik %v", randomNumber),
		}
	
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}
}
