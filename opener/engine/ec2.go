package engine

import (
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/slok/ec2-opener/rule"
)

// Ec2 representes the ec2 opener logic
type Ec2 struct {
	// Ec2 session
	session *ec2.EC2

	// Ec2 instances
	instances []*ec2.Instance
}

// getInstancesByID private method to get the instances from AWS querying by IDs
func (e *Ec2) getInstancesByID(ids []string) []*ec2.Instance {
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
		resp, err := e.session.DescribeInstances(params)
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

// NewEc2ByIDs creates an Ec2 engine based on instance IDs
func NewEc2ByIDs(ids []string) (*Ec2, error) {
	// Simple engine object with an EC2 session
	svc := ec2.New(session.New(), &aws.Config{})
	if svc == nil {
		return nil, errors.New("Could not connecto with AWS")
	}
	e := &Ec2{
		session: svc,
	}

	// start populating the object
	is := e.getInstancesByID(ids)
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
func NewEc2ByTags(tags []string) (*Ec2, error) {
	return nil, nil
}

// Open opens the rules on ec2 instances
func (e *Ec2) Open(rules []*rule.Rule) error {
	return nil
}

// Close closes the rules on ec2 instnaces
func (e *Ec2) Close() error {
	return nil
}

// Clean cleans the generated stuff to open the ec2 instances
func (e *Ec2) Clean() error {
	return nil
}
