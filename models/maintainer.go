package models

type Maintainer struct {
	Name  string `yaml:"name" validate:"required,max=35"`
	Email string `yaml:"email" validate:"required,email,max=50"`
}
