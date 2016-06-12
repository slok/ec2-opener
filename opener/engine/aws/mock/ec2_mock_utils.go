package mock

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"

	"github.com/slok/ec2-opener/opener/engine/aws/mock/sdk"
)

// SetDescribeInstancesSDK mocks describe instances call to EC2 SDK
func SetDescribeInstancesSDK(t *testing.T, mockMatcher *mock_ec2iface.MockEC2API, instanceIDs []string) {
	// Out API mock instances
	instances := make([]*ec2.Instance, len(instanceIDs))
	for idx, i := range instanceIDs {
		instances[idx] = &ec2.Instance{InstanceId: aws.String(i)}
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
			if aws.StringValue(i) != instanceIDs[idx] {
				t.Fatalf("Received wrong instance ID, got %s; want %s", aws.StringValue(i), instanceIDs[idx])
			}
		}

	}).Return(result, nil)
}
