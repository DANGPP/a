package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/vault/api"
)

func getKubeConfigFromVault(vaultAddr, token, secretPath string) (string, error) {
	config := &api.Config{Address: vaultAddr}
	client, err := api.NewClient(config)
	if err != nil {
		return "", fmt.Errorf("create vault client: %w", err)
	}
	client.SetToken(token)

	// Lưu ý: vault kv v2 => đường dẫn sẽ là "secret/data/<your-path>"
	secret, err := client.Logical().Read(fmt.Sprintf("secret/data/%s", secretPath))
	if err != nil {
		return "", fmt.Errorf("read secret: %w", err)
	}
	if secret == nil || secret.Data["data"] == nil {
		return "", fmt.Errorf("secret not found or empty")
	}
	data := secret.Data["data"].(map[string]interface{})
	kcfg, ok := data["config"].(string)
	if !ok {
		return "", fmt.Errorf("config key not found in Vault data")
	}
	return kcfg, nil
}

func main() {
	vaultAddr := "http://127.0.0.1:8205"
	vaultToken := ""
	secretPath := "dev-cluster"

	fmt.Println("🔐 Đang lấy kubeconfig từ Vault...")

	kubeconfig, err := getKubeConfigFromVault(vaultAddr, vaultToken, secretPath)
	if err != nil {
		log.Fatalf("❌ Lỗi: %v", err)
	}

	// ✅ Lưu trong biến kubeconfig để dùng luôn, không cần ghi ra file
	fmt.Println("✅ Kubeconfig:")
	fmt.Println(kubeconfig)

	// Bạn có thể dùng kubeconfig như sau:
	// config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
}
