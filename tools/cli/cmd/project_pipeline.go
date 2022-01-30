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
	"github.com/erda-project/erda/tools/cli/command"
)

var PROJECTPIPELINE = command.Command{
	Name:       "pipeline",
	ParentName: "PROJECT",
	ShortHelp:  "list pipelines in the project",
	Example:    "$ erda-cli project pipeline",
	Flags: []command.Flag{
		command.StringFlag{Short: "", Name: "name", Doc: "name of the pipeline", DefaultValue: ""},
		command.StringFlag{Short: "", Name: "application", Doc: "the name of a application", DefaultValue: ""},
		command.StringFlag{Short: "", Name: "ref", Doc: "the branch name", DefaultValue: ""},
		command.StringFlag{Short: "", Name: "filename", Doc: "the filename of pipeline yaml", DefaultValue: ""},
	},
	Run: ProjectPipeline,
}

func ProjectPipeline(ctx *command.Context, name, application, ref, filename string) error {

	var org, project string
	var orgId, projectId uint64

	org, orgId, err := getOrgId(ctx, org, orgId)
	if err != nil {
		return err
	}

	project, projectId, err = getProjectId(ctx, orgId, project, projectId)
	if err != nil {
		return err
	}

	ctx.Succ("TODO", project)
	return nil
}
