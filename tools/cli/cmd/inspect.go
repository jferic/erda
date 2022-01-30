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

var INSPECT = command.Command{
	Name:      "inspect",
	ShortHelp: "inspect org/project/application",
	Example:   "$ erda-cli inspect --application=<name>",
	Flags: []command.Flag{
		//command.Uint64Flag{Short: "", Name: "org-id", Doc: "the id of an organization", DefaultValue: 0},
		//command.Uint64Flag{Short: "", Name: "project-id", Doc: "the id of a project", DefaultValue: 0},
		command.StringFlag{Short: "", Name: "org", Doc: "the name of an organization", DefaultValue: ""},
		command.StringFlag{Short: "", Name: "project", Doc: "the name of a project", DefaultValue: ""},
		command.StringFlag{Short: "", Name: "application", Doc: "the name of an application ", DefaultValue: ""},
		//command.Uint64Flag{Short: "", Name: "application-id", Doc: "the id of an application ", DefaultValue: 0},
	},
	Run: Inspect,
}

func Inspect(ctx *command.Context, //orgId, projectId, applicationId uint64,
	org, project, application string) error {

	var err error
	if org != "" {
		err = OrgInspect(ctx, org)
	} else if project != "" {
		err = ProjectInspect(ctx, project)
	} else if application != "" {
		err = ApplicationInspect(ctx, application, false)
	} else if ctx.CurrentApplication.Name != "" {
		err = ApplicationInspect(ctx, ctx.CurrentApplication.Name, false)
	} else if ctx.CurrentProject.Name != "" {
		err = ProjectInspect(ctx, project)
	}
	if err != nil {
		return err
	}

	return nil
}
