package pkg

import "gopkg.in/yaml.v2"

/**
deployments:
  projects:
    - name: simpleserver
      seelf_key: 2mCQJBS8MKviu360tn76rypG98M
      github_url: "https://github.com/nanikjava/simpleserver"
**/

type Deployment struct {
	Deployments struct {
		Projects []Project `yaml:"projects"`
	} `yaml:"deployments"`
}

type Project struct {
	Name  string `yaml:"name"`
	SeelfKey  string `yaml:"seelf_key"`
	GithubURL string `yaml:"github_url"`
}


func ParseConfigs(yamlData []byte) (*Deployment, error) {
	var result Deployment

	// Unmarshal the YAML into Go structures
	err := yaml.Unmarshal(yamlData, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

