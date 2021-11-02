package main

import (
	"net/http"

	"github.com/Mossaka/Application-Metadata-API-Server/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"gopkg.in/yaml.v3"
)

type MetadataPreview struct {
	Id string `yaml:"id"`
	models.Metadata
}

func get_metadata(c *gin.Context) {
	id := c.Param("id")
	raw_metadata, err := DB.Get(id)
	if err != nil {
		c.YAML(http.StatusNotFound, gin.H{"UserError": err.Error()})
		return
	}
	var metadata models.Metadata
	err = yaml.Unmarshal(raw_metadata, &metadata)
	if err != nil {
		c.YAML(http.StatusInternalServerError, gin.H{"InternalServerError": err.Error()})
		return
	}
	c.YAML(http.StatusOK, metadata)
}

func list_metadata(c *gin.Context) {
	metadata_id := c.Param("id")
	title := c.Query("title")
	version := c.Query("version")
	company := c.Query("company")
	website := c.Query("website")
	source := c.Query("source")
	license := c.Query("license")
	maintainer_name := c.Query("maintainer_name")
	maintainer_email := c.Query("maintainer_email")

	var metadata_list_filtered []MetadataPreview
	all_metadata := DB.GetAll()
	for id, raw_metadata := range all_metadata {
		var metadata models.Metadata
		err := yaml.Unmarshal(raw_metadata, &metadata)
		if err != nil {
			c.YAML(http.StatusInternalServerError, gin.H{"InternalServerError": err.Error()})
			return
		}
		if metadata_id != "" && metadata_id != id {
			continue
		}
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
		to_append := MetadataPreview{Id: id, Metadata: metadata}
		metadata_list_filtered = append(metadata_list_filtered, to_append)
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
	raw_metadata, err := yaml.Marshal(new_metadata)
	if err != nil {
		c.YAML(http.StatusInternalServerError, gin.H{"InternalServerError": err.Error()})
		return
	}
	DB.Add(raw_metadata)
	c.YAML(http.StatusOK, new_metadata)
}

func put_metadata(c *gin.Context) {
	id := c.Param("id")
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
	raw_metadata, err := yaml.Marshal(new_metadata)
	if err != nil {
		c.YAML(http.StatusInternalServerError, gin.H{"InternalServerError": err.Error()})
		return
	}
	DB.Set(id, raw_metadata)
	c.YAML(http.StatusOK, new_metadata)
}

func delete_metadata(c *gin.Context) {
	id := c.Param("id")
	DB.Delete(id)
	c.YAML(http.StatusOK, gin.H{"Deleted": id})
}
