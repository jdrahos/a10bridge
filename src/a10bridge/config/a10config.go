package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

//A10Config configuration of a10 instances to work with
type A10Config struct {
	Instances []A10Instance `yaml:"instances"`
}

type A10Instances []A10Instance

func (s A10Instances) Len() int {
	return len(s)
}
func (s A10Instances) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s A10Instances) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

//A10Instance a10 instance configuration
type A10Instance struct {
	Name       string `yaml:"name"`
	APIUrl     string `yaml:"apiUrl"`
	APIVersion int    `yaml:"apiVersion"`
	UserName   string `yaml:"userName"`
	Password   string `yaml:"password"`
}

func readA10Configuration(configFilePath string) (*A10Config, error) {
	var a10config A10Config
	fileContent, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(fileContent, &a10config)
	if err != nil {
		return nil, err
	}
	return &a10config, nil
}
