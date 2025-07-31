package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"

	// "helm.sh/helm/v3/pkg/release"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig      = "C:/Users/Admin/.kube/config"
	clientset       *kubernetes.Clientset
	ctx             = context.TODO()
	settings        = cli.New()
	chartName       = "nginx"
	chartRepo       = "https://charts.bitnami.com/bitnami"
	targetNamespace = "custom-nginx"
	helmReleaseName = "my-release"
	loadedChart     *chart.Chart
)

func init() {
	// kết nối và tạo đối tượng k8s
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("%v", err)
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Tạo namespace nếu cần
	log.Printf("🔍 Kiểm tra namespace '%s'...", targetNamespace)
	_, err = clientset.CoreV1().Namespaces().Get(ctx, targetNamespace, metav1.GetOptions{})
	if err != nil {
		log.Printf("📁 Namespace '%s' chưa tồn tại. Đang tạo...", targetNamespace)
		_, err = clientset.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: targetNamespace,
			},
		}, metav1.CreateOptions{})
		if err != nil {
			log.Fatalf("❌ Không thể tạo namespace: %v", err)
		}
		log.Printf("✅ Đã tạo namespace: %s", targetNamespace)
	} else {
		log.Printf("✅ Namespace '%s' đã tồn tại", targetNamespace)
	}

	// Cấu hình Helm
	log.Println("⚙️ Khởi tạo Helm...")
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), targetNamespace, "", log.Printf); err != nil {
		log.Fatalf("❌ Lỗi init Helm: %v", err)
	}

	install := action.NewInstall(actionConfig)
	install.ReleaseName = helmReleaseName
	install.Namespace = targetNamespace
	install.ChartPathOptions.RepoURL = chartRepo

	// Tải chart
	log.Printf("🌐 Đang tìm chart '%s' từ repo: %s", chartName, chartRepo)
	chartPath, err := install.ChartPathOptions.LocateChart(chartName, settings)
	if err != nil {
		log.Fatalf("❌ Không tìm được chart: %v", err)
	}
	log.Printf("📥 Đã tải chart tại: %s", chartPath)

	// Load chart
	log.Println("📂 Đang load chart...")
	loadedChart, err = loader.Load(chartPath)
	if err != nil {
		log.Fatalf("❌ Không thể load chart: %v", err)
	}
}

func deploy(c *gin.Context) {
	namespace := c.Param("namespace")

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, "", log.Printf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	install := action.NewInstall(actionConfig)
	install.ReleaseName = helmReleaseName
	install.Namespace = namespace

	release, err := install.Run(loadedChart, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "✅ Triển khai thành công",
		"release": release.Name,
		"status":  release.Info.Status.String(),
	})
}

func update(c *gin.Context) {
	namespace := c.Param("namespace")
	releaseName := c.Param("release")

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, "", log.Printf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	upgrade := action.NewUpgrade(actionConfig)
	upgrade.Namespace = namespace
	upgrade.ChartPathOptions.RepoURL = chartRepo // <- thêm dòng này

	// 🆕 Tải lại chart mới từ Harbor mỗi lần update
	chartPath, err := upgrade.ChartPathOptions.LocateChart(chartName, settings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không tìm thấy chart: " + err.Error()})
		return
	}
	latestChart, err := loader.Load(chartPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể load chart: " + err.Error()})
		return
	}

	release, err := upgrade.Run(releaseName, latestChart, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cập nhật thất bại: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "🔁 Cập nhật thành công",
		"release": release.Name,
		"status":  release.Info.Status.String(),
	})
}

func delete(c *gin.Context) {
	namespace := c.Param("namespace")
	releaseName := c.Param("release")

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, "", log.Printf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	uninstall := action.NewUninstall(actionConfig)

	resp, err := uninstall.Run(releaseName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "🗑️ Xóa thành công",
		"release": releaseName,
		"info":    resp.Info,
	})
}

func main() {
	r := gin.Default()
	r.POST("/api/deploy/:namespace", deploy)
	r.POST("/api/update/:namespace/:release", update)
	r.DELETE("/api/delete/:namespace/:release", delete)

	r.Run(":8080")
}
