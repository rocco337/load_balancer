package main

import "container/ring"

var (
	instance *RoundRobinStrategy
)

//RoundRobinStrategy ...
type RoundRobinStrategy struct {
	nodesList *ring.Ring
}

//Instance ...
func (roundRobin *RoundRobinStrategy) Instance() *RoundRobinStrategy {

	if instance == nil {
		instance = new(RoundRobinStrategy)
	}

	return instance
}

//Init ...
func (roundRobin *RoundRobinStrategy) Init(nodes []string) {
	roundRobin.nodesList = ring.New(len(nodes))

	for i := 0; i < len(nodes); i++ {
		roundRobin.nodesList.Value = nodes[i]
		roundRobin.nodesList = roundRobin.nodesList.Next()
	}
}

//Next ...
func (roundRobin *RoundRobinStrategy) Next() string {
	defer func() {
		roundRobin.nodesList = roundRobin.nodesList.Next()
	}()

	return roundRobin.nodesList.Value.(string)
}
