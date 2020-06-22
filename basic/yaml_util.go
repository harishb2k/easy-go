package basic

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// Example Usage:
//
// type localConfig struct {
//    EmissaryConfiguration Configuration `yaml:"emissaryConfiguration"`
// }
//
// config := localConfig{}
// basic.ReadYaml("app.yml", &config)
//

func ReadYaml(file string, object interface{}) (err error) {
	data, err := ioutil.ReadFile(file)
	if err == nil {
		err = yaml.Unmarshal(data, object)
		return err
	} else {
		return err
	}
}
