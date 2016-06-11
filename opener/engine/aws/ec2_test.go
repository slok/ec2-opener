package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/slok/ec2-opener/opener/engine/aws/mock/sdk"
)

func TestDescribeInstancesByID(t *testing.T) {
	// Create mock for our EC2 engine
	ctrl := gomock.NewController(t)
	mockEC2 := mock_ec2iface.NewMockEC2API(ctrl)
	engine, err := NewEc2("")
	if err != nil {
		t.Error(err)
	}
	engine.client = mockEC2
	defer ctrl.Finish()
	// Out API mock instances
	expectedIds := []string{"i-mock1", "i-mock2", "i-mock3"}
	instance1 := &ec2.Instance{InstanceId: aws.String(expectedIds[0])}
	instance2 := &ec2.Instance{InstanceId: aws.String(expectedIds[1])}
	instance3 := &ec2.Instance{InstanceId: aws.String(expectedIds[2])}
	reservation := &ec2.Reservation{
		Instances: []*ec2.Instance{instance1, instance2, instance3},
	}
	result := &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{reservation},
	}

	mockEC2.EXPECT().DescribeInstances(gomock.Any()).Do(func(input interface{}) {
		gotInstances := input.(*ec2.DescribeInstancesInput)
		// Check API received parameters are ok
		if len(gotInstances.InstanceIds) != len(expectedIds) {
			t.Fatalf("Received different parameters that expected, got %d; want %d", len(gotInstances.InstanceIds), len(expectedIds))
		}

		for idx, i := range gotInstances.InstanceIds {
			if aws.StringValue(i) != expectedIds[idx] {
				t.Fatalf("Received wrong instance ID, got %s; want %s", aws.StringValue(i), expectedIds[idx])
			}
		}

	}).Return(result, nil)

	// Check
	output := engine.describeInstancesByID(expectedIds)
	if len(output) != len(expectedIds) {
		t.Errorf("Received wrong number of instances from AWS, got %d; want %d", len(output), len(expectedIds))
	}

	for idx, i := range output {
		if aws.StringValue(i.InstanceId) != expectedIds[idx] {
			t.Errorf("Wrong Instance ID, got %s; want %s", aws.StringValue(i.InstanceId), expectedIds[idx])
		}
	}
}
