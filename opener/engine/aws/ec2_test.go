package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/mock/gomock"
	"github.com/slok/ec2-opener/opener/engine/aws/mock"
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
	// Our API mock instances
	expectedIds := []string{"i-mock1", "i-mock2", "i-mock3"}
	mock.SetDescribeInstancesSDK(t, mockEC2, expectedIds)

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
