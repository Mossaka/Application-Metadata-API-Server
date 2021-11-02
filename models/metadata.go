package models

type Metadata struct {
	Title       string       `yaml:"title" validate:"required,max=35"`
	Version     string       `yaml:"version" validate:"required,max=25"`
	Maintainers []Maintainer `validate:"required,dive,required"`
	Company     string       `yaml:"company" validate:"required,max=35"`
	Website     string       `yaml:"website" validate:"required,url,max=2048"`
	Source      string       `yaml:"source" validate:"required,url,max=2048"`
	License     string       `yaml:"license" validate:"required,max=50"`
	Description string       `yaml:"description" validate:"required,max=2048"`
}
