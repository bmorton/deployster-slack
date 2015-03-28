package main

import (
	"fmt"
	"log"
	"path"
	"strconv"

	"github.com/coreos/go-etcd/etcd"
)

type Plugin struct {
	client       *etcd.Client
	Name         string
	eventsChan   chan string
	watchChan    chan *etcd.Response
	stopChan     chan bool
	handler      Handler
	ReplayUnseen bool
}

type Handler interface {
	Handle(string)
}

func New(name string, client *etcd.Client) *Plugin {
	watchChan := make(chan *etcd.Response, 100)
	eventsChan := make(chan string, 100)
	stopChan := make(chan bool)

	return &Plugin{
		Name:         name,
		eventsChan:   eventsChan,
		watchChan:    watchChan,
		stopChan:     stopChan,
		client:       client,
		ReplayUnseen: false,
	}
}

func (p *Plugin) SetHandler(newHandler Handler) {
	p.handler = newHandler
}

func (p *Plugin) Run() error {
	go p.loop()

	if p.ReplayUnseen {
		p.replayUnseen()
	}

	return p.watch()
}

func (p *Plugin) loop() {
	for {
		select {
		case w := <-p.watchChan:
			p.send(w.Node)
		case e := <-p.eventsChan:
			p.handler.Handle(e)
		}
	}
}

func (p *Plugin) replayUnseen() error {
	log.Println("Replaying unseen events...")
	events, err := p.client.Get("/deployster/events", true, true)
	if err != nil {
		return err
	}

	for _, event := range events.Node.Nodes {
		p.send(event)
	}

	return nil
}

func (p *Plugin) watch() error {
	log.Println("Watching for new events...")
	_, err := p.client.Watch("/deployster/events", 0, true, p.watchChan, nil)
	return err
}

func (p *Plugin) send(node *etcd.Node) {
	timestamp := keyToTimestamp(node.Key)
	event := node.Value

	if timestamp != 0 && timestamp > p.LastSeen() {
		p.eventsChan <- event
		p.UpdateLastSeen(timestamp)
	}
	return
}

func (p *Plugin) LastSeen() uint32 {
	last, err := p.client.Get(fmt.Sprintf("/deployster/plugins/%s/last_seen", p.Name), true, true)
	if err != nil {
		return 0
	}

	lastInt, err := strconv.ParseUint(last.Node.Value, 10, 32)
	if err != nil {
		return 0
	}

	return uint32(lastInt)
}

func (p *Plugin) UpdateLastSeen(timestamp uint32) error {
	_, err := p.client.Set(fmt.Sprintf("/deployster/plugins/%s/last_seen", p.Name), strconv.FormatUint(uint64(timestamp), 10), 0)
	return err
}

func keyToTimestamp(timestamp string) uint32 {
	conv, err := strconv.ParseUint(path.Base(timestamp), 10, 32)
	if err != nil {
		return 0
	}

	return uint32(conv)
}
