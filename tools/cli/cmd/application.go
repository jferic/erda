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
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/modules/core-services/types"
	"github.com/erda-project/erda/pkg/terminal/table"
	"github.com/erda-project/erda/tools/cli/command"
	"github.com/erda-project/erda/tools/cli/common"
	"github.com/erda-project/erda/tools/cli/utils"
)

var APPLICATION = command.Command{
	Name:      "application",
	ShortHelp: "list applications",
	Example:   "$ erda-cli application --project=<name>",
	Flags: []command.Flag{
		command.BoolFlag{Short: "", Name: "no-headers", Doc: "if true, don't print headers (default print headers)", DefaultValue: false},
		command.IntFlag{Short: "", Name: "page-size", Doc: "the number of page size", DefaultValue: 10},
		command.BoolFlag{Short: "", Name: "with-owner", Doc: "if true, return owners of applications", DefaultValue: false},
	},
	Run: GetApplications,
}

func GetApplications(ctx *command.Context, noHeaders bool, pageSize int, withOwner bool) error {
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

	num := 0
	err = utils.PagingView(
		func(pageNo, pageSize int) (bool, error) {
			pagingApplication, err := common.GetPagingApplications(ctx, orgId, projectId, pageNo, pageSize)
			if err != nil {
				return false, err
			}

			data := [][]string{}
			for _, p := range pagingApplication.List {
				line := []string{
					strconv.FormatUint(p.ID, 10),
					p.Name,
					p.DisplayName,
				}

				if withOwner {
					var ns []string
					ms, err := common.GetMembers(ctx, apistructs.AppScope, p.ID, []string{types.RoleAppOwner})
					if err != nil {
						return false, err
					}
					for _, m := range ms {
						ns = append(ns, m.Nick)
					}
					line = append(line, strings.Join(ns, ","))
				}

				line = append(line, p.Desc)

				data = append(data, line)
			}

			t := table.NewTable()
			if !noHeaders {
				h := []string{
					"ApplicationID", "Name", "DisplayName",
				}
				if withOwner {
					h = append(h, "Owner")
				}
				h = append(h, "Description")
				t.Header(h)
			}
			err = t.Data(data).Flush()
			if err != nil {
				return false, err
			}

			num += len(pagingApplication.List)
			return pagingApplication.Total > num, nil
		}, "Continue to display applications?", pageSize, command.Interactive)
	if err != nil {
		return err
	}

	return nil
}

func checkApplicationParam(application string, applicationId uint64) {
	if application != "" && applicationId != 0 {
		fmt.Println("Both --application and --application-id are set, we will only use name set by --application")
	}
}

func getApplicationId(ctx *command.Context, orgId, projectId uint64, application string, applicationId uint64) (string, uint64, error) {
	if application != "" {
		// TODO get no projectid
		appId, err := common.GetApplicationIdByName(ctx, orgId, projectId, application)
		if err != nil {
			return application, applicationId, err
		}
		applicationId = appId
	}

	if application == "" && ctx.CurrentApplication.Name == "" {
		return application, applicationId, errors.New("Invalid application name")
	}

	if application == "" && ctx.CurrentApplication.Name != "" {
		application = ctx.CurrentApplication.Name
	}

	if applicationId <= 0 && ctx.CurrentApplication.ID <= 0 && application != "" {
		appId, err := common.GetApplicationIdByName(ctx, orgId, projectId, application)
		if err != nil {
			return application, applicationId, err
		}
		ctx.CurrentApplication.ID = appId
		applicationId = appId
	}

	if applicationId <= 0 && ctx.CurrentApplication.ID <= 0 {
		return application, applicationId, errors.New("Invalid application id")
	}

	if applicationId == 0 && ctx.CurrentApplication.ID > 0 {
		applicationId = ctx.CurrentApplication.ID
	}

	return application, applicationId, nil
}
