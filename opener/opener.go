package opener

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/slok/ec2-opener/engine"
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
	Rules  []*Rule
	Engine engine.Engine
	Status Status
}

func randomID() string {
	s := rand.NewSource(time.Now().UnixNano())
	name := fmt.Sprintf("ec2-opener-%d", rand.New(s).Intn(10000000))
	return name
}

// NewOpener Creates a new Opener
func NewOpener(rules []*Rule, engine engine.Engine) (*Opener, error) {
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
	o.Status = Open
	return nil
}
