package client

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Ec2Client is an ec2 client to perform the calls to the ec2 AWS API
type Ec2Client interface {
	DescribeInstancesByID(ids []string) []*ec2.Instance
}

// Ec2APIClient is the client that will be contacting to the AWS API
type Ec2APIClient struct {
	session *ec2.EC2
}

// NewEc2APIClient Creates and returns a new ec2APIClient
func NewEc2APIClient() (*Ec2APIClient, error) {
	// Simple engine object with an EC2 session
	e := &Ec2APIClient{}
	e.session = ec2.New(session.New(), &aws.Config{})
	if e.session == nil {
		return nil, errors.New("Could not connecto with AWS")
	}

	return e, nil
}

// DescribeInstancesByID gets the instances from AWS querying by IDs
func (e *Ec2APIClient) DescribeInstancesByID(ids []string) []*ec2.Instance {
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
