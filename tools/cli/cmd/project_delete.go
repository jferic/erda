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
	"github.com/erda-project/erda/tools/cli/common"
)

var PROJECTDELETE = command.Command{
	Name:       "delete",
	ParentName: "PROJECT",
	ShortHelp:  "delete project",
	Example:    "erda-cli project delete --project-id=<id>",
	Flags: []command.Flag{
		command.BoolFlag{Short: "", Name: "clear", Doc: "if true, clear runtimes and addon first", DefaultValue: false},
		command.IntFlag{Short: "", Name: "wait-runtime", Doc: "minutes to wait runtimes deleted", DefaultValue: 3},
		command.IntFlag{Short: "", Name: "wait-addon", Doc: "minutes to wait addons deleted", DefaultValue: 3},
	},
	Run: ProjectDelete,
}

func ProjectDelete(ctx *command.Context, clear bool, waitRuntime, waitAddon int) error {
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

	if clear {
		err = clearProject(ctx, orgId, projectId, "", waitRuntime, waitAddon, true, true)
		if err != nil {
			return err
		}
	}

	err = common.DeleteProject(ctx, orgId, projectId)
	if err != nil {
		return err
	}

	if project == "" {
		project = ctx.CurrentProject.Name
	}
	ctx.Succ("Project '%s' deleted.", project)
	return nil
}
