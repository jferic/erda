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

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/tools/cli/command"
	"github.com/erda-project/erda/tools/cli/utils"
)

type ProjectPipeline struct {
	ID               string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Name             string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Creator          string `protobuf:"bytes,3,opt,name=creator,proto3" json:"creator,omitempty"`
	Category         string `protobuf:"bytes,4,opt,name=category,proto3" json:"category,omitempty"`
	SourceType       string `protobuf:"bytes,7,opt,name=sourceType,proto3" json:"sourceType,omitempty"`
	Remote           string `protobuf:"bytes,8,opt,name=remote,proto3" json:"remote,omitempty"`
	Ref              string `protobuf:"bytes,9,opt,name=ref,proto3" json:"ref,omitempty"`
	Path             string `protobuf:"bytes,10,opt,name=path,proto3" json:"path,omitempty"`
	FileName         string `protobuf:"bytes,11,opt,name=fileName,proto3" json:"fileName,omitempty"`
	PipelineSourceId string `protobuf:"bytes,12,opt,name=pipelineSourceId,proto3" json:"pipelineSourceId,omitempty"`
}

type CreateProjectPipelineResponse struct {
	apistructs.Header
	Data CreateProjectPipelineResponseData `json:"data"`
}

type CreateProjectPipelineResponseData struct {
	ProjectPipeline *ProjectPipeline `protobuf:"bytes,1,opt,name=ProjectPipeline,proto3" json:"ProjectPipeline,omitempty"`
}

func CreateProjectPipeline(ctx *command.Context, orgId, projectId, applicationId uint64, name, sourceType,
	ref, path, filename string) (*ProjectPipeline, error) {
	var resp CreateProjectPipelineResponse
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
		return nil, fmt.Errorf(utils.FormatErrMsg("post",
			fmt.Sprintf("failed to unmarshal projects post response ("+err.Error()+")"), false))
	}

	if !resp.Success {
		return nil, fmt.Errorf(utils.FormatErrMsg("post project pipeline",
			fmt.Sprintf("failed to request, error code: %s, error message: %s",
				resp.Error.Code, resp.Error.Msg), false))
	}

	fmt.Printf("%v", resp)

	return resp.Data.ProjectPipeline, nil
}

//func GetProjectPipelines(ctx *command.Context, orgId, projectId uint64) ([]*pb.PipelineYmlList, error) {
//	var resp pb.ListAppPipelineYmlResponse
//	var b bytes.Buffer
//
//	response, err := ctx.Get().Path("/api/project-pipeline/actions/get-pipeline-yml-list").
//		Param("joined", "true").
//		Param("appId", strconv.FormatUint(orgId, 10)).
//		Do().Body(&b)
//	if err != nil {
//		return nil, fmt.Errorf(
//			utils.FormatErrMsg("list", "failed to request ("+err.Error()+")", false))
//	}
//
//	if !response.IsOK() {
//		return nil, fmt.Errorf(utils.FormatErrMsg("list",
//			fmt.Sprintf("failed to request, status-code: %d, content-type: %s, raw bod: %s",
//				response.StatusCode(), response.ResponseHeader("Content-Type"), b.String()), false))
//	}
//
//	if err := json.Unmarshal(b.Bytes(), &resp); err != nil {
//		return nil, fmt.Errorf(utils.FormatErrMsg("list",
//			fmt.Sprintf("failed to unmarshal projects list response ("+err.Error()+")"), false))
//	}
//
//	return resp.Result, nil
//}
