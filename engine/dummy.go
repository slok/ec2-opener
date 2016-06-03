package engine

// Dummy will be used to use it in the tests
type Dummy struct {
}

// NewDummy creates a new dummy engine object
func NewDummy() (*Dummy, error) {
	d := &Dummy{}
	return d, nil
}
