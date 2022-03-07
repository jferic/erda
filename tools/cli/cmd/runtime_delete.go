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
	"github.com/pkg/errors"

	"github.com/erda-project/erda/tools/cli/command"
	"github.com/erda-project/erda/tools/cli/common"
)

var RUNTIMEDELETE = command.Command{
	Name:       "delete",
	ParentName: "RUNTIME",
	ShortHelp:  "delete runtime",
	Example:    "erda-cli runtime delete",
	Flags: []command.Flag{
		command.Uint64Flag{Short: "", Name: "runtime-id", Doc: "the id of a runtime", DefaultValue: 0},
	},
	Run: DeleteRuntime,
}

func DeleteRuntime(ctx *command.Context, runtimeId uint64) error {
	var org string
	var orgId uint64
	org, orgId, err := getOrgId(ctx, org, orgId)
	if err != nil {
		return err
	}

	if runtimeId <= 0 {
		return errors.New("Invalid runtime id")
	}

	err = common.DeleteRuntime(ctx, orgId, runtimeId)
	if err != nil {
		return err
	}

	ctx.Succ("Runtime deleted.")
	return nil
}
