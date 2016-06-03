package engine

import "github.com/slok/ec2-opener/rule"

// Dummy will be used to use it in the tests
type Dummy struct {
	OpenRules map[string]*rule.Rule
}

// NewDummy creates a new dummy engine object
func NewDummy() (*Dummy, error) {
	d := &Dummy{
		OpenRules: map[string]*rule.Rule{},
	}
	return d, nil
}

// Open stores the rules on memory until they are closed
func (d *Dummy) Open(rules []*rule.Rule) error {

	// Add rules to the openned memory storage
	for _, r := range rules {
		d.OpenRules[r.String()] = r
	}

	return nil
}
