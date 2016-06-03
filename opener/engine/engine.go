package engine

import "github.com/slok/ec2-opener/rule"

// Engine representes the interface every engine need to implement in order so an opener can use it
type Engine interface {
	// Open rules on the target
	Open(rules []*rule.Rule) error
}
