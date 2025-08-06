package main

import (
	"encoding/json"
	"io"

	// "log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var UrlHead = "http://harbor.local:30080/api/v2.0"

type Project struct {
	ProjectID int               `json:"project_id"`
	Name      string            `json:"name"`
	Metadata  map[string]string `json:"metadata"`
}
type Tag struct {
	Name string `json:"name"`
}

type ArtifactShort struct {
	Digest   string `json:"digest"`
	Size     int64  `json:"size"`
	PushTime string `json:"push_time"`
	Tags     []Tag  `json:"tags"`
}

func getHealth(c *gin.Context) {
	url := UrlHead + "/health"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	// c.Data(resp.StatusCode, "appication/json", body)
	c.JSON(resp.StatusCode, body)

}
func getRepoArtifacts(c *gin.Context) {
	project := c.Param("project")
	repo := c.Param("repo")
	url := UrlHead + "/projects/" + project + "/repositories/" + repo + "/artifacts"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}
	req.Header.Add("accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "request failed"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read response"})
		return
	}

	var artifacts []ArtifactShort
	if err := json.Unmarshal(body, &artifacts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse response"})
		return
	}

	c.JSON(http.StatusOK, artifacts)
}

func getRepositories(c *gin.Context) {
	project := c.Param("project")
	url := UrlHead + "/projects/" + project + "/repositories"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

func getPublicProjects(c *gin.Context) {
	url := UrlHead + "/projects"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}
	req.Header.Add("accept", "application/json")
	req.SetBasicAuth("admin", "Harbor12345") // nếu cần login

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to call Harbor API"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		c.JSON(resp.StatusCode, gin.H{"error": string(bodyBytes)})
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read response"})
		return
	}

	var projects []Project
	if err := json.Unmarshal(bodyBytes, &projects); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse JSON"})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func main() {
	r := gin.Default()

	r.GET("/harbor/projects", getPublicProjects)

	r.GET("/harbor/health", getHealth)

	r.GET("/harbor/:project/:repo/artifacts", getRepoArtifacts)

	r.GET("/harbor/:project/repositories", getRepositories)

	r.Run(":8080") // chạy tại http://localhost:8080
}
