package main

import (
	"context"
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig = "C:/Users/Admin/.kube/config"
	clientset  *kubernetes.Clientset
	ctx        = context.TODO()
)

func init() {
	// Load config từ file kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("❌ Failed to load kubeconfig: %v", err)
	}

	// Tạo clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("❌ Failed to create clientset: %v", err)
	}
}

func getNodes(c *gin.Context) {
	nodes, err := clientset.CoreV1().Nodes().List(ctx, v1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nodes.Items) // Trả về danh sách node
}
func getNamespaces(c *gin.Context) {
	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, namespaces)
}
func createNamespace(c *gin.Context) {
	type NamespaceRequest struct {
		Name string `json:"name"`
	}

	var req NamespaceRequest
	if err := c.BindJSON(&req); err != nil || req.Name == "" {
		c.JSON(400, gin.H{"error": "Invalid JSON or missing namespace name"})
		return
	}

	ns := &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: req.Name,
		},
	}

	createdNS, err := clientset.CoreV1().Namespaces().Create(ctx, ns, v1.CreateOptions{})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, createdNS)
}
func deleteNamespace(c *gin.Context) {
	name := c.Param("name")

	err := clientset.CoreV1().Namespaces().Delete(ctx, name, v1.DeleteOptions{})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Namespace deleted", "namespace": name})
}

func main() {
	r := gin.Default()

	// GET localhost:8080/api/nodes
	r.GET("/api/nodes", getNodes)
	r.GET("/api/namespaces", getNamespaces)
	r.POST("/api/namespaces", createNamespace)
	r.DELETE("/api/namespaces/:name", deleteNamespace)

	// log.Println("🚀 Server running on http://localhost:8080")
	r.Run(":8080")
}
