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

	"github.com/pkg/errors"

	"github.com/erda-project/erda/tools/cli/command"
	"github.com/erda-project/erda/tools/cli/dicedir"
)

var CONFIGW = command.Command{
	Name:       "use-context",
	ParentName: "CONFIG",
	ShortHelp:  "use context in config file for Erda CLI",
	Example:    "$ erda-cli config use-context <name>",
	Args: []command.Arg{
		command.StringArg{}.Name("name"),
	},
	Run: ConfigOpsWUseCtx,
}

func ConfigOpsWUseCtx(ctx *command.Context, name string) error {
	err := configOpsW("use-context", name, "", "", "")
	if err != nil {
		return err
	}

	ctx.Succ(fmt.Sprintf("Use context \"%s\".", name))
	return nil
}

func configOpsW(ops, name, server, org, platform string) error {
	file, conf, err := command.GetConfig()
	if err != nil && err != dicedir.NotExist {
		return err
	}
	switch ops {
	case "set-platform":
		if server == "" {
			return errors.New("Must set server by --server")
		}
		setPlatform(conf, name, server, org)
	case "set-context":
		if platform == "" {
			return errors.New("Must set platform by --platform")
		}
		setContext(conf, name, platform)
	case "use-context":
		err = useContext(conf, name)
		if err != nil {
			return err
		}
	case "delete-platform":
		deletePlatform(conf, name)
	case "delete-context":
		deleteContext(conf, name)
	default:
		return errors.New(ops + " ops not found")
	}

	err = command.SetConfig(file, conf)
	if err != nil {
		return err
	}

	return nil
}

func setPlatform(conf *command.Config, name, server, org string) {
	notExist := true

	for _, p := range conf.Platforms {
		if p.Name == name {
			p.Server = server
			notExist = false
		}
	}

	if notExist {
		conf.Platforms = append(conf.Platforms, &command.Platform{
			name, server,
			&command.OrgInfo{Name: org},
		})
	}
}

func deletePlatform(conf *command.Config, name string) {
	var ps []*command.Platform
	for _, p := range conf.Platforms {
		if p.Name != name {
			ps = append(ps, p)
		}
	}
	conf.Platforms = ps
}

func setContext(conf *command.Config, name, platform string) {
	notExist := true
	for _, c := range conf.Contexts {
		if c.Name == name {
			c.PlatformName = platform
			notExist = false
		}
	}

	if notExist {
		conf.Contexts = append(conf.Contexts, &command.Ctx{name, platform})
	}
}

func deleteContext(conf *command.Config, name string) {
	var cs []*command.Ctx
	for _, c := range conf.Contexts {
		if c.Name != name {
			cs = append(cs, c)
		}
	}
	conf.Contexts = cs
}

func useContext(conf *command.Config, name string) error {
	for _, c := range conf.Contexts {
		if c.Name == name {
			conf.CurrentContext = name
			return nil
		}
	}

	return errors.New(fmt.Sprintf("context %s not found", name))
}