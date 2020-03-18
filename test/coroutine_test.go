package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	stream "github.com/mcmakler/stream-chat-go"
	"gopkg.in/go-playground/assert.v1"
)

const (
	APIKey    = os.Getenv("STREAM_CHAT_API_KEY")
	APISecret = os.Getenv("STREAM_CHAT_API_SECRET")
)

var (
	testClient *stream.Client
)

func init() {
	var err error
	// replace with correct key and secret
	testClient, err = stream.NewClient(APIKey, []byte(APISecret))
	if err != nil {
		fmt.Println("init getstream error:", err)
	}
}

func updateUsers(users ...*stream.User) error {
	_, err := testClient.UpdateUsers(users...)
	if err != nil {
		fmt.Println("update users error:", err)
	}
	return err
}

func createChannel(channelID string, users ...*stream.User) error {
	members := []string{}
	for _, user := range users {
		members = append(members, user.ID)
	}

	fmt.Println("creating channel ", channelID)
	_, err := testClient.CreateChannel("team", channelID, "admin", map[string]interface{}{
		"members": members,
	})
	if err != nil {
		fmt.Println("create channel error:", err)
		return err
	}

	return err
}

func sendMessage(channelID, text string) (*stream.Message, error) {
	channels, err := testClient.QueryChannels(&stream.QueryOption{
		Filter: map[string]interface{}{
			"id":   channelID,
			"type": "team",
		},
	})
	if err != nil {
		return nil, err
	}
	if len(channels) == 1 {
		c := channels[0]

		req := &stream.Message{
			Text: text,
		}
		return c.SendMessage(req, "admin")
	}
	return nil, fmt.Errorf("channel does not exist")
}

func TestUpateUser(t *testing.T) {
	err := updateUsers(&stream.User{
		ID: "test",
	}, &stream.User{
		ID: "admin",
	})
	assert.Equal(t, nil, err)
}

func TestCreateChaneel(t *testing.T) {
	err := createChannel("new_channel", &stream.User{
		ID: "test",
	}, &stream.User{
		ID: "admin",
	})
	assert.Equal(t, nil, err)
}

func TestSendMessages(t *testing.T) {
	go send(t)
	go send(t)
	go send(t)

	time.Sleep(time.Second * 3)
}

func send(t *testing.T) {
	_, err := sendMessage("new_channel", "some text")
	assert.Equal(t, nil, err)
}
