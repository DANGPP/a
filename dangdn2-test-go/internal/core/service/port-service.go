package service

type HelmService interface {
	Deploy(namespace, releaseName string) (map[string]interface{}, error)
	Update(namespace, releaseName string) (map[string]interface{}, error)
	Delete(namespace, releaseName string) (map[string]interface{}, error)
}
