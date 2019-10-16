package kubernetes

import (
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
)

// This is docker image provided by fx/contrib/docker_packer
// it can build a Docker image with give Docker project source codes encoded with base64
// check the detail fx/contrib/docker_packer/main.go
const image = "metrue/fx-docker"

func injectInitContainer(name string, deployment *appsv1.Deployment) *appsv1.Deployment {
	configMapHasToBeReady := true
	valueInConfigMapHasToBeReady := true
	initContainer := v1.Container{
		Name:            "fx-docker-build-c",
		Image:           image,
		ImagePullPolicy: v1.PullAlways,
		Command: []string{
			"/bin/sh",
			"-c",
			"/usr/bin/docker_packer $(APP_META) " + name,
		}, // Maybe it can be passed by Binary data from config map
		// Args:    []string{"${APP_META}"}, // function source codes and name
		VolumeMounts: []v1.VolumeMount{
			v1.VolumeMount{
				Name:      "dockersock",
				MountPath: "/var/run/docker.sock",
			},
		},
		Env: []v1.EnvVar{
			v1.EnvVar{
				Name: ConfigMap.AppMetaEnvName,
				ValueFrom: &v1.EnvVarSource{
					ConfigMapKeyRef: &v1.ConfigMapKeySelector{
						LocalObjectReference: v1.LocalObjectReference{Name: name},
						Key:                  ConfigMap.AppMetaEnvName,
						Optional:             &valueInConfigMapHasToBeReady,
					},
				},
			},
		},
		EnvFrom: []v1.EnvFromSource{
			v1.EnvFromSource{
				ConfigMapRef: &v1.ConfigMapEnvSource{
					LocalObjectReference: v1.LocalObjectReference{
						Name: name,
					},
					Optional: &configMapHasToBeReady,
				},
			},
		},
	}

	volumes := []v1.Volume{
		v1.Volume{
			Name: "dockersock",
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/var/run/docker.sock",
				},
			},
		},
	}
	deployment.Spec.Template.Spec.InitContainers = []apiv1.Container{initContainer}
	deployment.Spec.Template.Spec.Volumes = volumes
	return deployment
}
