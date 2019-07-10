package pubsub

import (
	"github.com/eddieowens/axon"
	"github.com/eddieowens/ranvier/server/app/model"
	"strings"
	"sync"
)

const Key = "PubSub"

type PubSub interface {
	Publish(topic string, config *model.ConfigEvent)
	Subscribe(topic string) chan model.ConfigEvent
}

type pubSubImpl struct {
	topics map[string][]chan model.ConfigEvent
	lock   sync.RWMutex
}

func (p *pubSubImpl) Publish(topic string, config *model.ConfigEvent) {
	topic = strings.ToLower(topic)
	p.lock.RLock()
	defer p.lock.RUnlock()
	topicChannels := p.topics[topic]
	for _, c := range topicChannels {
		c <- *config
	}

	allChannels := p.topics["*"]
	for _, c := range allChannels {
		c <- *config
	}
}

func (p *pubSubImpl) Subscribe(topic string) chan model.ConfigEvent {
	c := make(chan model.ConfigEvent)
	topic = strings.ToLower(topic)
	p.lock.Lock()
	defer p.lock.Unlock()
	topicChannels := p.topics[topic]
	p.topics[topic] = p.addChannel(topicChannels, c)
	return c
}

func (p pubSubImpl) addChannel(channels []chan model.ConfigEvent, c chan model.ConfigEvent) []chan model.ConfigEvent {
	if channels == nil {
		channels = make([]chan model.ConfigEvent, 1)
		channels[0] = c
	} else {
		channels = append(channels, c)
	}
	return channels
}

func pubSubFactory(_ axon.Injector, _ axon.Args) axon.Instance {
	p := &pubSubImpl{}

	p.topics = make(map[string][]chan model.ConfigEvent)
	p.lock = sync.RWMutex{}

	return axon.StructPtr(p)
}
