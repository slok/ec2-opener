package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/slok/ec2-opener/opener/engine/aws/mock"
	"github.com/slok/ec2-opener/opener/engine/aws/mock/sdk"
	"github.com/slok/ec2-opener/rule"
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
	expectedInstances := []*mock.Instance{
		&mock.Instance{ID: "i-mock1", VpcID: "vpc1"},
		&mock.Instance{ID: "i-mock2", VpcID: "vpc1"},
		&mock.Instance{ID: "i-mock3", VpcID: "vpc1"}}
	mock.SetDescribeInstancesSDK(t, mockEC2, expectedInstances)

	expectedIds := make([]string, len(expectedInstances))
	for i, v := range expectedInstances {
		expectedIds[i] = v.ID
	}

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

func TestInitWithoutInstancesOrTags(t *testing.T) {
	engine, err := NewEc2("")
	if err != nil {
		t.Error(err)
	}

	err = engine.InitByInstancesOrTags(nil, nil)
	if err == nil {
		t.Error("Initialization without instances or tags should return and error")
	}

}

func TestInitWithInstances(t *testing.T) {
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
	expectedInstances := []*mock.Instance{
		&mock.Instance{ID: "i-mock1", VpcID: "vpc1"},
		&mock.Instance{ID: "i-mock2", VpcID: "vpc1"},
		&mock.Instance{ID: "i-mock3", VpcID: "vpc1"},
		&mock.Instance{ID: "i-mock4", VpcID: "vpc1"}}
	mock.SetDescribeInstancesSDK(t, mockEC2, expectedInstances)

	expectedIds := make([]string, len(expectedInstances))
	for i, v := range expectedInstances {
		expectedIds[i] = v.ID
	}

	engine.InitByInstancesOrTags(expectedIds, nil)

	if len(engine.instances) != len(expectedIds) {
		t.Errorf("Wrong number of instances from AWS, got %d; want %d", len(engine.instances), len(expectedIds))
	}

	for idx, i := range engine.instances {
		if aws.StringValue(i.InstanceId) != expectedIds[idx] {
			t.Errorf("Wrong Instance ID, got %s; want %s", aws.StringValue(i.InstanceId), expectedIds[idx])
		}
	}
}

func TestCreateSecurityGroupsSingleVPC(t *testing.T) {
	// Create mock for our EC2 engine
	ctrl := gomock.NewController(t)
	mockEC2 := mock_ec2iface.NewMockEC2API(ctrl)

	engine, err := NewEc2("")
	if err != nil {
		t.Error(err)
	}
	engine.client = mockEC2
	defer ctrl.Finish()

	// Mock ec2 API
	mock.SetCreateSecurityGroupSDK(t, mockEC2)

	// Set our instances
	engine.instances = []*ec2.Instance{
		&ec2.Instance{InstanceId: aws.String("i-mock1"), VpcId: aws.String("vpc1")},
		&ec2.Instance{InstanceId: aws.String("i-mock2"), VpcId: aws.String("vpc1")},
	}
	// Create the required security groups
	if err := engine.createSecurityGroups([]*rule.Rule{}); err != nil {
		t.Errorf("Failed creating security groups: %s", err)
	}

	if len(engine.createdSGPerVPC) != 1 {
		t.Errorf("Should create one security group, intead created %d", len(engine.createdSGPerVPC))
	}
}

func TestCreateSecurityGroupsMultipleVPC(t *testing.T) {
	// Create mock for our EC2 engine
	ctrl := gomock.NewController(t)
	mockEC2 := mock_ec2iface.NewMockEC2API(ctrl)

	engine, err := NewEc2("")
	if err != nil {
		t.Error(err)
	}
	engine.client = mockEC2
	defer ctrl.Finish()

	// Mock ec2 API
	mock.SetCreateSecurityGroupSDK(t, mockEC2)

	// Set our instances
	engine.instances = []*ec2.Instance{
		&ec2.Instance{InstanceId: aws.String("i-mock1"), VpcId: aws.String("vpc1")},
		&ec2.Instance{InstanceId: aws.String("i-mock2"), VpcId: aws.String("vpc2")},
	}
	// Create the required security groups
	if err := engine.createSecurityGroups([]*rule.Rule{}); err != nil {
		t.Errorf("Failed creating security groups: %s", err)
	}

	if len(engine.createdSGPerVPC) != len(engine.instances) {
		t.Errorf("Invalid number of security groups created, got: %d; want: %d", len(engine.createdSGPerVPC), len(engine.instances))
	}
}

