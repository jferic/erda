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

package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	pb "github.com/erda-project/erda-proto-go/dop/projectpipeline/pb"
	"github.com/erda-project/erda/tools/cli/command"
	"github.com/erda-project/erda/tools/cli/utils"
)

func CreateProjectPipeline(ctx *command.Context, orgId, projectId, applicationId uint64, name, sourceType,
	ref, path, filename string) (*pb.ProjectPipeline, error) {
	var resp pb.CreateProjectPipelineResponse
	var b bytes.Buffer

	response, err := ctx.Post().Path("/api/project-pipeline").
		Header("Org-ID", strconv.FormatUint(orgId, 10)).
		Param("projectID", strconv.FormatUint(projectId, 10)).
		Param("appID", strconv.FormatUint(applicationId, 10)).
		Param("name", name).Param("sourceType", sourceType).
		Param("ref", ref).Param("path", path).Param("fileName", filename).
		Do().Body(&b)
	if err != nil {
		return nil, fmt.Errorf(
			utils.FormatErrMsg("post", "failed to request ("+err.Error()+")", false))
	}

	if !response.IsOK() {
		return nil, fmt.Errorf(utils.FormatErrMsg("post",
			fmt.Sprintf("failed to request, status-code: %d, content-type: %s, raw bod: %s",
				response.StatusCode(), response.ResponseHeader("Content-Type"), b.String()), false))
	}

	if err := json.Unmarshal(b.Bytes(), &resp); err != nil {
		return nil, fmt.Errorf(utils.FormatErrMsg("list",
			fmt.Sprintf("failed to unmarshal projects post response ("+err.Error()+")"), false))
	}

	return resp.ProjectPipeline, nil
}

func GetProjectPipelines(ctx *command.Context, orgId, projectId uint64) ([]*pb.PipelineYmlList, error) {
	var resp pb.ListAppPipelineYmlResponse
	var b bytes.Buffer

	response, err := ctx.Get().Path("/api/project-pipeline/actions/get-pipeline-yml-list").
		Param("joined", "true").
		Param("appId", strconv.FormatUint(orgId, 10)).
		Do().Body(&b)
	if err != nil {
		return nil, fmt.Errorf(
			utils.FormatErrMsg("list", "failed to request ("+err.Error()+")", false))
	}

	if !response.IsOK() {
		return nil, fmt.Errorf(utils.FormatErrMsg("list",
			fmt.Sprintf("failed to request, status-code: %d, content-type: %s, raw bod: %s",
				response.StatusCode(), response.ResponseHeader("Content-Type"), b.String()), false))
	}

	if err := json.Unmarshal(b.Bytes(), &resp); err != nil {
		return nil, fmt.Errorf(utils.FormatErrMsg("list",
			fmt.Sprintf("failed to unmarshal projects list response ("+err.Error()+")"), false))
	}

	return resp.Result, nil
}
