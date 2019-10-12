package kubernetes

import (
	"strconv"

	"github.com/metrue/fx/constants"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

func generateServiceSpec(
	namespace string,
	name string,
	typ string,
	ports []int32,
	selector map[string]string,
) *apiv1.Service {
	servicePorts := []apiv1.ServicePort{
		apiv1.ServicePort{
			Name:       "http",
			Protocol:   apiv1.ProtocolTCP,
			Port:       80,
			TargetPort: intstr.FromInt(int(constants.FxContainerExposePort)),
		},
		apiv1.ServicePort{
			Name:       "https",
			Protocol:   apiv1.ProtocolTCP,
			Port:       443,
			TargetPort: intstr.FromInt(int(constants.FxContainerExposePort)),
		},
	}
	// Append custom Port
	for index, port := range ports {
		servicePorts = append(servicePorts, apiv1.ServicePort{
			Name:       "custom-port-" + strconv.Itoa(index),
			Protocol:   apiv1.ProtocolTCP,
			Port:       port,
			TargetPort: intstr.FromInt(int(3000)),
		})
	}

	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			ClusterName: namespace,
		},
		Spec: apiv1.ServiceSpec{
			Ports:    servicePorts,
			Type:     apiv1.ServiceType(typ),
			Selector: selector,
		},
	}
}

// CreateService create a service
func (k *K8S) CreateService(
	namespace string,
	name string,
	typ string,
	ports []int32,
	selector map[string]string,
) (*apiv1.Service, error) {
	service := generateServiceSpec(namespace, name, typ, ports, selector)
	createdService, err := k.CoreV1().Services(namespace).Create(service)
	if err != nil {
		return nil, err
	}

	return createdService, nil
}

// UpdateService update a service
// TODO this method is not perfect yet, should refactor later
func (k *K8S) UpdateService(
	namespace string,
	name string,
	typ string,
	ports []int32,
	selector map[string]string,
) (*apiv1.Service, error) {
	svc, err := k.GetService(namespace, name)
	if err != nil {
		return nil, err
	}
	svc.Spec.Selector = selector
	svc.Spec.Type = apiv1.ServiceType(typ)
	return k.CoreV1().Services(namespace).Update(svc)
}

// DeleteService a service
func (k *K8S) DeleteService(namespace string, name string) error {
	// TODO figure out the elegant way to delete a service
	options := &metav1.DeleteOptions{}
	return k.CoreV1().Services(namespace).Delete(name, options)
}

// GetService get a service
func (k *K8S) GetService(namespace string, name string) (*apiv1.Service, error) {
	return k.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
}
