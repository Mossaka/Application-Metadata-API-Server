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
		Company:     "Upbound Inc.",
		Website:     "https://upbound.io",
		Source:      "https://github.com/upbound/repo",
		License:     "Apache-2.0",
		Description: "### Why app 2 is the best\nBecause it simply is...",
	},
	{
		Title:   "Valid App 1",
		Version: "0.0.1",
		Maintainers: []maintainer{
			{
				Name:  "firstmaintainer app1",
				Email: "apptwo@hotmail.com",
			},
			{
				Name:  "secondmaintainer app1",
				Email: "secondmaintainer@gmail.com",
			},
		},
		Company:     "Random Inc.",
		Website:     "https://website.com",
		Source:      "https://github.com/random/repo",
		License:     "Apache-1.0",
		Description: "### Why app 2 is the best\nBecause it simply is...",
	},
}

func get_metadata(c *gin.Context) {
	title := c.Query("title")
	version := c.Query("version")
	company := c.Query("company")
	website := c.Query("website")
	source := c.Query("source")
	license := c.Query("license")
	maintainer_name := c.Query("maintainer_name")
	maintainer_email := c.Query("maintainer_email")

	var metadata_list_filtered []metadata
	for _, metadata := range metadata_list {
		if title != "" && metadata.Title != title {
			continue
		}
		if version != "" && metadata.Version != version {
			continue
		}
		if company != "" && metadata.Company != company {
			continue
		}
		if website != "" && metadata.Website != website {
			continue
		}
		if source != "" && metadata.Source != source {
			continue
		}
		if license != "" && metadata.License != license {
			continue
		}
		for _, maintainer := range metadata.Maintainers {
			if maintainer_name != "" && maintainer.Name != maintainer_name {
				continue
			}
			if maintainer_email != "" && maintainer.Email != maintainer_email {
				continue
			}
		}
		metadata_list_filtered = append(metadata_list_filtered, metadata)
	}
	c.YAML(http.StatusOK, metadata_list_filtered)
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
