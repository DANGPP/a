package main

import (
	"context"
	"log"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"

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
	chartName       = "nginx" // Chú ý: KHÔNG có "bitnami/" ở đây
	chartRepo       = "https://charts.bitnami.com/bitnami"
	targetNamespace = "custom-nginx"
	helmReleaseName = "my-release"
)

func main() {
	log.Println("📦 Đang khởi tạo cấu hình Kubernetes...")

	// Load kubeconfig và khởi tạo clientset
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("❌ Lỗi kubeconfig: %v", err)
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("❌ Lỗi tạo clientset: %v", err)
	}
	settings.KubeConfig = kubeconfig

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
	chart, err := loader.Load(chartPath)
	if err != nil {
		log.Fatalf("❌ Không thể load chart: %v", err)
	}

	// Cài đặt chart
	log.Printf("🚀 Đang cài Helm release '%s' vào namespace '%s'...", helmReleaseName, targetNamespace)
	release, err := install.Run(chart, nil)
	if err != nil {
		log.Fatalf("❌ Lỗi khi cài chart: %v", err)
	}

	log.Printf("✅ Đã cài xong release: '%s' trong namespace: '%s'", release.Name, release.Namespace)
}
