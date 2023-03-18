package pipelines

type PipelineDef struct {
	Tasks []Step `yaml:"tasks"`
}

type Step struct {
	Image  string `yaml:"image"`
	Script string `yaml:"script"`
	Name   string `yaml:"name"`
}