func TestCreateSecurityGroupsError(t *testing.T) {
	// Create mock for our EC2 engine
	ctrl := gomock.NewController(t)
	mockEC2 := mock_ec2iface.NewMockEC2API(ctrl)

	engine, err := NewEc2("")
	if err != nil {
		t.Error(err)
	}
	engine.client = mockEC2
	defer ctrl.Finish()

	// Mock ec2 API
	// Fail on third time, will only create 2 security groups
	errorTime := 3
	mock.SetCreateSecurityGroupWithErrorSDK(t, mockEC2, errorTime)

	// Set our instances
	engine.instances = []*ec2.Instance{
		&ec2.Instance{InstanceId: aws.String("i-mock1"), VpcId: aws.String("vpc1")},
		&ec2.Instance{InstanceId: aws.String("i-mock2"), VpcId: aws.String("vpc2")},
		&ec2.Instance{InstanceId: aws.String("i-mock3"), VpcId: aws.String("vpc1")},
		&ec2.Instance{InstanceId: aws.String("i-mock4"), VpcId: aws.String("vpc3")},
	}
	// Create the required security groups
	if err := engine.createSecurityGroups([]*rule.Rule{}); err == nil {
		t.Errorf("Should return an error")
	}

	expectedNumber := errorTime - 1
	if len(engine.createdSGPerVPC) != expectedNumber {
		t.Errorf("Invalid number of security groups created, got: %d; want: %d", len(engine.createdSGPerVPC), expectedNumber)
	}
}

func TestSetSecurityGroupRules(t *testing.T) {
	// Create mock for our EC2 engine
	ctrl := gomock.NewController(t)
	mockEC2 := mock_ec2iface.NewMockEC2API(ctrl)

	engine, err := NewEc2("")
	if err != nil {
		t.Error(err)
	}
	engine.client = mockEC2
	defer ctrl.Finish()

	// Mock ec2 API
	mock.SetAuthorizeSecurityGroupIngressSDK(t, mockEC2)

	// Set our security groups
	engine.createdSGPerVPC = map[string]string{
		"vpc1": "sg-1251",
		"vpc2": "sg-1252",
		"vpc3": "sg-1253",
		"vpc4": "sg-1254",
	}

	rs := []*rule.Rule{
		&rule.Rule{Protocol: rule.TCP, CIDR: "0.0.0.0/0", Port: 22},
		&rule.Rule{Protocol: rule.TCP, CIDR: "0.0.0.0/0", Port: 80},
		&rule.Rule{Protocol: rule.TCP, CIDR: "0.0.0.0/0", Port: 443},
	}
	err = engine.setSecurityGroupRules(rs)
	if err != nil {
		t.Errorf("Security Group rules set shouldn't fail: %s", err)
	}
}

func TestSetSecurityGroupRulesError(t *testing.T) {

	tests := []struct {
		SGs   map[string]string
		rules []*rule.Rule
	}{
		{
			SGs:   map[string]string{},
			rules: []*rule.Rule{&rule.Rule{Protocol: rule.TCP, CIDR: "0.0.0.0/0", Port: 22}},
		},
		{
			SGs:   map[string]string{"vpc1": "sg-1251", "vpc2": "sg-1252"},
			rules: []*rule.Rule{},
		},
	}

	// Create mock for our EC2 engine
	ctrl := gomock.NewController(t)
	mockEC2 := mock_ec2iface.NewMockEC2API(ctrl)

	for _, test := range tests {

		engine, err := NewEc2("")
		if err != nil {
			t.Error(err)
		}
		engine.client = mockEC2
		defer ctrl.Finish()

		// Mock ec2 API
		mock.SetAuthorizeSecurityGroupIngressSDK(t, mockEC2)

		// Set our security groups
		engine.createdSGPerVPC = test.SGs
		err = engine.setSecurityGroupRules(test.rules)
		if err == nil {
			t.Errorf("%+v\n - Security Group rules set should fail, it dind't", test)
		}
	}

}
