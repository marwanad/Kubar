package restore

import (
	"fmt"
	"os"

	"io/ioutil"
	"path/filepath"

	"github.com/marwanad/kubar/export"
	"github.com/marwanad/kubar/utils"
)

// Restore ...
func Restore(path string) error {
	if err := restoreGlobalResources(path); err != nil {
		return err
	}

	if err := restoreKubeResources(path); err != nil {
		return err
	}
	fmt.Println("Sucessfully restored Kubernetes resources")
	return nil
}

func restoreGlobalResources(path string) error {
	for _, rsc := range export.GlobalResources {
		if err := applyKubectl(path, fmt.Sprintf("%v.yaml", rsc)); err != nil {
			return err
		}
	}
	return nil
}

func restoreKubeResources(path string) error {
	initialchildPaths := []string{"kube-system", "kube-public", "default"}

	for _, childPath := range initialchildPaths {
		if err := applyKubectl(path, childPath); err != nil {
			return err
		}
	}

	// search for custom childPaths and apply the kube resources
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() && !utils.Contains(initialchildPaths, f.Name()) {
			if err := applyKubectl(path, f.Name()); err != nil {
				return err
			}
		}
	}
	return nil
}

func applyKubectl(basePath string, childPath string) error {
	rscPath := filepath.Join(basePath, childPath)
	if _, err := os.Stat(rscPath); os.IsNotExist(err) {
		fmt.Printf("[%v] Resource does not exist. Skipping restoration\n", childPath)
		return nil
	}
	fmt.Printf("[%v] Restoring resources\n", childPath)

	kubectlArgs := []string{
		"apply",
		"-n", childPath,
		"-f", filepath.Join(basePath, childPath),
	}
	out, err := utils.ExecKubectlCmd(kubectlArgs)

	if err != nil {
		return fmt.Errorf("output: %s, error: %v", string(out), err)
	}
	fmt.Printf("output: %s\n", string(out))

	return nil
}
