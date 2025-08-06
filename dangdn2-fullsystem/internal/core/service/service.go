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

	"github.com/hashicorp/vault/api"

	"database/sql"

	_ "github.com/lib/pq"
)

type helmService struct {
	clientset  *kubernetes.Clientset
	kubeconfig string
	settings   *cli.EnvSettings
	chartName  string
	chartRepo  string
}

type HelmService interface {
	DeployRelease(namespace, releaseName string) (map[string]interface{}, error)
	UpdateRelease(namespace, releaseName string) (map[string]interface{}, error)
	DeleteRelease(namespace, releaseName string) (map[string]interface{}, error)
	GetRelease(namespace, releaseName string) (map[string]interface{}, error)
}

func getVaultSecretPathFromDB(clusterName string) (string, error) {
	dbURL := "host=127.0.0.1 port=5433 user=dangdn2 password=1 dbname=db-part-k8s sslmode=disable"

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return "", fmt.Errorf("connect db: %w", err)
	}
	defer db.Close()

	var path string
	query := `SELECT path FROM k8s_configs WHERE name = $1 LIMIT 1`
	err = db.QueryRow(query, clusterName).Scan(&path)
	if err != nil {
		return "", fmt.Errorf("query path: %w", err)
	}

	return path, nil
}

func getKubeConfigFromVault(vaultAddr, token, secretPath string) (string, error) {
	config := &api.Config{Address: vaultAddr}
	client, err := api.NewClient(config)
	if err != nil {
		return "", fmt.Errorf("create vault client: %w", err)
	}
	client.SetToken(token)

	secret, err := client.Logical().Read("secret/data/" + secretPath)
	if err != nil {
		return "", fmt.Errorf("read secret: %w", err)
	}
	if secret == nil || secret.Data == nil || secret.Data["data"] == nil {
		return "", fmt.Errorf("secret not found or empty")
	}

	data := secret.Data["data"].(map[string]interface{})
	kcfg, ok := data["config"].(string)
	if !ok {
		return "", fmt.Errorf("'config' key not found or invalid in Vault data")
	}
	return kcfg, nil
}

func NewHelmService() (HelmService, error) {
	clusterName := "k8s-destop"
	secretPath, _ := getVaultSecretPathFromDB(clusterName)
	vaultAddr := ""
	vaultToken := ""
	// secretPath := "dev-cluster"
	kubeconfig, _ := getKubeConfigFromVault(vaultAddr, vaultToken, secretPath)

	// config, _ := clientcmd.BuildConfigFromFlags("", kubeconfqig)
	config, _ := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))

	clientset, _ := kubernetes.NewForConfig(config)

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

func (h *helmService) DeployRelease(namespace, releaseName string) (map[string]interface{}, error) {
	err := h.ensureNamespace(namespace)
	if err != nil {
		return nil, err
	}

	actionConfig := new(action.Configuration)
	err = actionConfig.Init(h.settings.RESTClientGetter(), namespace, "", log.Printf)
	if err != nil {
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
func (h *helmService) UpdateRelease(namespace, releaseName string) (map[string]interface{}, error) {
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

func (h *helmService) DeleteRelease(namespace, releaseName string) (map[string]interface{}, error) {
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

func (h *helmService) GetRelease(namespace, releaseName string) (map[string]interface{}, error) {
	actioConfig := new(action.Configuration)
	if err := actioConfig.Init(h.settings.RESTClientGetter(), namespace, "", log.Printf); err != nil {
		return nil, err
	}
	get := action.NewGet(actioConfig)
	release, err := get.Run(releaseName)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"name":         release.Name,
		"namespace":    release.Namespace,
		"chart":        release.Chart.Metadata.Name,
		"chartVersion": release.Chart.Metadata.Version,
		"status":       release.Info.Status.String(),
		"updated":      release.Info.LastDeployed.Time,
		"notes":        release.Info.Notes,
		"values":       release.Config,
	}, nil
}
