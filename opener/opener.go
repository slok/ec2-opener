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
	// Clean status
	Clean
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
	case Clean:
		return "clean"
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
		Status: Clean,
	}

	return o, nil
}

// Open is the action of oppenning the listeners on the atarget
func (o *Opener) Open() error {
	// already open
	if o.Status == Open {
		return nil
	}

	// check if previous status is ok
	if o.Status != Clean {
		return fmt.Errorf("you need to be on clean state before opening, current status: %s", o.Status)
	}

	// Open with the engine
	if err := o.Engine.Open(o.Rules); err != nil {
		return fmt.Errorf("error opening rules: %s", err)
	}

	// All ok, set status to open
	o.Status = Open

	return nil
}

// Close is the action of closing open listeners on the target
func (o *Opener) Close() error {
	// already closed
	if o.Status == Close {
		return nil
	}

	// check if previous status is ok
	if o.Status != Open {
		return fmt.Errorf("you need to be on open state before closing, current status: %s", o.Status)
	}

	// Close with the engine
	if err := o.Engine.Close(); err != nil {
		return fmt.Errorf("error closing rules: %s", err)
	}

	// All ok, set status
	o.Status = Close
	return nil
}

// Clean is the action of cleaning up all the stuff made when opening and closing on target
func (o *Opener) Clean() error {
	// already closed
	if o.Status == Clean {
		return nil
	}

	// check if previous status is ok
	if o.Status != Close {
		return fmt.Errorf("you need to be on closed state before cleaning up, current status: %s", o.Status)
	}

	// Clean with the engine
	if err := o.Engine.Clean(); err != nil {
		return fmt.Errorf("error cleaning rules: %s", err)
	}

	o.Status = Clean
	return nil

}
