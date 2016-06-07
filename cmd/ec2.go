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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
	Run: func(cmd *cobra.Command, args []string) {
	},
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
