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

	"github.com/erda-project/erda/tools/cli/command"
	"github.com/erda-project/erda/tools/cli/common"
	"github.com/erda-project/erda/tools/cli/format"
	"github.com/erda-project/erda/tools/cli/prettyjson"
)

var ADDONINSPECT = command.Command{
	Name:       "inspect",
	ParentName: "ADDON",
	ShortHelp:  "inspect addon",
	Example:    "$ erda-cli addon inspect --addon-id=<id>",
	Flags: []command.Flag{
		command.Uint64Flag{Short: "", Name: "org-id", Doc: "the id of an organization", DefaultValue: 0},
		command.StringFlag{Short: "", Name: "org", Doc: "the name of an organization", DefaultValue: ""},
		command.StringFlag{Short: "", Name: "addon-id", Doc: "the id of an addon", DefaultValue: ""},
	},
	Run: AddonInspect,
}

func AddonInspect(ctx *command.Context, orgId uint64, org, addonId string) error {
	checkOrgParam(org, orgId)
	orgId, err := getOrgId(ctx, org, orgId)
	if err != nil {
		return err
	}

	resp, err := common.GetAddonDetail(ctx, orgId, addonId)
	if err != nil {
		return err
	}

	s, err := prettyjson.Marshal(resp)
	if err != nil {
		return fmt.Errorf(format.FormatErrMsg("application inspect",
			"failed to prettyjson marshal application data ("+err.Error()+")", false))
	}
	fmt.Println(string(s))

	return nil
}