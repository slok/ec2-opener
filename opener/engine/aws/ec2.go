package aws

import (
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"

	"github.com/slok/ec2-opener/rule"
)

// go:generate mockgen -source vendor/github.com/aws/aws-sdk-go/service/ec2/ec2iface/interface.go  -destination opener/engine/aws/mock/sdk/ec2iface_mock.go

// Ec2Engine representes the ec2 opener logic
type Ec2Engine struct {
	// Ec2 client
	client ec2iface.EC2API

	// Ec2 instances
	instances []*ec2.Instance
}

// describeInstancesByID gets the instances from AWS querying by IDs
func (e *Ec2Engine) describeInstancesByID(ids []string) []*ec2.Instance {
	result := []*ec2.Instance{}

	// Convert params for AWS
	ec2IDs := []*string{}
	for _, i := range ids {
		ec2IDs = append(ec2IDs, aws.String(i))
	}
	// Get the instances with the API
	logrus.Debugf("Getting %s instances", ids)
	nextToken := aws.String("")
	params := &ec2.DescribeInstancesInput{}

	// Make the calls (paginated)
	for nextToken != nil {
		// If there is a next token then we only need this toke in the call params
		if *nextToken != "" {
			params.NextToken = nextToken
		} else {
			params.InstanceIds = ec2IDs
		}
		// Call!
		resp, err := e.client.DescribeInstances(params)
		if err != nil {
			logrus.Error(err)
			return result
		}
		// Get the instances and append to our result
		for _, r := range resp.Reservations {
			for _, i := range r.Instances {
				result = append(result, i)
			}
		}
		// more pages?
		nextToken = resp.NextToken
	}

	logrus.Debugf("Restrieved %d instances", len(result))
	return result
}

// InitByInstancesOrTags will init the engine using instances or tags depending
// on the params
func (e *Ec2Engine) InitByInstancesOrTags(instanceIds, tags []string) error {
	// instance Ids have priority
	if len(instanceIds) != 0 {
		logrus.Debug("Populating EC2 engine by instances IDs")
		// start populating the object
		is := e.describeInstancesByID(instanceIds)

		if len(is) == 0 {
			return fmt.Errorf("No instances found")
		}
		e.instances = is
		return nil
	}

	if len(tags) != 0 {
		logrus.Debug("Populating EC2 engine by tags")
		return nil
	}

	return errors.New("Could not initialize engine, wrong instances IDs or tags")

}

// NewEc2 creates an Ec2 engine
func NewEc2(region string) (*Ec2Engine, error) {

	client := ec2.New(session.New(), &aws.Config{
		Region: aws.String(region),
	})

	if client == nil {
		return nil, errors.New("Could not connect with AWS")
	}

	e := &Ec2Engine{
		client: client,
	}

	return e, nil
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
