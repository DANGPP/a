package anotherService

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/hashicorp/vault/api"

	"github.com/google/uuid"
)

type Another struct {
	VaultAddr  string `json:"vaultAddr"`
	VaultToken string `json:"vaultToken"`
	Path       string `json:"path"`
	Client     *api.Client
}

func NewAnotherService(vaultAddress, vaultToken, path string) (*Another, error) {
	config := &api.Config{Address: vaultAddress}
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	client.SetToken(vaultToken)

	ano := &Another{
		VaultAddr:  vaultAddress,
		VaultToken: vaultToken,
		Path:       path,
		Client:     client,
	}
	return ano, nil
}

// 1. gen và save secret key vào vault
func (a *Another) GenSecretKey() (string, string, error) {
	// 1. Sinh UUID
	id := uuid.NewString()

	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "ko gen ddc secretkey", "ko gen ddc secretkey", err
	}
	secretKey := hex.EncodeToString(bytes)

	// // 3. Tạo client Vault
	// config := &api.Config{Address: a.VaultAddr}
	// client, err := api.NewClient(config)
	// if err != nil {
	// 	return "failed to create Vault client", "failed to create Vault client", err
	// }
	// client.SetToken(a.VaultToken)

	// 4. Lưu vào Vault (KV v2 yêu cầu bọc trong "data")
	_, err = a.Client.Logical().Write("secret/data/"+a.Path+"/"+id, map[string]interface{}{
		"data": map[string]interface{}{
			"secretKey": secretKey,
		},
	})
	if err != nil {
		return "failed to write secret to Vault", "failed to write secret to Vault", err
	}

	return id, secretKey, nil
}

// get secretkey thong qua uuid
func (a *Another) GetSecretKey(uuid string) (string, error) {
	// config := &api.Config{Address: a.VaultAddr}
	// client, err := api.NewClient(config)
	// if err != nil {
	// 	return "", err
	// }
	// client.SetToken(a.VaultToken)

	vaultPath := "secret/data/" + a.Path + "/" + uuid
	secret, err := a.Client.Logical().Read(vaultPath)
	if err != nil {
		return "", err
	}
	if secret == nil || secret.Data == nil {
		return "secret not found at path", nil
	}
	// Lấy data bên trong
	dataRaw, ok := secret.Data["data"]
	if ok == false {

	}

	data, ok := dataRaw.(map[string]interface{})
	if ok == false {

	}
	// Lấy secretKey (đúng tên trường trong Vault)
	key, ok := data["secretKey"].(string)
	if ok == false {

	}

	return key, nil
}

// get full uuid
func (a *Another) GetFullUUIDkey() ([]string, error) {
	vaultPath := "secret/metadata/" + a.Path
	secret, err := a.Client.Logical().List(vaultPath)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no keys found at path %s", vaultPath)
	}

	keysRaw, ok := secret.Data["keys"]
	if !ok {
		return nil, fmt.Errorf("no keys field in secret metadata")
	}

	keys, ok := keysRaw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid keys format")
	}

	var result []string
	for _, k := range keys {
		if ks, ok := k.(string); ok {
			result = append(result, ks)
		}
	}

	return result, nil
}
