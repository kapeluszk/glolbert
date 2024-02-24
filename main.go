package main

import (
	b "dcBot/bot"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

var VoteCount = make(map[string]int)
var ChanMap = make(map[string]chan bool)

var (
	BotToken = flag.String("token", "", "discord bot token")
	//ServerId = flag.String("guild", "", "discord server id")
)

func main() {
	flag.Parse()
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
			if len(users) < 2 {
				_, _ = s.ChannelMessageSend(m.ChannelID, "zbyt mało użytkowników żeby stworzyć 2 drużyny!")
				return
			}
			var nicknames []string
			members, err := s.GuildMembers(guild.ID, "", 1000)
			if err != nil {
				log.Fatal(err)
			}
			for _, member := range members {
				fmt.Println(member.User.ID)
			}
			for _, uid := range users {
				for _, member := range members {
					if uid == member.User.ID {
						nickname := member.Nick
						if nickname == "" || nickname == "0" {
							nickname = member.User.Username
						}
						nicknames = append(nicknames, nickname)
					}
				}
			}
			nicknames = b.ShuffleStrings(nicknames)
			mid := len(nicknames) / 2
			team1 := nicknames[:mid]
			team2 := nicknames[mid:]

			res1 := strings.Join(team1, "\n")
			res2 := strings.Join(team2, "\n")

			result := fmt.Sprintf("Wylosowałem drużyny!\n Drużyna 1: %s\n Drużyna 2: %s\n", res1, res2)

			_, _ = s.ChannelMessageSend(m.ChannelID, result)
		} else if command[0] == "!votekick" {
			kicked := command[1]
			fmt.Println(kicked)
			var kickedID = ""

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

			members, err := s.GuildMembers(guild.ID, "", 1000)

			if err != nil {
				log.Fatal(err)
			}
			for _, member := range members {
				if kicked == member.Nick {
					fmt.Printf("match!: %s", member.Nick)
					kickedID = member.User.ID
				} else if kicked == member.User.Username {
					kickedID = member.User.ID
				}
			}
			fmt.Println(kickedID)

			var mutex = sync.Mutex{}

			users := b.CurrentVcMembers(s, m)

			voteThreshold := len(users) / 2

			_, ok := ChanMap[kickedID]
			if !ok {
				ChanMap[kickedID] = make(chan bool)
				go b.VotingHandler(ChanMap[kickedID], &mutex, kickedID, VoteCount, s, guild, ChanMap, m.ChannelID)
			}
			b.AddVote(kickedID, VoteCount, &mutex, voteThreshold, ChanMap, s, m.ChannelID)
			fmt.Println(ChanMap)
			fmt.Println(VoteCount[kickedID])
		} else if command[0] == "!test" {

		}
	}
}
