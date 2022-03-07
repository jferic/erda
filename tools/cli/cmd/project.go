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
	"github.com/erda-project/erda/modules/dop/types"
	"github.com/erda-project/erda/pkg/terminal/table"
	"github.com/erda-project/erda/tools/cli/command"
	"github.com/erda-project/erda/tools/cli/common"
	"github.com/erda-project/erda/tools/cli/utils"
)

var PROJECT = command.Command{
	Name:      "project",
	ShortHelp: "list projects",
	Example:   "$ erda-cli project",
	Flags: []command.Flag{
		command.BoolFlag{Short: "", Name: "no-headers", Doc: "if true, don't print headers (default print headers)", DefaultValue: false},
		command.IntFlag{Short: "", Name: "page-size", Doc: "the number of page size", DefaultValue: 10},
		command.BoolFlag{Short: "", Name: "with-owner", Doc: "if true, return owners of projects", DefaultValue: false},
	},
	Run: GetProjects,
}

func GetProjects(ctx *command.Context, noHeaders bool, pageSize int, withOwner bool) error {
	var org string
	var orgId uint64

	org, orgId, err := getOrgId(ctx, org, orgId)
	if err != nil {
		return err
	}

	num := 0
	err = utils.PagingView(func(pageNo, pageSize int) (bool, error) {
		pagingProject, err := common.GetPagingProjects(ctx, orgId, pageNo, pageSize)
		if err != nil {
			return false, err
		}

		data := [][]string{}
		for _, p := range pagingProject.List {
			current := " "
			if p.Name == ctx.CurrentProject.Name {
				current = "*"
			}
			line := []string{
				current,
				strconv.FormatUint(p.ID, 10),
				p.Name,
				p.DisplayName,
			}

			if withOwner {
				var ns []string
				ms, err := common.GetMembers(ctx, apistructs.ProjectScope, p.ID, []string{types.RoleProjectOwner})
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
			headers := []string{
				"Current", "ProjectID", "Name", "DisplayName",
			}
			if withOwner {
				headers = append(headers, "Owner")
			}
			headers = append(headers, "Description")

			t.Header(headers)
		}
		err = t.Data(data).Flush()
		if err != nil {
			return false, err
		}

		num += len(pagingProject.List)
		return pagingProject.Total > num, nil
	}, "Continue to display project?", pageSize, command.Interactive)
	if err != nil {
		return err
	}

	return nil
}

func getProjectId(ctx *command.Context, orgId uint64, project string, projectId uint64) (string, uint64, error) {
	if project != "" {
		pId, err := common.GetProjectIdByName(ctx, orgId, project)
		if err != nil {
			return project, projectId, err
		}
		projectId = pId
	}

	if project == "" && ctx.CurrentProject.Name == "" {
		return project, projectId, errors.New("Invalid project name")
	}

	if project == "" && ctx.CurrentProject.Name != "" {
		project = ctx.CurrentProject.Name
	}

	if projectId <= 0 && ctx.CurrentProject.ID <= 0 && project != "" {
		pId, err := common.GetProjectIdByName(ctx, orgId, project)
		if err != nil {
			return project, projectId, err
		}
		ctx.CurrentProject.ID = pId
		projectId = pId
	}

	if projectId <= 0 && ctx.CurrentProject.ID <= 0 {
		return project, projectId, errors.New("Invalid project id")
	}

	if projectId == 0 && ctx.CurrentProject.ID > 0 {
		projectId = ctx.CurrentProject.ID
	}

	return project, projectId, nil
}

func checkProjectParam(project string, projectId uint64) {
	if project != "" && projectId != 0 {
		fmt.Println("Both --project and --project-id are set, we will only use name set by --project")
	}
}
