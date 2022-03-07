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

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/terminal/table"
	"github.com/erda-project/erda/tools/cli/command"
	"github.com/erda-project/erda/tools/cli/common"
	"github.com/erda-project/erda/tools/cli/utils"
)

var APPLICATIONMEMBER = command.Command{
	Name:       "member",
	ParentName: "APPLICATION",
	ShortHelp:  "display members of the application",
	Example:    "$ erda-cli application member --application=<name>",
	Flags: []command.Flag{
		command.BoolFlag{Short: "", Name: "no-headers", Doc: "if true, don't print headers (default print headers)", DefaultValue: false},
		command.StringFlag{Short: "", Name: "application", Doc: "the name of an application ", DefaultValue: ""},
		command.IntFlag{Short: "", Name: "page-size", Doc: "the number of page size", DefaultValue: 10},
		command.StringListFlag{Short: "", Name: "roles", Doc: "roles to list", DefaultValue: nil},
	},
	Run: ApplicationMember,
}

func ApplicationMember(ctx *command.Context, noHeaders bool, application string, pageSize int, roles []string) error {
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

	num := 0
	utils.PagingView(func(pageNo, pageSize int) (bool, error) {
		pagingMembers, err := common.GetPagingMembers(ctx, apistructs.AppScope, applicationId, roles, pageNo, pageSize)
		if err != nil {
			fmt.Println(err)
			return false, err
		}

		data := [][]string{}
		for _, m := range pagingMembers.List {
			data = append(data, []string{
				m.Nick,
				m.Name,
				m.Email,
				m.Mobile,
				strings.Join(m.Roles, ","),
			})
		}

		t := table.NewTable()
		if !noHeaders {
			t.Header([]string{
				"Nick", "Name", "Email", "Mobile", "Roles",
			})
		}
		err = t.Data(data).Flush()
		if err != nil {
			return false, err
		}
		num += len(pagingMembers.List)
		return pagingMembers.Total > num, nil
	}, "Continue to display members?", pageSize, command.Interactive)

	return nil
}
