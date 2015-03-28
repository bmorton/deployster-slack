package main

import (
	"fmt"

	"github.com/coreos/go-etcd/etcd"
)

func main() {
	client := etcd.NewClient([]string{"http://127.0.0.1:2379"})
	plugin := New("slack", client)

	plugin.SetHandler(&EventHandler{})
	plugin.ReplayUnseen = true

	plugin.Run()
}

type EventHandler struct{}

func (e *EventHandler) Handle(event string) {
	fmt.Println(event)
}
