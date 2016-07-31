package mock

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"

	"github.com/slok/ec2-opener/opener/engine/aws/mock/sdk"
)

// Instance is a mock object to create the instance
type Instance struct {
	ID    string
	VpcID string
}

// SetDescribeInstancesSDK mocks describe instances call to EC2 SDK
func SetDescribeInstancesSDK(t *testing.T, mockMatcher *mock_ec2iface.MockEC2API, instanceIDs []*Instance) {
	// Out API mock instances
	instances := make([]*ec2.Instance, len(instanceIDs))
	for idx, i := range instanceIDs {
		instances[idx] = &ec2.Instance{InstanceId: aws.String(i.ID), VpcId: aws.String(i.VpcID)}
	}

	reservation := &ec2.Reservation{
		Instances: instances,
	}
	result := &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{reservation},
	}

	mockMatcher.EXPECT().DescribeInstances(gomock.Any()).Do(func(input interface{}) {
		gotInstances := input.(*ec2.DescribeInstancesInput)
		// Check API received parameters are ok
		if len(gotInstances.InstanceIds) != len(instanceIDs) {
			t.Fatalf("Received different parameters that expected, got %d; want %d", len(gotInstances.InstanceIds), len(instanceIDs))
		}

		for idx, i := range gotInstances.InstanceIds {
			if aws.StringValue(i) != instanceIDs[idx].ID {
				t.Fatalf("Received wrong instance ID, got %s; want %s", aws.StringValue(i), instanceIDs[idx].ID)
			}
		}

	}).AnyTimes().Return(result, nil)
}

// SetCreateSecurityGroupSDK mocks creating security group on an VPC
func SetCreateSecurityGroupSDK(t *testing.T, mockMatcher *mock_ec2iface.MockEC2API) {
	result := &ec2.CreateSecurityGroupOutput{}

	mockMatcher.EXPECT().CreateSecurityGroup(gomock.Any()).Do(func(input interface{}) {
		sgInput := input.(*ec2.CreateSecurityGroupInput)
		if aws.StringValue(sgInput.GroupName) == "" || sgInput.GroupName == nil {
			t.Fatalf("Received wrong group name parameter")
		}

		// Set the group id
		result.GroupId = aws.String(fmt.Sprintf("%s-id", aws.StringValue(sgInput.GroupName)))

	}).AnyTimes().Return(result, nil)

}

// SetCreateSecurityGroupWithErrorSDK mocks creating security group on an VPC and error on the X call
func SetCreateSecurityGroupWithErrorSDK(t *testing.T, mockMatcher *mock_ec2iface.MockEC2API, errorTime int) {
	result := &ec2.CreateSecurityGroupOutput{}

	// This is the number of calls that will return ok
	call1 := mockMatcher.EXPECT().CreateSecurityGroup(gomock.Any()).Do(func(input interface{}) {
		sgInput := input.(*ec2.CreateSecurityGroupInput)
		if aws.StringValue(sgInput.GroupName) == "" || sgInput.GroupName == nil {
			t.Fatalf("Received wrong group name parameter")
		}

		// Set the group id
		result.GroupId = aws.String(fmt.Sprintf("%s-id", aws.StringValue(sgInput.GroupName)))

	}).MaxTimes(errorTime-1).Return(result, nil)

	// This time will fail
	mockMatcher.EXPECT().CreateSecurityGroup(gomock.Any()).After(call1).Return(nil, errors.New("Error on call"))

}

// SetAuthorizeSecurityGroupIngressSDK mocks the set of rules on security groups
func SetAuthorizeSecurityGroupIngressSDK(t *testing.T, mockMatcher *mock_ec2iface.MockEC2API) {
	result := &ec2.AuthorizeSecurityGroupIngressOutput{}
	mockMatcher.EXPECT().AuthorizeSecurityGroupIngress(gomock.Any()).Do(func(input interface{}) {
		sgInput := input.(*ec2.AuthorizeSecurityGroupIngressInput)
		if len(sgInput.IpPermissions) == 0 {
			t.Fatalf("Received empty permissions")
		}

		if aws.StringValue(sgInput.GroupId) == "" {
			t.Fatalf("Received empty Group ID")
		}
	}).AnyTimes().Return(result, nil)
}
