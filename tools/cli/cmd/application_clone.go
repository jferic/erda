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
	"net/url"

	"github.com/spf13/cobra"

	"github.com/pkg/errors"

	"github.com/erda-project/erda/tools/cli/command"
	"github.com/erda-project/erda/tools/cli/common"
	"github.com/erda-project/erda/tools/cli/utils"
)

var APPLICATIONCLONE = command.Command{
	Name:       "clone",
	ParentName: "APPLICATION",
	ShortHelp:  "clone the application",
	Example:    "$ erda-cli application clone <name>",
	Args: []command.Arg{
		command.StringArg{}.Name("application"),
	},
	ValidArgsFunction: ArgApplicationCompletion,
	Run:               ApplicationClone,
}

func ArgApplicationCompletion(ctx *cobra.Command, args []string, toComplete string, application string) []string {
	var comps []string

	err := command.PrepareCtx(ctx, args)
	if err != nil {
		return comps
	}

	var org, project string
	var orgId, projectId uint64

	c := command.GetContext()
	org, orgId, err = getOrgId(c, org, orgId)
	if err != nil {
		return comps
	}

	project, projectId, err = getProjectId(c, orgId, project, projectId)
	if err != nil {
		return comps
	}

	appList, err := common.GetApplications(c, orgId, projectId)
	if err != nil {
		return comps
	}

	for _, a := range appList {
		comps = append(comps, a.Name)
	}

	return comps
}

func ApplicationClone(ctx *command.Context, application string) error {
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

	a, err := common.GetApplicationDetail(ctx, orgId, projectId, applicationId)
	if err != nil {
		return err
	}

	u, err := url.Parse(ctx.CurrentOpenApiHost)
	if err != nil {
		return err
	}

	repo := fmt.Sprintf("%s://%s", u.Scheme, a.GitRepoNew)

	_, pInfo, err := command.GetProjectConfig()
	if err != nil {
		if err == utils.NotExist {
			return errors.New("current workspace is not an erda project.")
		}
		return err
	}

	dir := application
	err = cloneApplication(pInfo, a, repo, dir)
	if err != nil {
		return err
	}

	ctx.Succ("Application '%s' cloned.", application)
	return nil
}
