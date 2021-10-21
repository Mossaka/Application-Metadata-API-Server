package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type maintainer struct {
	Name  string `yaml:"name" validate:"required"`
	Email string `yaml:"email" validate:"required,email"`
}

type metadata struct {
	Title       string       `yaml:"title" validate:"required"`
	Version     string       `yaml:"version" validate:"required"`
	Maintainers []maintainer `validate:"required,dive,required"`
	Company     string       `yaml:"company" validate:"required"`
	Website     string       `yaml:"website" validate:"required,url"`
	Source      string       `yaml:"source" validate:"required,url"`
	License     string       `yaml:"license" validate:"required"`
	Description string       `yaml:"description" validate:"required"`
}

var metadata_list = []metadata{
	{
		Title:   "Valid App 2",
		Version: "1.0.1",
		Maintainers: []maintainer{
			{
				Name:  "AppTwo Maintainer",
				Email: "apptwo@hotmail.com",
			},
		},
		Company:     "Random Inc.",
		Website:     "https://website.com",
		Source:      "https://github.com/upbound/repo",
		License:     "Apache-2.0",
		Description: "### Why app 2 is the best\nBecause it simply is...",
	},
}

func get_metadata(c *gin.Context) {
	c.YAML(http.StatusOK, metadata_list)
}

func post_metadata(c *gin.Context) {
	var new_metadata metadata
	if err := c.BindYAML(&new_metadata); err != nil {
		c.YAML(http.StatusBadRequest, gin.H{"UserError": err.Error()})
		return
	}
	validate := validator.New()
	if err := validate.Struct(&new_metadata); err != nil {
		c.YAML(http.StatusBadRequest, gin.H{"UserError": err.Error()})
		return
	}
	metadata_list = append(metadata_list, new_metadata)
	c.YAML(http.StatusOK, new_metadata)
}

func main() {
	router := gin.Default()
	router.GET("/api/metadata", get_metadata)
	router.POST("/api/metadata", post_metadata)

	router.Run("localhost:8080")
}
