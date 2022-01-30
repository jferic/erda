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
	"os"

	"github.com/pkg/errors"

	"github.com/erda-project/erda/tools/cli/command"
	"github.com/erda-project/erda/tools/cli/common"
	"github.com/erda-project/erda/tools/cli/utils"
)

var PROJECTINIT = command.Command{
	Name:       "init",
	ParentName: "PROJECT",
	ShortHelp:  "init project",
	Example:    "$ erda-cli project inspect --project=<name>",
	Flags: []command.Flag{
		command.StringFlag{Short: "", Name: "org", Doc: "the name of an organization", DefaultValue: ""},
		command.StringFlag{Short: "", Name: "project", Doc: "the name of a project", DefaultValue: ""},
	},
	Run: ProjectInit,
}

func ProjectInit(ctx *command.Context, org, project string) error {
	_, _, err := command.GetProjectConfig()
	if err == nil {
		return errors.New("project already inited.")
	} else if err != utils.NotExist {
		return err
	}

	o, err := common.GetOrgDetail(ctx, org)
	if err != nil {
		return err
	}

	pId, err := common.GetProjectIdByName(ctx, o.ID, project)
	if err != nil {
		return err
	}

	pInfo := command.ProjectInfo{
		Version:   command.ConfigVersion,
		Server:    ctx.CurrentOpenApiHost,
		Org:       org,
		OrgId:     o.ID,
		Project:   project,
		ProjectId: pId,
	}

	appList, err := common.GetApplications(ctx, o.ID, pId)
	if err != nil {
		return err
	}
	for _, a := range appList {
		aInfo := command.ApplicationInfo{a.Name, a.ID}
		pInfo.Applications = append(pInfo.Applications, aInfo)
	}

	for _, d := range []string{
		fmt.Sprintf("%s", utils.GlobalErdaDir),
	} {
		err = os.MkdirAll(d, 0755)
		if err != nil {
			return err
		}
	}

	err = command.SetProjectConfig(".erda.d/config", &pInfo)
	if err != nil {
		return err
	}
	ctx.Succ("Project '%s' cloned.", project)

	return nil
}
