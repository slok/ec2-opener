package aws

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/slok/ec2-opener/opener/engine/aws/client"
	"github.com/slok/ec2-opener/rule"
)

// Ec2Engine representes the ec2 opener logic
type Ec2Engine struct {
	// Ec2 client
	client client.Ec2Client

	// Ec2 instances
	instances []*ec2.Instance
}

// NewEc2ByIDs creates an Ec2 engine based on instance IDs
func NewEc2ByIDs(ids []string) (*Ec2Engine, error) {

	client, err := client.NewEc2APIClient()
	if err != nil {
		logrus.Error(err)
		return nil, fmt.Errorf("Could not connect to AWS")
	}

	e := &Ec2Engine{
		client: client,
	}

	// start populating the object
	is := e.client.DescribeInstancesByID(ids)

	if len(is) == 0 {
		return nil, fmt.Errorf("No instances found")
	}
	e.instances = is

	// Always return the first one, it hosuld only be one
	//return resp.Reservations[0].Instances[0], nil

	// Search for the IDs
	return e, nil
}

// NewEc2ByTags creates an Ec2 engine based on instance tags
func NewEc2ByTags(tags []string) (*Ec2Engine, error) {
	return nil, nil
}

// Open opens the rules on ec2 instances
func (e *Ec2Engine) Open(rules []*rule.Rule) error {
	return nil
}

// Close closes the rules on ec2 instnaces
func (e *Ec2Engine) Close() error {
	return nil
}

// Clean cleans the generated stuff to open the ec2 instances
func (e *Ec2Engine) Clean() error {
	return nil
}
