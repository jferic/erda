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
	"github.com/erda-project/erda/tools/cli/dicedir"
	"github.com/erda-project/erda/tools/cli/format"
)

//func GetOrgDetail(ctx *command.Context, orgIdorName string) (apistructs.OrgFetchResponse, error) {
//	var resp apistructs.OrgFetchResponse
//	var b bytes.Buffer
//
//	if orgIdorName == "" {
//		return apistructs.OrgFetchResponse{}, fmt.Errorf(format.FormatErrMsg("get organization detail",
//			"invalid required parameter organization", false))
//	}
//
//	response, err := ctx.Get().Path(fmt.Sprintf("/api/orgs/%s", orgIdorName)).Do().Body(&b)
//	if err != nil {
//		return apistructs.OrgFetchResponse{}, fmt.Errorf(format.FormatErrMsg(
//			"get organization detail", "failed to request ("+err.Error()+")", false))
//	}
//
//	if !response.IsOK() {
//		return apistructs.OrgFetchResponse{}, fmt.Errorf(format.FormatErrMsg("get organization detail",
//			fmt.Sprintf("failed to request, status-code: %d, content-type: %s, raw bod: %s",
//				response.StatusCode(), response.ResponseHeader("Content-Type"), b.String()), false))
//	}
//
//	if err := json.Unmarshal(b.Bytes(), &resp); err != nil {
//		return apistructs.OrgFetchResponse{}, fmt.Errorf(format.FormatErrMsg("get organization detail",
//			fmt.Sprintf("failed to unmarshal organization detail response ("+err.Error()+")"), false))
//	}
//
//	if !resp.Success {
//		return apistructs.OrgFetchResponse{}, fmt.Errorf(format.FormatErrMsg("get organization detail",
//			fmt.Sprintf("failed to request, error code: %s, error message: %s",
//				resp.Error.Code, resp.Error.Msg), false))
//	}
//
//	return resp, nil
//}

func GetMembers(ctx *command.Context, scopeType apistructs.ScopeType, scopeId uint64, roles []string) ([]apistructs.Member, error) {
	var orgs []apistructs.Member
	err := dicedir.PagingAll(func(pageNo, pageSize int) (bool, error) {

		page, err := GetPagingMembers(ctx, scopeType, scopeId, roles, pageNo, pageSize)
		if err != nil {
			return false, err
		}
		orgs = append(orgs, page.List...)

		if page.Total > len(orgs) {
			return true, nil
		}
		return false, nil
	}, 20)
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func GetPagingMembers(ctx *command.Context, scopeType apistructs.ScopeType, scopeId uint64, roles []string, pageNo, pageSize int) (apistructs.MemberList, error) {
	var resp apistructs.MemberListResponse
	var b bytes.Buffer

	req := ctx.Get().Path("/api/members").
		Param("scopeId", strconv.FormatUint(scopeId, 10)).
		Param("scopeType", string(scopeType)).
		Param("pageNo", strconv.Itoa(pageNo)).
		Param("pageSize", strconv.Itoa(pageSize))
	for _, r := range roles {
		req.Param("roles", r)
	}
	response, err := req.Do().Body(&b)
	if err != nil {
		return apistructs.MemberList{}, fmt.Errorf(
			format.FormatErrMsg("members", "failed to request ("+err.Error()+")", false))
	}

	if !response.IsOK() {
		return apistructs.MemberList{}, fmt.Errorf(format.FormatErrMsg("members",
			fmt.Sprintf("failed to request, status-code: %d, content-type: %s, raw bod: %s",
				response.StatusCode(), response.ResponseHeader("Content-Type"), b.String()), false))
	}

	if err := json.Unmarshal(b.Bytes(), &resp); err != nil {
		return apistructs.MemberList{}, fmt.Errorf(format.FormatErrMsg("members",
			fmt.Sprintf("failed to unmarshal organizations list response ("+err.Error()+")"), false))
	}

	if !resp.Success {
		return apistructs.MemberList{}, fmt.Errorf(format.FormatErrMsg("members",
			fmt.Sprintf("error code(%s), error message(%s)", resp.Error.Code, resp.Error.Msg), false))
	}

	return resp.Data, nil
}