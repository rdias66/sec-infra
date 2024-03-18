package utils

import (
	"fmt"
	"os/exec"
)

func InstallPackage(pkg string) error {
	fmt.Printf("Attempting to install package %s...\n", pkg)
	cmd := exec.Command("apt-get", "install", "-y", pkg)
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Package %s installed.\n", pkg)
	return nil
}

func RunCommand(command string) error {
	cmd := exec.Command("bash", "-c", command)
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
	return nil
}