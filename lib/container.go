package lib

import (
	"fmt"
	"os/exec"
)

// EnsureDockerWorking ensures docker is avaiable and function properly. Return error if not.
func EnsureDockerWorking() error {
	if output, err := exec.Command("docker", "ps").CombinedOutput(); err != nil {
		return fmt.Errorf("Failed to run docker command, do you have docker installed? do you need sudo?. Here is the output of docker ps: \n%s", output)
	}
	return nil
}

// IsContainerExist checks whether a container with the `name` already exists
func IsContainerExist(name string) bool {
	err := exec.Command("docker", "inspect", name).Run()
	return err == nil
}

// RunContainer runs a container with the `name`, `args`, and `imageName`. `imageName` can contains tag.
// Returns error if any step fails.
func RunContainer(name string, args []string, imageName string) error {
	cmdArgs := []string{"run", "--name", name}
	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, imageName)
	if output, err := exec.Command("docker", cmdArgs...).CombinedOutput(); err != nil {
		return fmt.Errorf("Failed to run docker %s. Here is the output: \n%s", name, output)
	}
	return nil
}

// RemoveContainer removes a container with the `name`. Returns error if any step fails.
func RemoveContainer(name string) error {
	if output, err := exec.Command("docker", "stop", name).CombinedOutput(); err != nil {
		return fmt.Errorf("Failed to stop docker %s. Here is the output: \n%s", name, output)
	}

	if output, err := exec.Command("docker", "rm", "-v", name).CombinedOutput(); err != nil {
		return fmt.Errorf("Failed to remove docker %s. Here is the output: \n%s", name, output)
	}

	return nil
}
