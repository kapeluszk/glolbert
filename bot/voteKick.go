package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"sync"
	"time"
)

func AddVote(kickedID string, VoteCount map[string]int, mutex *sync.Mutex, voteThreshold int, ChanMap map[string]chan bool, s *discordgo.Session, channelID string) {
	mutex.Lock()
	defer mutex.Unlock()

	count := VoteCount[kickedID]
	println(count)

	VoteCount[kickedID] += 1

	if VoteCount[kickedID] < voteThreshold {
		response := fmt.Sprintf("na wyrzucenie użytkownika oddano obecnie %d na %d wymaganych głosów", VoteCount[kickedID], voteThreshold)
		_, _ = s.ChannelMessageSend(channelID, response)
	}

	if VoteCount[kickedID] >= voteThreshold {
		close(ChanMap[kickedID])
	}
}

func VotingHandler(done chan bool, mutex *sync.Mutex, kickedID string, VoteCount map[string]int, s *discordgo.Session, guild *discordgo.Guild, ChanMap map[string]chan bool, channelID string) {
	select {
	case <-done:

		err := s.GuildMemberMove(guild.ID, kickedID, nil)
		if err != nil {
			log.Fatalf("failed to kick: %v", err)
		}

		response := fmt.Sprintf("Głosowanie udane, wyrzucam użytkownika")
		_, _ = s.ChannelMessageSend(channelID, response)

		mutex.Lock()
		defer mutex.Unlock()

		delete(VoteCount, kickedID)
		delete(ChanMap, kickedID)

	case <-time.After(180 * time.Second):

		response := fmt.Sprintf("Nie uzyskano wymaganej liczby głosów w 180s - zamykam głosowanie")
		_, _ = s.ChannelMessageSend(channelID, response)

		mutex.Lock()
		defer mutex.Unlock()

		delete(VoteCount, kickedID)
		delete(ChanMap, kickedID)
	}
}
