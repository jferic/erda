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
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/erda-project/erda/tools/cli/utils"
	"github.com/mholt/archiver"

	"github.com/erda-project/erda/pkg/loop"

	"gopkg.in/yaml.v3"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/tools/cli/command"
	"github.com/erda-project/erda/tools/cli/common"
	"github.com/pkg/errors"
)

var PROJECTCREATE = command.Command{
	Name:       "create",
	ParentName: "PROJECT",
	ShortHelp:  "create project",
	Example:    "erda-cli project create --name=<name>",
	Flags: []command.Flag{
		command.StringFlag{Short: "", Name: "org", Doc: "the name of an organization", DefaultValue: ""},
		command.StringFlag{Short: "n", Name: "name", Doc: "the name of an application ", DefaultValue: ""},
		command.StringFlag{Short: "d", Name: "description", Doc: "description of the project", DefaultValue: ""},
		command.StringFlag{Short: "", Name: "init-package", Doc: "package for init the project", DefaultValue: ""},
		command.Uint64Flag{Short: "", Name: "wait-import", Doc: "minutes to wait package to be import", DefaultValue: 1},
	},
	Run: ProjectCreate,
}

func ProjectCreate(ctx *command.Context, org, project, desc, pkg string, waitImport uint64) error {
	var orgId uint64
	org, orgId, err := getOrgId(ctx, org, orgId)
	if err != nil {
		return err
	}

	var values map[string]interface{}
	if pkg != "" {
		s, err := os.Stat(pkg)
		if err != nil {
			return errors.Errorf("Invalid package %v", err)
		}

		if s.IsDir() {
			files, err := utils.ListDir(pkg)
			if err != nil {
				return err
			}
			zipTmpFile, err := ioutil.TempFile("", "project-package-*.zip")
			if err != nil {
				return err
			}
			defer os.Remove(zipTmpFile.Name())
			defer zipTmpFile.Close()
			err = archiver.Zip.Write(zipTmpFile, files)
			if err != nil {
				return err
			}
			pkg = zipTmpFile.Name()
		}
		if strings.HasSuffix(pkg, ".zip") {
			values, err = readValues(pkg)
			if err != nil {
				return errors.Errorf("Invalid package %v", err)
			}
		} else {
			return errors.Errorf("Invalid package %v, neither a dirctory nor a zip file", err)
		}
	}

	var resourceConfigs *apistructs.ResourceConfigs
	if values != nil {
		resourceConfigs = apistructs.NewResourceConfigs()
		err = parseResources(resourceConfigs, values)
		if err != nil {
			return err
		}
	}

	ctx.Info("Devops project %s creating...", project)
	projectId, err := common.CreateProject(ctx, orgId, project, desc, resourceConfigs)
	if err != nil {
		return err
	}
	ctx.Info("Devops project %s created.", project)

	ctx.Info("Msp project %s creating...", project)
	_, err = common.CreateMSPProject(ctx, projectId, project)
	if err != nil {
		return err
	}
	ctx.Info("Msp project %s created.", project)

	if pkg != "" {
		ctx.Info("Project package importing...")
		fileId, err := common.ImportPackage(ctx, orgId, projectId, pkg)
		if err != nil {
			return errors.Errorf("Import package %s failed %v", pkg, err)
		}

		loop := loop.New(loop.WithMaxTimes(6*waitImport), loop.WithInterval(10*time.Second))
		err = loop.Do(func() (bool, error) {
			record, err := common.GetRecord(ctx, orgId, fileId)
			if err != nil {
				return false, err
			}
			if record.State == apistructs.FileRecordStateFail {
				return true, errors.Errorf("Import package %s failed, error %s", pkg, record.ErrorInfo)
			}
			if record.State == apistructs.FileRecordStateSuccess {
				return true, nil
			}

			return false, nil
		})
		if err != nil {
			return err
		}
		record, err := common.GetRecord(ctx, orgId, fileId)
		if err != nil {
			return err
		}
		if record.State == apistructs.FileRecordStatePending ||
			record.State == apistructs.FileRecordStateProcessing {
			return errors.Errorf("Import package %s timeout.", pkg)
		}
		ctx.Info("Project package imported.")
	}

	ctx.Succ("Project '%s' created.", project)
	return nil
}

func parseResources(resourceConfigs *apistructs.ResourceConfigs, values map[string]interface{}) error {
	for k, v := range values {
		if v == "" {
			return errors.Errorf("Invalid package, found value of '%s' not configed", k)
		}

		splits := strings.SplitN(k, ".", 4)

		if len(splits) == 4 && splits[0] == "values" && splits[2] == "cluster" {
			env := splits[1]
			if splits[3] == "name" {
				resourceConfigs.GetClusterConfig(env).ClusterName = fmt.Sprintf("%v", v)
			} else if splits[3] == "quota.cpuQuota" {
				cpuQuotaStr := fmt.Sprintf("%v", v)
				cpuQuota, err := strconv.ParseFloat(cpuQuotaStr, 64)
				if err != nil {
					return errors.Errorf("Invalid package, found value of '%s' not a float", k)
				}
				resourceConfigs.GetClusterConfig(env).CPUQuota = cpuQuota
			} else if splits[3] == "quota.memoryQuota" {
				memoryQuotaStr := fmt.Sprintf("%v", v)
				memoryQuota, err := strconv.ParseFloat(memoryQuotaStr, 64)
				if err != nil {
					return errors.Errorf("Invalid package, found value of '%s' not a float", k)
				}
				resourceConfigs.GetClusterConfig(env).MemQuota = memoryQuota
			}
		}
	}
	return nil
}

func readValues(pkg string) (map[string]interface{}, error) {
	zipReader, err := zip.OpenReader(pkg)
	if err != nil {
		return nil, err
	}
	valueYml, err := zipReader.Open("values.yml")
	if err != nil {
		return nil, err
	}
	defer valueYml.Close()

	yamlBytes, err := io.ReadAll(valueYml)
	if err != nil {
		return nil, err
	}

	values := map[string]interface{}{}
	if err := yaml.Unmarshal(yamlBytes, values); err != nil {
		return nil, err
	}

	return values, nil
}
