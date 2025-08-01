package service

type HelmService interface {
	Deploy(namespace string) (map[string]interface{}, error)
	Update(namespace, release string) (map[string]interface{}, error)
	Delete(namespace, release string) (map[string]interface{}, error)
}
