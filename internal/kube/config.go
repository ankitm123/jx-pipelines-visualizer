package kube

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewConfig(kubeConfigPath string) (*rest.Config, error) {
	// first, let's try to see if we are running in a pod in a cluster
	config, err := rest.InClusterConfig()
	if err == nil {
		_ = rest.SetKubernetesDefaults(config)
		return config, nil
	}

	// otherwise, fallback to using our kubeconfig path
	config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build kube config from %s: %w", kubeConfigPath, err)
	}

	_ = rest.SetKubernetesDefaults(config)
	return config, nil
}

func DefaultKubeConfigPath() string {
	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		return kubeconfig
	}

	home, _ := homedir.Dir()
	if home != "" {
		return filepath.Join(home, ".kube", "config")
	}

	wd, _ := os.Getwd()
	if wd != "" {
		return filepath.Join(wd, ".kube", "config")
	}

	return ""
}
