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
	"strings"

	"github.com/erda-project/erda/tools/cli/command"
	"github.com/erda-project/erda/tools/cli/common"
)

var PROJECTPIPELINECreate = command.Command{
	Name:       "pipeline",
	ParentName: "PROJECTPIPELINE",
	ShortHelp:  "create a pipeline in the project",
	Example:    "$ erda-cli project pipeline create",
	Flags: []command.Flag{
		command.StringFlag{Short: "", Name: "name", Doc: "name of the pipeline", DefaultValue: ""},
		command.StringFlag{Short: "", Name: "application", Doc: "the name of a application", DefaultValue: ""},
		command.StringFlag{Short: "", Name: "branch", Doc: "the branch name", DefaultValue: ""},
		command.StringFlag{Short: "", Name: "filename", Doc: "the filename of pipeline yaml", DefaultValue: ""},
	},
	Run: ProjectPipelineCreate,
}

func ProjectPipelineCreate(ctx *command.Context, name, application, branch, filename string) error {

	var org, project string
	var orgId, projectId, applicationId uint64

	org, orgId, err := getOrgId(ctx, org, orgId)
	if err != nil {
		return err
	}

	project, projectId, err = getProjectId(ctx, orgId, project, projectId)
	if err != nil {
		return err
	}

	application, applicationId, err = getApplicationId(ctx, orgId, projectId, application, applicationId)
	if err != nil {
		return err
	}

	path := ""
	idx := strings.LastIndex(filename, "/")
	if idx != -1 {
		path = filename[:idx]
		filename = filename[idx:]
	}

	pp, err := common.CreateProjectPipeline(ctx, orgId, projectId, applicationId, name, "erda", branch, path, filename)
	if err != nil {
		return err
	}

	ctx.Succ("Project pipeline %s created", pp.Name)
	return nil
}
