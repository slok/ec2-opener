package opener

import (
	"testing"

	"github.com/slok/ec2-opener/engine"
)

func TestNewOpener(t *testing.T) {
	tests := []struct {
		rules  []*Rule
		engine engine.Engine
	}{
		{
			rules:  []*Rule{&Rule{}},
			engine: &engine.Dummy{},
		},
	}

	for _, test := range tests {
		want, err := NewOpener(test.rules, test.engine)
		// Check for errors
		if err != nil {
			t.Errorf("For '%v', got error: %s", test, err)
		}

		if want == nil {
			t.Errorf("For '%v', got nil value", test)
		}
	}
}

func TestOpenerOpen(t *testing.T) {
	e, _ := engine.NewDummy()
	o, err := NewOpener(nil, e)

	if err != nil {
		t.Errorf("Got error while creating opener: %s", err)
	}

	if o.Status != Close {
		t.Errorf("Got wrong status: %s; want: %s", o.Status, Close)
	}

	// Open
	o.Open()

	if o.Status != Open {
		t.Errorf("Got wrong status: %s; want: %s", o.Status, Open)
	}
}
