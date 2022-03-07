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

var APPLICATIONOPEN = command.Command{
	Name:       "open",
	ParentName: "APPLICATION",
	ShortHelp:  "open the application page in browser",
	Example:    "$ erda-cli application open --application=<name>",
	Flags: []command.Flag{
		command.StringFlag{Short: "", Name: "application", Doc: "the name of an application ", DefaultValue: ""},
	},
	Run: ApplicationOpen,
}

func ApplicationOpen(ctx *command.Context, application string) error {
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

	err = common.Open(ctx, common.AppEntity, org, orgId, projectId, applicationId)
	if err != nil {
		return err
	}

	if application == "" {
		application = ctx.CurrentApplication.Name
	}
	ctx.Succ("Open application '%s' in browser.", application)
	return nil
}
