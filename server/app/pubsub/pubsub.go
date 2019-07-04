package pubsub

import (
	"github.com/eddieowens/axon"
	"github.com/eddieowens/ranvier/server/app/model"
	"sync"
)

const PubSubKey = "PubSub"

type PubSub interface {
	Publish(topic string, config *model.Config)
	Subscribe(topic string) chan model.Config
}

type pubSubImpl struct {
	topics map[string][]chan model.Config
	lock   sync.RWMutex
}

func (p *pubSubImpl) Publish(topic string, config *model.Config) {
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

func (p *pubSubImpl) Subscribe(topic string) chan model.Config {
	c := make(chan model.Config)
	p.lock.Lock()
	defer p.lock.Unlock()
	topicChannels := p.topics[topic]
	p.topics[topic] = p.addChannel(topicChannels, c)
	return c
}

func (p pubSubImpl) addChannel(channels []chan model.Config, c chan model.Config) []chan model.Config {
	if channels == nil {
		channels = make([]chan model.Config, 1)
		channels[0] = c
	} else {
		channels = append(channels, c)
	}
	return channels
}

func pubSubFactory(_ axon.Injector, _ axon.Args) axon.Instance {
	p := &pubSubImpl{}

	p.topics = make(map[string][]chan model.Config)
	p.lock = sync.RWMutex{}

	return axon.StructPtr(p)
}
