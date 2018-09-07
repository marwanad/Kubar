package export

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/marwanad/kubar/utils"
)

// GlobalResources ...
var GlobalResources = []string{"namespace", "storageclasses", "crd"}

// Export ...
func Export(path string) error {
	namespaces, err := getKubeNamespaces()
	if err != nil {
		return fmt.Errorf("Getting namespaces failed with the following error: %v", err)
	}

	globalResources, err := getGlobalResources()
	if err != nil {
		return fmt.Errorf("Getting global resources failed with the following error: %v", err)
	}
	cleanedUpGlobalResources, err := CleanupResourceFields(globalResources)

	if err != nil {
		return fmt.Errorf("Cleaning up global resources fields failed with the following error: %v", err)
	}

	for key, rsrc := range cleanedUpGlobalResources {
		utils.WriteToYAMLFile(rsrc, path, key)
	}

	for _, namespace := range strings.Split(namespaces[1:len(namespaces)-1], " ") {
		rsrcs, err := getKubeResources(namespace)
		if err != nil {
			return fmt.Errorf("[%v] Getting resources failed with the following error: %v", namespace, err)
		}
		cleanedUpResources, err := CleanupResourceFields(rsrcs)
		if err != nil {
			return fmt.Errorf("[%v] Cleaning up resources fields failed with the following error: %v", namespace, err)
		}
		for key, rsrc := range cleanedUpResources {
			utils.WriteToYAMLFile(rsrc, filepath.Join(path, fmt.Sprintf(namespace)), key)
		}
	}
	return nil
}

func getKubeNamespaces() (string, error) {
	kubectlArgs := []string{
		"get",
		"ns", "-o",
		"jsonpath='{.items[*].metadata.name}'",
	}

	out, err := utils.ExecKubectlCmd(kubectlArgs)

	if err != nil {
		return "", err
	}
	return string(out), nil
}

func getGlobalResources() (map[string][]byte, error) {
	globalResources := make(map[string][]byte)

	for _, resource := range GlobalResources {
		fmt.Printf("Querying global resource type: [%v]\n", resource)
		kubectlArgs := []string{
			"get",
			"--export", "-o=json",
			resource,
		}
		out, err := utils.ExecKubectlCmd(kubectlArgs)
		if err != nil {
			return nil, err
		}
		globalResources[resource] = out
	}
	return globalResources, nil
}

func getKubeResources(namespace string) (map[string][]byte, error) {
	resourceTypes := []string{"configmaps", "customresourcedefinition", "cronjobs", "daemonsets",
		"deployments", "ingresses", "replicasets", "role", "rolebinding", "secret", "services",
		"statefulsets", "storageclasses"}
	resources := make(map[string][]byte)

	for _, resource := range resourceTypes {
		fmt.Printf("[%v] Querying resource type: [%v]\n", namespace, resource)
		kubectlArgs := []string{
			"--namespace",
			namespace, "get",
			"--export", "-o=json",
			resource,
		}
		out, err := utils.ExecKubectlCmd(kubectlArgs)
		if err != nil {
			return nil, err
		}
		resources[resource] = out
	}
	return resources, nil
}
