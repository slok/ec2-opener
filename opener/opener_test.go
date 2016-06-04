package opener

import (
	"testing"

	"github.com/slok/ec2-opener/opener/engine"
	"github.com/slok/ec2-opener/rule"
)

func TestNewOpener(t *testing.T) {
	tests := []struct {
		rules  []*rule.Rule
		engine engine.Engine
	}{
		{
			rules:  []*rule.Rule{&rule.Rule{}},
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
	rs := []*rule.Rule{
		&rule.Rule{
			CIDR:     "0.0.0.0/0",
			Protocol: rule.TCP,
			Port:     22,
		},
		&rule.Rule{
			CIDR:     "54.12.45.0/24",
			Protocol: rule.UDP,
			Port:     9873,
		},
	}
	o, err := NewOpener(rs, e)

	if err != nil {
		t.Errorf("Got error while creating opener: %s", err)
	}

	if o.Status != Clean {
		t.Errorf("Got wrong status: %s; want: %s", o.Status, Clean)
	}

	// Open
	o.Open()

	if o.Status != Open {
		t.Errorf("Got wrong status: %s; want: %s", o.Status, Open)
	}

	// Check openned rules
	for _, r := range rs {
		if got, ok := e.OpenRules[r.String()]; !ok {
			t.Errorf("Rule %s, not opened", r)
		} else {
			if got.String() != r.String() {
				t.Errorf("Not expected rule: got: %s; want: %s", got, r)
			}
		}
	}
}

func TestOpenerClose(t *testing.T) {
	e, _ := engine.NewDummy()
	rs := []*rule.Rule{
		&rule.Rule{
			CIDR:     "0.0.0.0/0",
			Protocol: rule.TCP,
			Port:     22,
		},
		&rule.Rule{
			CIDR:     "54.12.45.0/24",
			Protocol: rule.UDP,
			Port:     9873,
		},
	}
	o, _ := NewOpener(rs, e)
	o.Open()

	if o.Status != Open {
		t.Errorf("Got wrong status: %s; want: %s", o.Status, Open)
	}

	o.Close()

	if o.Status != Close {
		t.Errorf("Got wrong status: %s; want: %s", o.Status, Close)
	}

	// Check closed rules
	for _, r := range rs {
		if got, ok := e.CloseRules[r.String()]; !ok {
			t.Errorf("Rule %s, not closed", r)
		} else {
			if got.String() != r.String() {
				t.Errorf("Not expected rule: got: %s; want: %s", got, r)
			}
		}
	}

}
