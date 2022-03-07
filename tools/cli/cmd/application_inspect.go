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
	"strings"

	"github.com/erda-project/erda/tools/cli/command"
	"github.com/erda-project/erda/tools/cli/common"
	"github.com/erda-project/erda/tools/cli/utils"
)

var APPLICATIONINSPECT = command.Command{
	Name:       "inspect",
	ParentName: "APPLICATION",
	ShortHelp:  "inspect application",
	Example:    "$ erda-cli application inspect --application=<name>",
	Flags: []command.Flag{
		command.StringFlag{Short: "", Name: "application", Doc: "the name of an application ", DefaultValue: ""},
		command.BoolFlag{Short: "", Name: "only-repo", Doc: "If true, only show git repo url", DefaultValue: false},
	},
	Run: ApplicationInspect,
}

func ApplicationInspect(ctx *command.Context, application string, onlyRepo bool) error {
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

	resp, err := common.GetApplicationDetail(ctx, orgId, projectId, applicationId)
	if err != nil {
		return err
	}

	if onlyRepo {
		repoUrl := resp.GitRepoNew
		if !(strings.HasPrefix(repoUrl, "http://") || strings.HasPrefix(repoUrl, "https://")) {
			repoUrl = fmt.Sprintf("https://%s", repoUrl)
		}
		fmt.Println(repoUrl)
	} else {
		s, err := utils.Marshal(resp)
		if err != nil {
			return fmt.Errorf(utils.FormatErrMsg("application inspect",
				"failed to prettyjson marshal application data ("+err.Error()+")", false))
		}

		fmt.Println(string(s))
	}

	return nil
}
