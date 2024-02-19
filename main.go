package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	BotToken = flag.String("token", "", "discord bot token")
	//ServerId = flag.String("guild", "", "discord server id")
)

func main() {
	dc, err := discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	dc.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	dc.AddHandler(messageCreate)

	dc.Identify.Intents = discordgo.IntentsAll

	err = dc.Open()
	if err != nil {
		log.Fatalf("Cannot login: %v", err)
	}
	defer dc.Close()

	fmt.Println("bot is now running - press ctrl + c to stop")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!") {
		var command []string
		command = strings.Split(m.Content, " ")

		if command[0] == "!teams" {
			channel, err := s.State.Channel(m.ChannelID)
			if err != nil {
				log.Fatalf("Could not find channel: %v", err)
				return
			}
			guild, err := s.State.Guild(channel.GuildID)
			if err != nil {
				log.Fatalf("Could not find guild: %v", err)
				return
			}

			var users []string

			for _, vs := range guild.VoiceStates {
				if vs.UserID == m.Author.ID {
					for _, vss := range guild.VoiceStates {
						if vss.ChannelID == vs.ChannelID {
							users = append(users, vss.UserID)
						}
					}
				}
			}
			mid := len(users) / 2
			team1 := users[:mid]
			team2 := users[mid:]

			res1 := strings.Join(team1, "\n")
			res2 := strings.Join(team2, "\n")

			result := fmt.Sprintf("Wylosowałem drużyny!\n Drużyna 1: %s\n Drużyna 2: %s", res1, res2)

			_, _ = s.ChannelMessageSend(m.ChannelID, result)
		}
	}
}
