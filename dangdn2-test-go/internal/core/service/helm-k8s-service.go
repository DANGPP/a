// core/service/service.go
package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type helmService struct {
	clientset  *kubernetes.Clientset
	kubeconfig string
	settings   *cli.EnvSettings
	chartName  string
	chartRepo  string
}

func NewHelmService() (HelmService, error) {
	kubeconfig := "./config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	svc := &helmService{
		clientset:  clientset,
		kubeconfig: kubeconfig,
		settings:   cli.New(),
		chartName:  "nginx",
		chartRepo:  "https://charts.bitnami.com/bitnami",
	}

	return svc, nil
}

func (h *helmService) ensureNamespace(namespace string) error {
	_, err := h.clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err == nil {
		return nil
	}

	_, err = h.clientset.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: namespace},
	}, metav1.CreateOptions{})

	return err
}

func (h *helmService) Deploy(namespace, releaseName string) (map[string]interface{}, error) {
	if err := h.ensureNamespace(namespace); err != nil {
		return nil, fmt.Errorf("cannot ensure namespace: %w", err)
	}

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(h.settings.RESTClientGetter(), namespace, "", log.Printf); err != nil {
		return nil, err
	}

	install := action.NewInstall(actionConfig)
	install.ReleaseName = releaseName
	install.Namespace = namespace
	install.ChartPathOptions.RepoURL = h.chartRepo

	chartPath, err := install.ChartPathOptions.LocateChart(h.chartName, h.settings)
	if err != nil {
		return nil, err
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}

	release, err := install.Run(chart, nil)
	if err != nil {
		return nil, err
	}

	go func() {
		time.Sleep(2 * time.Second)
		_ = os.Remove(chartPath)
	}()

	return map[string]interface{}{
		"message": "✅ Triển khai thành công",
		"release": release.Name,
		"status":  release.Info.Status.String(),
	}, nil
}

func (h *helmService) Update(namespace, releaseName string) (map[string]interface{}, error) {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(h.settings.RESTClientGetter(), namespace, "", log.Printf); err != nil {
		return nil, err
	}

	upgrade := action.NewUpgrade(actionConfig)
	upgrade.Namespace = namespace
	upgrade.ChartPathOptions.RepoURL = h.chartRepo

	chartPath, err := upgrade.ChartPathOptions.LocateChart(h.chartName, h.settings)
	if err != nil {
		return nil, err
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}

	release, err := upgrade.Run(releaseName, chart, nil)
	if err != nil {
		return nil, err
	}

	go func() {
		time.Sleep(2 * time.Second)
		_ = os.Remove(chartPath)
	}()

	return map[string]interface{}{
		"message": "🔁 Cập nhật thành công",
		"release": release.Name,
		"status":  release.Info.Status.String(),
	}, nil
}

func (h *helmService) Delete(namespace, releaseName string) (map[string]interface{}, error) {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(h.settings.RESTClientGetter(), namespace, "", log.Printf); err != nil {
		return nil, err
	}

	uninstall := action.NewUninstall(actionConfig)
	resp, err := uninstall.Run(releaseName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "🗑️ Xóa thành công",
		"release": releaseName,
		"info":    resp.Info,
	}, nil
}
