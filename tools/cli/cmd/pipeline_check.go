// Copyright (c) 2021 Terminus, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/erda-project/erda/pkg/parser/pipelineyml"
	"github.com/erda-project/erda/pkg/terminal/color_str"
	"github.com/erda-project/erda/tools/cli/command"
)

var PIPELINECHECK = command.Command{
	Name:       "check",
	ParentName: "PIPELINE",
	ShortHelp:  "validate pipeline.yml",
	Example:    "$ erda-cli pipeline check -f .erda/pipelines/pipeline.yml",
	Flags: []command.Flag{
		command.StringFlag{
			"f",
			"filename",
			"specify the path of pipeline.yml file",
			".erda/pipelines/pipeline.yml"},
	},
	RegisterFlagCompletionFunc: map[string]interface{}{"filename": YmlCompletion},
	Run:                        PipelineCheck,
}

func YmlCompletion(ctx *cobra.Command, args []string, toComplete string, filename string) []string {
	comps := []string{}

	p, err := getWorkspacePipelines()
	if err == nil {
		comps = p
	}
	return comps
}

func PipelineCheck(ctx *command.Context, ymlfile string) error {
	var b []byte
	if ymlfile == "-" {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		b = bytes
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		b, err = ioutil.ReadFile(filepath.Join(wd, ymlfile))
		if err != nil {
			return err
		}
	}

	_, err := pipelineyml.New(b)
	if err != nil {
		return err
	}

	fmt.Println(color_str.Green("✔ "), ymlfile, "ok!")
	return nil
}
