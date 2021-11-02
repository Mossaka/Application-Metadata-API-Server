package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/Mossaka/Application-Metadata-API-Server/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestGetMetadata(t *testing.T) {
	ts := httptest.NewServer(setupServer())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/v1/metadata")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code is not 200. Status code: %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	m := []models.Metadata{}
	err = yaml.Unmarshal([]byte(b), &m)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(m), 2)
	assert.Equal(t, m[0].Maintainers[0].Name, "AppTwo Maintainer")
	assert.Equal(t, m[1].Maintainers[0].Name, "firstmaintainer app1")
}

func TestGetMetadataFilter(t *testing.T) {
	ts := httptest.NewServer(setupServer())
	defer ts.Close()

	urls := []string{"/v1/metadata?title=Valid%20App%202", "/v1/metadata?maintainer_email=secondmaintainer%40gmail.com"}

	for _, url := range urls {
		resp, err := http.Get(ts.URL + url)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Status code is not 200. Status code: %d", resp.StatusCode)
		}

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		m := []models.Metadata{}
		err = yaml.Unmarshal([]byte(b), &m)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, len(m), 1)
		raw_metadata := db.GetAll()[m[0].Title+m[0].Version]
		var metadata models.Metadata
		err = yaml.Unmarshal([]byte(raw_metadata), &metadata)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, m[0], metadata)
	}
}

func TestPostMetadata(t *testing.T) {
	ts := httptest.NewServer(setupServer())
	defer ts.Close()

	url := ts.URL + "/v1/metadata"
	payload := `title: Test App
maintainers:
- name: Test Maintainer
  email: testmaintain@google.com
- name: Test Maintainer 2
  email: testmaintainer2@google.com
company: Test Company
website: http://test.com
description: Test Description
source: https://github.com/random/repo
version: 1.0.0
license: MIT`
	resp, err := http.Post(url, "application/x-yaml", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code is not 200. Status code: %d", resp.StatusCode)
	}

	assert.Equal(t, len(db.GetAll()), 3)
	raw_metadata := db.GetAll()["Test App1.0.0"]
	var metadata models.Metadata
	err = yaml.Unmarshal([]byte(raw_metadata), &metadata)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, metadata.Maintainers[0].Name, "Test Maintainer")
}

func TestPostInvalidMetadata(t *testing.T) {
	ts := httptest.NewServer(setupServer())
	defer ts.Close()

	url := ts.URL + "/v1/metadata"
	payload := `title: App w/ Invalid maintainer email
version: 1.0.1
maintainers:
- name: Firstname Lastname
  email: apptwohotmail.com
company: Upbound Inc.
website: https://upbound.io
source: https://github.com/upbound/repo
license: Apache-2.0
description: |
  ### blob of markdown
  More markdown`
	resp, err := http.Post(url, "application/x-yaml", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Status code is not 400. Status code: %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, string(b), "Field validation for ''Email''\n  failed")
}

func TestPostInvalidMetadata2(t *testing.T) {
	ts := httptest.NewServer(setupServer())
	defer ts.Close()

	url := ts.URL + "/v1/metadata"
	payload := `title: App w/ missing version
maintainers:
- name: first last
  email: email@hotmail.com
- name: first last
  email: email@gmail.com
company: Company Inc.
website: https://website.com
source: https://github.com/company/repo
license: Apache-2.0
description: |
  ### blob of markdown
  More markdown`
	resp, err := http.Post(url, "application/x-yaml", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Status code is not 400. Status code: %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, string(b), "Field validation for ''Version'' failed")
}

func TestConcurrentPost(t *testing.T) {
	var wg sync.WaitGroup
	ts := httptest.NewServer(setupServer())
	defer ts.Close()

	url := ts.URL + "/v1/metadata"
	for i := 0; i < 10; i++ {
		payload := fmt.Sprintf(`title: Test App %d
maintainers:
  - name: Test Maintainer
    email: testmaintain@google.com
  - name: Test Maintainer 2
    email: testmaintainer2@google.com
company: Test Company
website: http://test.com
description: Test Description
source: https://github.com/random/repo
version: 1.0.0
license: MIT`, i)
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			defer wg.Done()
			http.Post(url, "application/x-yaml", strings.NewReader(payload))
			fmt.Printf("Done %d", i)
		}(&wg, i)
	}
	wg.Wait()
	assert.Equal(t, len(db.GetAll()), 12)
}
