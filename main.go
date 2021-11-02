package main

import (
	"github.com/Mossaka/Application-Metadata-API-Server/models"
	"github.com/Mossaka/Application-Metadata-API-Server/persister"
	"gopkg.in/yaml.v3"

	"github.com/gin-gonic/gin"
)

var DB *persister.Persister

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

func setupServer() *gin.Engine {
	router := gin.Default()
	DB = persister.NewPersister()
	for _, metadata := range metadata_list {
		raw_metadata, err := yaml.Marshal(metadata)
		if err != nil {
			panic(err)
		}
		DB.Add(raw_metadata)
	}
	router.GET("/v1/metadata", list_metadata)
	router.GET("/v1/metadata/:id", get_metadata)
	router.POST("/v1/metadata", post_metadata)
	router.PUT("/v1/metadata/:id", put_metadata)
	router.DELETE("/v1/metadata/:id", delete_metadata)

	return router
}

func main() {
	setupServer().Run("localhost:8080")
}
