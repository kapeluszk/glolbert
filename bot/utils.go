package bot

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
)

func ShuffleStrings(inputSlice []string) []string {
	rand.Shuffle(len(inputSlice), func(i, j int) {
		inputSlice[i], inputSlice[j] = inputSlice[j], inputSlice[i]
	})
	return inputSlice
}

func CurrentVcMembers(s *discordgo.Session, m *discordgo.MessageCreate) []string {

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Fatalf("Could not find channel: %v", err)

	}
	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		log.Fatalf("Could not find guild: %v", err)

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
	return users
}
