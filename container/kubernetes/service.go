package kubernetes

import (
	"github.com/metrue/fx/constants"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

// CreateService create a service
func (k *K8S) CreateService(
	namespace string,
	name string,
	typ string,
	ports []int32,
	podsLabels map[string]string,
) (*v1.Service, error) {
	servicePorts := []v1.ServicePort{
		v1.ServicePort{
			Name:       "http",
			Protocol:   v1.ProtocolTCP,
			Port:       80,
			TargetPort: intstr.FromInt(int(constants.FxContainerExposePort)),
		},
		v1.ServicePort{
			Name:       "https",
			Protocol:   v1.ProtocolTCP,
			Port:       443,
			TargetPort: intstr.FromInt(int(constants.FxContainerExposePort)),
		},
	}
	// Append custom Port
	for _, port := range ports {
		servicePorts = append(servicePorts, v1.ServicePort{
			Name:       "custom",
			Protocol:   v1.ProtocolTCP,
			Port:       port,
			TargetPort: intstr.FromInt(int(3000)),
		})
	}

	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			ClusterName: namespace,
		},
		Spec: v1.ServiceSpec{
			Ports:    servicePorts,
			Type:     v1.ServiceType(typ),
			Selector: podsLabels,
		},
	}

	createdService, err := k.CoreV1().Services(namespace).Create(service)
	if err != nil {
		return nil, err
	}

	return createdService, nil
}

// DeleteService a service
func (k *K8S) DeleteService(namespace string, name string) error {
	// TODO figure out the elegant way to delete a service
	options := &metav1.DeleteOptions{}
	return k.CoreV1().Services(namespace).Delete(name, options)
}
