package models

type Metadata struct {
	Title       string       `yaml:"title" validate:"required"`
	Version     string       `yaml:"version" validate:"required"`
	Maintainers []Maintainer `validate:"required,dive,required"`
	Company     string       `yaml:"company" validate:"required"`
	Website     string       `yaml:"website" validate:"required,url"`
	Source      string       `yaml:"source" validate:"required,url"`
	License     string       `yaml:"license" validate:"required"`
	Description string       `yaml:"description" validate:"required"`
}
