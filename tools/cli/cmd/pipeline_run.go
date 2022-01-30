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
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pkg/errors"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/tools/cli/command"
	"github.com/erda-project/erda/tools/cli/common"
	"github.com/erda-project/erda/tools/cli/utils"
)

// BUILD command
var PIPELINERUN = command.Command{
	Name:       "run",
	ParentName: "PIPELINE",
	ShortHelp:  "create a pipeline and run it",
	Example:    "$ erda-cli pipeline run <path-to/pipeline.yml>",
	Args: []command.Arg{
		command.StringArg{}.Name("filename"),
	},
	Flags: []command.Flag{
		command.StringFlag{Short: "", Name: "branch", Doc: "branch to create pipeline, default is current branch", DefaultValue: ""},
		//command.StringFlag{Short: "f", Name: "filename", Doc: "filename for 'pipeline.yml'", DefaultValue: path.Join(utils.ProjectPipelineDir, "pipeline.yml")},
		command.BoolFlag{Short: "w", Name: "watch", Doc: "watch the status", DefaultValue: false},
	},
	ValidArgsFunction:          FilenameCompletion,
	RegisterFlagCompletionFunc: map[string]interface{}{"branch": BranchCompletion},
	Run:                        PipelineRun,
}

func FilenameCompletion(ctx *cobra.Command, args []string, toComplete string, filename, branch string, watch bool) []string {
	comps := []string{}
	if branch != "" {
		b, err := utils.GetWorkspaceBranch()
		if err != nil || branch != b {
			return comps
		}
	}

	p, err := getWorkspacePipelines()
	if err == nil {
		comps = p
	}
	return comps
}

func BranchCompletion(ctx *cobra.Command, args []string, toComplete string, filename, branch string, watch bool) []string {
	comps := []string{}

	c1 := exec.Command("git", "branch")
	c2 := exec.Command("cut", "-c", "3-")
	output, err := utils.PipeCmds(c1, c2)
	if err == nil {
		splites := strings.Split(output, "\n")
		for _, s := range splites {
			comps = append(comps, s)
		}
	}
	return comps
}

// Create an pipeline and run it
func PipelineRun(ctx *command.Context, filename, branch string, watch bool) error {
	// 1. check if .git dir exists in current directory
	// 2. parse current branch
	// 3. create pipeline, run it
	if _, err := os.Stat(".git"); err != nil {
		return errors.New("Current directory is not a local git repository")
	}

	dirty, err := utils.IsWorkspaceDirty()
	if err != nil {
		return err
	}
	if dirty {
		return errors.New("Changes should be committed first")
	}

	if branch == "" {
		b, err := utils.GetWorkspaceBranch()
		if err != nil {
			return err
		}
		branch = b
	}

	// fetch appID
	info, err := utils.GetWorkspaceInfo(command.Remote)
	if err != nil {
		return err
	}

	org, err := common.GetOrgDetail(ctx, info.Org)
	if err != nil {
		return err
	}

	repoStats, err := common.GetRepoStats(ctx, org.ID, info.Project, info.Application)
	if err != nil {
		return err
	}

	var (
		request      apistructs.PipelineCreateRequest
		pipelineResp apistructs.PipelineCreateResponse
	)
	request.AppID = uint64(repoStats.ApplicationID)
	request.Branch = branch
	request.Source = apistructs.PipelineSourceDice
	request.PipelineYmlSource = apistructs.PipelineYmlSourceGittar
	request.PipelineYmlName = filename
	request.AutoRun = true

	// create pipeline
	response, err := ctx.Post().Path("/api/cicds").JSONBody(request).Do().JSON(&pipelineResp)
	if err != nil {
		return err
	}
	if !response.IsOK() {
		return errors.Errorf("build fail, status code: %d, err: %+v", response.StatusCode(), pipelineResp.Error)
	}
	if !pipelineResp.Success {
		return errors.Errorf("build fail: %+v", pipelineResp.Error)
	}

	if watch {
		err = PipelineView(ctx, branch, pipelineResp.Data.ID, true)
		if err != nil {
			ctx.Fail("failed to watch status of pipeline %d", pipelineResp.Data.ID)
		}
	} else {
		ctx.Succ("run pipeline: %s for branch: %s, pipelineID: %d, you can view building status via `erda-cli pipeline view -i %d`",
			filename, branch, pipelineResp.Data.ID, pipelineResp.Data.ID)
	}

	return nil
}
