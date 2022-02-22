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
	"strings"

	"github.com/erda-project/erda/apistructs"

	"github.com/pkg/errors"

	"github.com/erda-project/erda/tools/cli/command"
	"github.com/erda-project/erda/tools/cli/common"
)

var PROJECTPUSH = command.Command{
	Name:       "push",
	ParentName: "PROJECT",
	ShortHelp:  "push project",
	Example:    "$ erda-cli project push --project-config=<path2file>",
	Flags: []command.Flag{
		command.StringFlag{Short: "", Name: "org", Doc: "the name of an organization", DefaultValue: ""},
		command.StringFlag{Short: "", Name: "project", Doc: "the name of a project", DefaultValue: ""},
		command.BoolFlag{Short: "", Name: "all", Doc: "If true, push all applications", DefaultValue: false},
		command.BoolFlag{Short: "", Name: "force", Doc: "If true, git push with --force flag", DefaultValue: false},
		command.StringListFlag{Short: "", Name: "application", Doc: "applications to push", DefaultValue: nil},
		command.StringFlag{Short: "", Name: "configfile", Doc: "config file contains applications", DefaultValue: ""},
	},
	Run: ProjectPush,
}

func ProjectPush(ctx *command.Context, org, project string, pushall, force bool, applications []string, configfile string) error {
	// TODO make it easy
	if command.ProjectConfig != "" && (org != "" || project != "") {
		return errors.Errorf("Must not both specify --project-config and --org,--project,--host")
	}

	if len(applications) > 0 && configfile != "" {
		return errors.New("Should not both set --application and --configfile")
	}

	if len(applications) == 0 && !pushall && configfile == "" {
		return errors.New("No application set to push.")
	}

	var orgId, projectId uint64
	org, orgId, err := getOrgId(ctx, org, orgId)
	if err != nil {
		return err
	}

	project, projectId, err = getProjectId(ctx, orgId, project, projectId)
	if err != nil {
		return err
	}

	existProjectList, err := common.GetApplications(ctx, orgId, projectId)
	existProjectNames := map[string]apistructs.ApplicationDTO{}
	for _, p := range existProjectList {
		existProjectNames[p.Name] = p
	}

	var applications2push []command.ApplicationInfo2
	_, c, err := command.GetProjectConfig()
	if err != nil {
		return errors.Errorf("Failed to get project config, %v", err)
	}

	if len(applications) > 0 {
		cMap := map[string]command.ApplicationInfo{}
		for _, a := range c.Applications {
			cMap[a.Application] = a
		}
		for _, app := range applications {
			if a, ok := cMap[app]; ok {
				a2 := command.ApplicationInfo2{ID: a.ApplicationId, Name: a.Application, Mode: a.Mode, Desc: a.Desc}
				applications2push = append(applications2push, a2)
			} else {
				return errors.Errorf("Failed to get application in local project.")
			}
		}
	} else if configfile != "" {
		config, err := command.GetProjectConfigFrom(configfile)
		if err != nil {
			return errors.Errorf("Failed to get application from config file %s", configfile)
		}
		for _, a := range config.Applications {
			a2 := command.ApplicationInfo2{ID: a.ApplicationId, Name: a.Application, Mode: a.Mode, Desc: a.Desc}
			applications2push = append(applications2push, a2)
		}
	} else if command.ProjectConfig != "" {
		applications2push = ctx.Applications
	} else if pushall {
		for _, a := range c.Applications {
			a2 := command.ApplicationInfo2{ID: a.ApplicationId, Name: a.Application, Mode: a.Mode, Desc: a.Desc}
			applications2push = append(applications2push, a2)
		}
	}

	if len(applications2push) == 0 {
		return errors.New("No application set to push.")
	}

	for _, a := range applications2push {
		if _, err := os.Stat(a.Name); err != nil {
			return errors.Errorf("Application %s is not found in current directory. You may change to root directory of the project.", a.Name)
		}

		var gitRepo string
		if p, ok := existProjectNames[a.Name]; ok {
			gitRepo = p.GitRepoNew
		} else {
			remoteApp, err := common.CreateApplication(ctx, projectId, a.Name, a.Mode, a.Desc)
			if err != nil {
				return err
			}
			gitRepo = remoteApp.GitRepoNew
		}

		ss := strings.Split(ctx.CurrentOpenApiHost, "://")
		if len(ss) < 1 {
			return errors.Errorf("Invalid openapi host %s", ctx.CurrentOpenApiHost)
		}
		repo := fmt.Sprintf("%s://%s", ss[0], gitRepo)

		err = common.PushApplication(a.Name, repo, force)
		if err != nil {
			return err
		}

		ctx.Info("Application '%s' pushed.", a.Name)
	}

	ctx.Succ("Project '%s' pushed to server %s.", project, ctx.CurrentOpenApiHost)
	return nil
}
