package main

import "sync"

var (
	stickySessioninstance *StickySessionStrategy
)

//StickySessionStrategy ...
type StickySessionStrategy struct {
	Sessions map[string]string

	roundRobinStrategy *RoundRobinStrategy
	lock               *sync.RWMutex
}

//Instance ...
func (stickySession *StickySessionStrategy) Instance() *StickySessionStrategy {

	if stickySessioninstance == nil {
		stickySessioninstance = new(StickySessionStrategy)
		stickySessioninstance.Sessions = make(map[string]string, 0)
		stickySessioninstance.lock = new(sync.RWMutex)
		stickySession.roundRobinStrategy = new(RoundRobinStrategy).Instance()
	}

	return stickySessioninstance
}

//Next ...
func (stickySession *StickySessionStrategy) Next(sessionID string) string {
	node := stickySession.getOrAddSession(sessionID, stickySession.roundRobinStrategy.Instance().Next)
	return node
}

func (stickySession *StickySessionStrategy) getOrAddSession(sessionID string, getNextNode func() string) string {
	stickySession.lock.Lock()
	defer stickySession.lock.Unlock()

	node, ok := stickySession.Sessions[sessionID]
	if !ok {
		node = getNextNode()
		stickySession.Sessions[sessionID] = node
	}
	return node
}
