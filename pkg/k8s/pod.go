package k8s

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"tektonclone/pkg/pipelines"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func CreatePod(ctx context.Context, repo string, steps pipelines.PipelineDef) (string, error) {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		return "", fmt.Errorf("home directory not found")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return "", err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	initContainers := []apiv1.Container{
		{
			Name:    "clone",
			Image:   "alpine/git",
			Command: []string{"git", "clone", repo, "/home/workspace"},
			VolumeMounts: []apiv1.VolumeMount{
				{
					Name:      "shared",
					MountPath: "/home/workspace",
				},
			},
		},
	}
	for _, s := range steps.Tasks {
		c := apiv1.Container{
			Name:       s.Name,
			Image:      s.Image,
			Command:    strings.Split(s.Script, " "),
			WorkingDir: "/home/workspace",
			VolumeMounts: []apiv1.VolumeMount{
				{
					Name:      "shared",
					MountPath: "/home/workspace",
				},
			},
		}

		initContainers = append(initContainers, c)
	}

	runnerPod := apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "trunner-",
			Namespace:    "default",
		},
		Spec: apiv1.PodSpec{
			RestartPolicy:  apiv1.RestartPolicyNever,
			InitContainers: initContainers,
			Containers: []apiv1.Container{
				{
					Name:    "exit",
					Image:   "alpine",
					Command: []string{"echo", "Done"},
				},
			},
			Volumes: []apiv1.Volume{
				{
					Name: "shared",
					VolumeSource: apiv1.VolumeSource{
						EmptyDir: &apiv1.EmptyDirVolumeSource{},
					},
				},
			},
		},
	}

	pod, err := clientset.CoreV1().Pods("default").Create(ctx, &runnerPod, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	return pod.Name, nil
}
