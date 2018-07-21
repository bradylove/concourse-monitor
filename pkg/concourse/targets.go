package concourse

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type FlyRC struct {
	Targets Targets
}

type Token struct {
	Type  string
	Value string
}

type Targets map[string]Target

type Target struct {
	API   string
	Team  string
	Name  string
	Token Token
}

func LoadTargets(path string) ([]Target, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var flyRC FlyRC
	if err := yaml.Unmarshal(data, &flyRC); err != nil {
		return nil, err
	}

	targets := make([]Target, 0)
	for k, t := range flyRC.Targets {
		t.Name = k
		targets = append(targets, t)
	}

	return targets, nil
}
