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

var OPEN = command.Command{
	Name:      "open",
	ShortHelp: "open the web page in browser",
	Example:   "$ erda-cli open --application=<name>",
	Flags: []command.Flag{
		command.StringFlag{Short: "", Name: "org", Doc: "the name of an organization", DefaultValue: ""},
		command.StringFlag{Short: "", Name: "project", Doc: "the name of a project", DefaultValue: ""},
		command.StringFlag{Short: "", Name: "application", Doc: "the name of an application ", DefaultValue: ""},
	},
	Run: Open,
}

func Open(ctx *command.Context, //orgId, projectId, applicationId uint64,
	org, project, application string) error {

	if org != "" {
		OrgOpen(ctx, org)
	} else if project != "" {
		ProjectOpen(ctx, project)
	} else if application != "" {
		ApplicationOpen(ctx, application)
	} else if ctx.CurrentApplication.Name != "" {
		ApplicationOpen(ctx, ctx.CurrentApplication.Name)
	} else if ctx.CurrentProject.Name != "" {
		ProjectOpen(ctx, project)
	}

	return nil
}
