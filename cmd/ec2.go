// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/slok/ec2-opener/opener"
	"github.com/slok/ec2-opener/opener/engine/aws"
	"github.com/slok/ec2-opener/rule"
)

func exit(opener *opener.Opener, code int) {
	logrus.Debugf("Cleaning up")
	// cleanup all the mess
	opener.Clean()

	// Finish program
	logrus.Debugf("Exiting")
	time.Sleep(1 * time.Second)
	os.Exit(code)
}

func ec2Main(cmd *cobra.Command, args []string) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugf("Call EC2 opener")
	if logrus.GetLevel() == logrus.DebugLevel {
		cmd.DebugFlags()
	}

	// Create opener
	e, err := aws.NewEc2ByIDs(instances)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	rs := []*rule.Rule{}
	o, err := opener.NewOpener(rs, e)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	// Open rules
	err = o.Open()
	if err != nil {
		logrus.Error(err)
		exit(o, 1)
	}

	// Listen until ^C
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt)

	logrus.Debugf("Press ctr+C to close the instance port")
	<-c
	err = o.Close()

	if err != nil {
		logrus.Error(err)
		exit(o, 1)
	}

	exit(o, 0)
}

// Command flags
var (
	instances []string
	tags      []string
)

// ec2Cmd represents the ec2 command
var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "EC2 engine for the target opening",
	Long:  ``,
	Run:   ec2Main,
}

func init() {
	RootCmd.AddCommand(ec2Cmd)

	// Check
	// https://github.com/spf13/viper/issues/112

	ec2Cmd.Flags().StringSliceVarP(&instances, "instances", "i", []string{}, `	Accepts a comma separated list of instances that are going to identify the instances.
				You can also use multiple flag statements like -i i-xxxxx -i i-yyyyy...; The presence
				of this argument invalidates '--tags' o '-t' argument`)
	viper.BindPFlag("instances", ec2Cmd.Flags().Lookup("instances"))

	ec2Cmd.Flags().StringSliceVarP(&tags, "tags", "t", []string{}, `	Accepts a comma separated list of tags that are going to use as filters. You can also
				use multiple flag statements like -t tag1 -t tag2...; The presence '--instances' or '-i'
				argument invalidates this argument`)
	viper.BindPFlag("tags", ec2Cmd.Flags().Lookup("tags"))
}
