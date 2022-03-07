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

	"github.com/erda-project/erda/tools/cli/common"

	"github.com/pkg/errors"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/tools/cli/command"
	"github.com/erda-project/erda/tools/cli/utils"
)

var APPLICATIONCREATE = command.Command{
	Name:       "create",
	ParentName: "APPLICATION",
	ShortHelp:  "create application",
	Example:    "$ erda-cli application create -n <name>",
	Flags: []command.Flag{
		command.StringFlag{Short: "n", Name: "name", Doc: "the name of an application ", DefaultValue: ""},
		command.StringFlag{Short: "m", Name: "mode",
			Doc:          "application type, available values：LIBRARY, SERVICE, BIGDATA, PROJECT_SERVICE",
			DefaultValue: "SERVICE"},
		command.StringFlag{Short: "d", Name: "description", Doc: "description of the application", DefaultValue: ""},
		command.StringFlag{Short: "s", Name: "sonarhost", Doc: "host url of sonarqube", DefaultValue: ""},
		command.StringFlag{Short: "t", Name: "sonartoken", Doc: "token of project in sonarqube", DefaultValue: ""},
		command.StringFlag{Short: "k", Name: "sonarproject", Doc: "project key in sonarqube", DefaultValue: ""},
	},
	Run: ApplicationCreate,
}

func ApplicationCreate(ctx *command.Context, name, mode, desc, sonarhost, sonartoken, sonarproject string) error {
	if name == "" {
		return errors.New("Invalid project name")
	}

	if err := apistructs.ApplicationMode(mode).CheckAppMode(); err != nil {
		return err
	}

	config, pInfo, err := command.GetProjectConfig()
	if err == utils.NotExist {
		return errors.New("Not in a project directory")
	} else if err != nil {
		return err
	}

	var project string
	var projectId uint64
	project, projectId, err = getProjectId(ctx, projectId, project, projectId)
	if err != nil {
		return err
	}

	app, err := common.CreateApplication(ctx, projectId, name, mode, desc, sonarhost, sonartoken, sonarproject)
	if err != nil {
		return err
	}

	var (
		sonarHost    string
		sonarToken   string
		sonarProject string
	)
	if app.SonarConfig != nil {
		sonarHost = app.SonarConfig.Host
		sonarToken = app.SonarConfig.Token
		sonarProject = app.SonarConfig.ProjectKey
	}
	// TODO clone project into pwd
	pInfo.Applications = append(pInfo.Applications, command.ApplicationInfo{
		app.Name, app.ID, app.Mode, app.Desc,
		sonarHost, sonarToken, sonarProject,
	})

	err = command.SetProjectConfig(config, pInfo)
	if err != nil {
		return err
	}

	ctx.Succ("Application created.")

	s, err := utils.Marshal(app)
	if err != nil {
		return fmt.Errorf(utils.FormatErrMsg("create",
			"failed to prettyjson marshal application data ("+err.Error()+")", false))
	}

	fmt.Println(string(s))

	// TODO print init gittar ..

	return nil
}
