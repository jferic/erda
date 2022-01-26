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
	"os"
	"os/exec"
	"strconv"

	"github.com/erda-project/erda/apistructs"

	"github.com/pkg/errors"

	"github.com/erda-project/erda/tools/cli/command"
	"github.com/erda-project/erda/tools/cli/common"
	"github.com/erda-project/erda/tools/cli/utils"
)

var CLONE = command.Command{
	Name:      "clone",
	ShortHelp: "clone project or application from Erda",
	Example:   "$ erda-cli clone https://erda.cloud/trial/dop/projects/599",
	Args: []command.Arg{
		command.StringArg{}.Name("url"),
	},
	Run: Clone,
}

func Clone(ctx *command.Context, ustr string) error {
	var org string
	var orgId uint64
	//var project string
	var projectId uint64
	//var application string
	var applicationId uint64

	u, err := url.Parse(ustr)
	if err != nil {
		return err
	}

	t, paths, err := utils.ClassifyURL(u.Path)
	switch t {
	case utils.ApplicatinURL:
		applicationId, err = strconv.ParseUint(paths[6], 10, 64)
		if err != nil {
			return errors.Errorf("Invalid erda url.")
		}
		fallthrough
	case utils.ProjectURL:
		org = paths[1]
		projectId, err = strconv.ParseUint(paths[4], 10, 64)
		if err != nil {
			return errors.Errorf("Invalid erda url.")
		}
		break
	default:
		return errors.Errorf("Invalid erda url.")
	}

	orgId, err = getOrgId(ctx, org, orgId)
	if err != nil {
		return err
	}

	// init project
	if t == utils.ProjectURL {

		_, _, err := command.GetProjectConfig()
		if err != nil && err != utils.NotExist {
			return err
		} else if err == nil {
			return errors.New("you are already in a erda project workspace.")
		}

		p, err := common.GetProjectDetail(ctx, orgId, projectId)
		if err != nil {
			return err
		}
		pInfo := command.ProjectInfo{
			command.ConfigVersion,
			"https://openapi.erda.cloud", // TODO
			org,
			orgId,
			p.Name,
			projectId,
		}

		for _, d := range []string{
			p.Name,
			fmt.Sprintf("%s/%s", p.Name, utils.GlobalErdaDir),
		} {
			err = os.MkdirAll(d, 0755)
			if err != nil {
				return err
			}
		}

		pconfig := fmt.Sprintf("%s/.erda.d/config", p.Name)
		err = command.SetProjectConfig(pconfig, &pInfo)
		if err != nil {
			return err
		}
		ctx.Succ("Project '%s' cloned.", p.Name)
	} else if t == utils.ApplicatinURL { // init application
		a, err := common.GetApplicationDetail(ctx, orgId, projectId, applicationId)
		if err != nil {
			return err
		}

		repo := fmt.Sprintf("%s://%s", u.Scheme, a.GitRepoNew)

		err = cloneApplication(a, repo)
		if err != nil {
			return err
		}

		ctx.Succ("Application '%s' cloned.", a.Name)
	}

	return nil
}

func cloneApplication(a apistructs.ApplicationDTO, repo string) error {
	_, pInfo, err := command.GetProjectConfig()
	if err != nil {
		if err == utils.NotExist {
			return errors.New("current workspace is not an erda project.")
		}
		return err
	}

	if pInfo.ProjectId != a.ProjectID || pInfo.OrgId != a.OrgID {
		return errors.Errorf("application %s/%s cloned is not belong to project %s in the current workspace",
			a.ProjectName, a.Name, pInfo.Project)
	}

	// clone code
	_, err = exec.Command("git", "clone", repo).Output()
	if err != nil {
		fmt.Printf("git clone repo err: %v", err)
		return err
	}
	err = createApplicationDir(*pInfo, a.Name, a.ID)
	if err != nil {
		return err
	}

	return nil
}

func createApplicationDir(pInfo command.ProjectInfo, name string, applicationId uint64) error {
	aInfo := command.ApplicationInfo{
		pInfo,
		name,
		applicationId,
	}

	for _, d := range []string{
		name,
		fmt.Sprintf("%s/%s", name, utils.GlobalErdaDir),
	} {
		err := os.MkdirAll(d, 0755)
		if err != nil {
			return err
		}
	}

	aconfig := fmt.Sprintf("%s/.erda.d/config", name)
	err := command.SetApplicationConfig(aconfig, &aInfo)
	if err != nil {
		return err
	}

	return nil
}
