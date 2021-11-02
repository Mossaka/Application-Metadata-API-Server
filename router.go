package main

import (
	"github.com/Mossaka/Application-Metadata-API-Server/persister"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

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
