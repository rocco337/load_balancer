package main

var (
	stickySessioninstance *StickySessionStrategy
)

//StickySessionStrategy ...
type StickySessionStrategy struct {
	Sessions map[string]string

	roundRobinStrategy *RoundRobinStrategy
}

//Instance ...
func (stickySession *StickySessionStrategy) Instance() *StickySessionStrategy {

	if stickySessioninstance == nil {
		stickySessioninstance = new(StickySessionStrategy)
		stickySessioninstance.Sessions = make(map[string]string, 0)
	}

	return stickySessioninstance
}

//GetNode ...
func (stickySession *StickySessionStrategy) Next(sessionId string) string {
	if node, ok := stickySession.Sessions[sessionId]; ok {
		return node
	} else {
		nextNode := stickySession.roundRobinStrategy.Instance().Next()
		stickySession.Sessions[sessionId] = nextNode
		return nextNode

	}

}
