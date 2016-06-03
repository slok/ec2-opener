// Package opener package contains all the logic to open access to an EC2 instance or instances
package opener

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/slok/ec2-opener/opener/engine"
	"github.com/slok/ec2-opener/rule"
)

// Status type indicates the status of the opener
type Status int

const (
	// Open status
	Open Status = iota
	// Close status
	Close
	// Error status
	Error
)

func (s Status) String() string {
	switch s {
	case Open:
		return "open"
	case Close:
		return "close"
	case Error:
		return "error"
	default:
		return "unknown"

	}
}

// Opener will desribe the way that will open the Instance(s)
type Opener struct {
	ID     string
	Rules  []*rule.Rule
	Engine engine.Engine
	Status Status
}

func randomID() string {
	s := rand.NewSource(time.Now().UnixNano())
	name := fmt.Sprintf("ec2-opener-%d", rand.New(s).Intn(10000000))
	return name
}

// NewOpener Creates a new Opener
func NewOpener(rules []*rule.Rule, engine engine.Engine) (*Opener, error) {
	id := randomID()
	o := &Opener{
		ID:     id,
		Rules:  rules,
		Engine: engine,
		Status: Close,
	}

	return o, nil
}

// Open is the action of oppenning the listeners on the atarget
func (o *Opener) Open() error {
	// Open with the engine
	if err := o.Engine.Open(o.Rules); err != nil {
		return fmt.Errorf("error opening rules: %s", err)
	}

	// All ok, set status to open
	o.Status = Open

	return nil
}
