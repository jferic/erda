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
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/tools/cli/command"
	"github.com/erda-project/erda/tools/cli/format"
	"github.com/erda-project/erda/tools/cli/prettyjson"
)

var APPLICATIONCREATE = command.Command{
	Name:       "create",
	ParentName: "APPLICATION",
	ShortHelp:  "create application",
	Example:    "$ erda-cli application create --project-id=<id> -n <name>",
	Flags: []command.Flag{
		command.Uint64Flag{Short: "", Name: "project-id", Doc: "the id of a project ", DefaultValue: 0},
		command.StringFlag{Short: "n", Name: "name", Doc: "the name of an application ", DefaultValue: ""},
		command.StringFlag{Short: "m", Name: "mode",
			Doc:          "application type, available values：LIBRARY, SERVICE, BIGDATA, PROJECT_SERVICE",
			DefaultValue: "SERVICE"},
		command.StringFlag{"d", "description", "description of the application", ""},
	},
	Run: ApplicationCreate,
}

func ApplicationCreate(ctx *command.Context, projectId uint64, name, mode, desc string) error {
	if name == "" {
		return errors.New("Invalid project name")
	}

	if err := apistructs.ApplicationMode(mode).CheckAppMode(); err != nil {
		return err
	}

	var request apistructs.ApplicationCreateRequest
	var response apistructs.ApplicationCreateResponse
	var b bytes.Buffer

	request.Name = name
	request.Mode = apistructs.ApplicationMode(mode)
	request.Desc = desc
	request.ProjectID = projectId

	resp, err := ctx.Post().Path("/api/applications").JSONBody(request).Do().Body(&b)
	if err != nil {
		return fmt.Errorf(
			format.FormatErrMsg("create", "failed to request ("+err.Error()+")", false))
	}

	if !resp.IsOK() {
		return fmt.Errorf(format.FormatErrMsg("create",
			fmt.Sprintf("failed to request, status-code: %d, content-type: %s, raw bod: %s",
				resp.StatusCode(), resp.ResponseHeader("Content-Type"), b.String()), false))
	}

	if err := json.Unmarshal(b.Bytes(), &response); err != nil {
		return fmt.Errorf(format.FormatErrMsg("create",
			fmt.Sprintf("failed to unmarshal application create response ("+err.Error()+")"), false))
	}

	if !response.Success {
		return fmt.Errorf(format.FormatErrMsg("create",
			fmt.Sprintf("failed to request, error code: %s, error message: %s",
				response.Error.Code, response.Error.Msg), false))
	}

	ctx.Succ("Application created.")

	s, err := prettyjson.Marshal(response.Data)
	if err != nil {
		return fmt.Errorf(format.FormatErrMsg("create",
			"failed to prettyjson marshal application data ("+err.Error()+")", false))
	}

	fmt.Println(string(s))
	return nil
}