package main

import (
	"net/http"

	"github.com/Mossaka/Application-Metadata-API-Server/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

var metadata_list = []models.Metadata{
	{
		Title:   "Valid App 2",
		Version: "1.0.1",
		Maintainers: []models.Maintainer{
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
		Maintainers: []models.Maintainer{
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

	var metadata_list_filtered []models.Metadata
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
		found_name := false
		found_email := false
		for _, maintainer := range metadata.Maintainers {
			if maintainer_name != "" && maintainer.Name == maintainer_name {
				found_name = true
			}
			if maintainer_email != "" && maintainer.Email == maintainer_email {
				found_email = true
			}
			if found_name && found_email {
				break
			}
		}
		if maintainer_name != "" && !found_name {
			continue
		}
		if maintainer_email != "" && !found_email {
			continue
		}
		metadata_list_filtered = append(metadata_list_filtered, metadata)
	}
	c.YAML(http.StatusOK, metadata_list_filtered)
}

func post_metadata(c *gin.Context) {
	var new_metadata models.Metadata
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

func setupServer() *gin.Engine {
	router := gin.Default()
	router.GET("/v1/metadata", get_metadata)
	router.POST("/v1/metadata", post_metadata)

	return router
}

func main() {
	setupServer().Run("localhost:8080")
}
