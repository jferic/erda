package command

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/erda-project/erda/tools/cli/utils"
)

var ConfigVersion string = "v0.0.1"

type ProjectInfo struct {
	Version      string            `yaml:"version"`
	Server       string            `yaml:"server"`
	Org          string            `yaml:"org"`
	OrgId        uint64            `yaml:"org_id"`
	Project      string            `yaml:"project"`
	ProjectId    uint64            `yaml:"project_id"`
	Applications []ApplicationInfo `yaml:"applications"`
}

type ApplicationInfo struct {
	Application   string `yaml:"application"`
	ApplicationId uint64 `yaml:"application_id"`
}

func GetProjectConfig() (string, *ProjectInfo, error) {
	info := ProjectInfo{Version: ConfigVersion}
	config, err := utils.FindProjectConfig()
	if err != nil {
		return config, &info, err
	}

	f, err := os.Open(config)
	if err != nil {
		return config, &info, err
	}
	if err := yaml.NewDecoder(f).Decode(&info); err != nil {
		os.Remove(config)
		return config, &info, err
	}

	return config, &info, nil
}

func SetProjectConfig(file string, conf *ProjectInfo) error {
	c, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, c, 0655)
	if err != nil {
		return err
	}

	return nil
}
