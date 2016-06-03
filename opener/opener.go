package opener

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/slok/ec2-opener/engine"
)

// Opener will desribe the way that will open the Instance(s)
type Opener struct {
	ID     string
	Rules  []*Rule
	Engine engine.Engine
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
	}

	return o, nil

}
