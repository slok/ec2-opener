package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	// EC2 session
	svc *ec2.EC2

	// Created security group
	sgID string

	// Instance
	instance *ec2.Instance
)

// ingressRule describes a rule for the security group
type ingressRule struct {
	port     int
	cidr     string
	protocol string
}

func exit(code int) {
	cleanup()
	os.Exit(code)
}

func cleanup() {
	fmt.Println("Cleaning...")
	// Unset security group
	sgs := []string{}
	for _, sg := range instance.SecurityGroups {
		sgs = append(sgs, *sg.GroupId)
	}
	setSecurityGroupsOnInstance(*instance.InstanceId, sgs)
	err := deleteSecurityGroup(sgID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not delete '%s' security group: %v\n", sgID, err)
		os.Exit(1)
	}
	fmt.Println("Clean complete, good bye")
}

//createSecurityGroup creates security group
func createSecurityGroup(vpcID string) (groupID string, err error) {

	s := rand.NewSource(time.Now().UnixNano())
	groupName := fmt.Sprintf("tmp-%d", rand.New(s).Intn(10000000))

	params := &ec2.CreateSecurityGroupInput{
		Description: aws.String("Temporal security group"),
		GroupName:   aws.String(groupName),
		VpcId:       aws.String(vpcID),
	}
	resp, err := svc.CreateSecurityGroup(params)

	if err != nil {
		return "", err
	}

	fmt.Printf("%s security group created\n", *resp.GroupId)
	return *resp.GroupId, nil
}

// setSecurityGoupRules Set the ingress rules to security group
func setSecurityGoupRules(sgID string, rules []*ingressRule) error {
	if len(rules) == 0 {
		return errors.New("Not ingress rules available")
	}

	perms := make([]*ec2.IpPermission, len(rules))
	for i, r := range rules {
		perm := &ec2.IpPermission{
			ToPort:     aws.Int64(int64(r.port)),
			FromPort:   aws.Int64(int64(r.port)),
			IpProtocol: aws.String(r.protocol),
			IpRanges:   []*ec2.IpRange{{CidrIp: aws.String(r.cidr)}},
		}
		perms[i] = perm
	}

	params := &ec2.AuthorizeSecurityGroupIngressInput{
		IpPermissions: perms,
		GroupId:       aws.String(sgID),
	}

	_, err := svc.AuthorizeSecurityGroupIngress(params)

	fmt.Printf("Set ingress rules for %s\n", sgID)
	return err
}

// deleteSecurityGroup deletes a security group
func deleteSecurityGroup(sgID string) error {
	params := &ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(sgID),
	}
	_, err := svc.DeleteSecurityGroup(params)
	fmt.Printf("%s Security group deleted\n", sgID)
	return err
}

// getEC2Instance gets the ec2 instance
func getEC2Instance(instanceID string) (*ec2.Instance, error) {
	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	}
	resp, err := svc.DescribeInstances(params)

	if err != nil {
		return nil, err
	}

	if len(resp.Reservations) == 0 {
		return nil, fmt.Errorf("No instance with '%s' id found", instanceID)
	}

	// Always return the first one, it hosuld only be one
	return resp.Reservations[0].Instances[0], nil
}

// setSecurityGroupsOnInstance Sets the security groups on the instance
func setSecurityGroupsOnInstance(instanceID string, securityGroups []string) error {
	gs := make([]*string, len(securityGroups))

	for i, g := range securityGroups {
		gs[i] = aws.String(g)
	}

	params := &ec2.ModifyInstanceAttributeInput{
		InstanceId: aws.String(instanceID),
		Groups:     gs,
	}
	_, err := svc.ModifyInstanceAttribute(params)
	fmt.Printf("Security group set on instance\n")
	return err
}

func main() {
	region := "eu-west-1"
	var awsConf *aws.Config

	// If region then overwrite
	if region != "" {
		awsConf = &aws.Config{Region: aws.String(region)}
	}

	// Open a new Ec2 session
	svc = ec2.New(session.New(), awsConf)

	instanceID := "i-016e456f844f9ac25"

	// Get instance
	var err error
	instance, err = getEC2Instance(instanceID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving instance: %v\n", err)
		exit(1)
	}

	// Create the security group and
	sgID, err = createSecurityGroup(*instance.VpcId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error openning ports: %v\n", err)
		exit(1)
	}

	// Add rules to security group
	rules := []*ingressRule{
		&ingressRule{
			port:     22,
			cidr:     "0.0.0.0/0",
			protocol: "TCP",
		},
		&ingressRule{
			port:     3389,
			cidr:     "0.0.0.0/0",
			protocol: "TCP",
		},
	}
	err = setSecurityGoupRules(sgID, rules)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting the port rules: %v\n", err)
		exit(1)
	}

	// Set group on instance
	sgs := []string{}
	for _, sg := range instance.SecurityGroups {
		sgs = append(sgs, *sg.GroupId)
	}
	// Add our new security group
	sgs = append(sgs, sgID)

	err = setSecurityGroupsOnInstance(*instance.InstanceId, sgs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting security group on instance: %v\n", err)
		exit(1)
	}

	// Listen until ^C
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt)

	fmt.Println("Press ctr+C to close the instance port...")
	<-c
	cleanup()

}
