package main

import (
	"fmt"
	"log"
	"os"

	"github.com/coreos/go-etcd/etcd"
	"github.com/nlopes/slack"
)

func main() {
	client := etcd.NewClient([]string{os.Getenv("ETCDCTL_PEERS")})
	plugin := New("slack", client)

	api := slack.New(os.Getenv("SLACK_TOKEN"))
	plugin.SetDeployHandler(&DeployEventHandler{client: api, channel: os.Getenv("SLACK_CHANNEL")})
	plugin.ReplayUnseen = true

	log.Println(plugin.Run())
}

type DeployEventHandler struct {
	client  *slack.Slack
	channel string
}

func (e *DeployEventHandler) Handle(event *DeployEvent) {
	params := slack.PostMessageParameters{
		Username: "Deployster",
	}
	message := fmt.Sprintf("Deploy started for %d instance(s) of %s (version: %s) at %s", event.InstanceCount, event.ServiceName, event.Version, event.Timestamp)
	channel, timestamp, err := e.client.PostMessage(e.channel, message, params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	log.Printf("Message successfully sent to channel %s at %s", channel, timestamp)
}
