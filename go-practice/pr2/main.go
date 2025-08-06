package main

import (
	// "net/http"
	"github.com/gin-gonic/gin"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"

	// "helm.sh/helm/v3/pkg/release"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

)
var(
	kubeconfig     ="C:/Users/dangt/.kube/config"
	clientset      *kubernetes.Clientset
)
func init(){
	// kết nối đến k8s:
	config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)
	// tạo đối tượng k8s:
	clientset,_:= kubernetes.NewForConfig(config)
}
func deployK(c *gin.Context){
	np := c.Param("namespace")
	re := c.Param("release")
}
func main(){
	r:= gin.Default()
	r.POST("/api/deploy/:namespace/:release", deployK)
	r.Run(":8080")
}