package kubernetes

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateConfigMap create a config map with data
func (k *K8S) CreateConfigMap(namespace string, name string, data map[string]string) (*apiv1.ConfigMap, error) {
	cm := &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: data,
	}
	return k.CoreV1().ConfigMaps(namespace).Create(cm)
}

// DeleteConfigMap delete a config map
func (k *K8S) DeleteConfigMap(namespace string, name string) error {
	return k.CoreV1().ConfigMaps(namespace).Delete(name, &metav1.DeleteOptions{})
}

// UpdateConfigMap update a config map
func (k *K8S) UpdateConfigMap(namespace string, name string, data map[string]string) (*apiv1.ConfigMap, error) {
	cm := &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: data,
	}
	return k.CoreV1().ConfigMaps(namespace).Update(cm)
}

// GetConfigMap get a config map
func (k *K8S) GetConfigMap(namespace string, name string) (*apiv1.ConfigMap, error) {
	return k.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})
}

// CreateOrUpdateConfigMap create or update a config map
func (k *K8S) CreateOrUpdateConfigMap(namespace string, name string, data map[string]string) (*apiv1.ConfigMap, error) {
	_, err := k.GetConfigMap(namespace, name)
	if err != nil {
		return k.CreateConfigMap(namespace, name, data)
	}
	return k.UpdateConfigMap(namespace, name, data)
}
