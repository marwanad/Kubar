package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ghodss/yaml"
)

//ExecKubectlCmd ...
func ExecKubectlCmd(args []string) ([]byte, error) {
	cmd := exec.Command("kubectl", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("output: %s, error: %v", string(out), err)
	}
	return out, nil
}

//WriteToYAMLFile ...
func WriteToYAMLFile(out []byte, path string, filename string) error {
	y, err := yaml.JSONToYAML(out)
	if err != nil {
		return err
	}
	return WriteToFile(y, path, filename)
}

//WriteToFile ...
func WriteToFile(out []byte, path string, filename string) error {
	CreateDirIfNotExist(path)
	file, err := os.OpenFile(
		filepath.Join(path, fmt.Sprintf("%v.yaml", filename)),
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write bytes to file
	_, err = file.Write(out)
	if err != nil {
		return err
	}
	return nil
}

//CreateDirIfNotExist ...
func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// Contains ...
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
