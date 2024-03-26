package ssh

import (
	"fmt"
)

func InstallSSH() error {
	fmt.Println("Attempting to install SSH service... \n")
	// Install SSH server
	if err := utils.InstallPackage("openssh-server"); err != nil {
		fmt.Println(err)
	}
	return nil
}

func ConfigSSH() error {
	fmt.Println("Attemption to configure SSH...\n")
	
	sshPort := "22"
	permitRootLogin := "no"
	passwordAuthentication := "yes"

	err := exec.Command("sed", "-i", "s/^Port .*/Port "+sshPort+"/", "/etc/ssh/sshd_config").Run()
	if err != nil {
		fmt.Println(err)
	}

	err = exec.Command("sed", "-i", "s/^PermitRootLogin .*/PermitRootLogin "+permitRootLogin+"/", "/etc/ssh/sshd_config").Run()
	if err != nil {
		fmt.Println(err)
	}

	err = exec.Command("sed", "-i", "s/^PasswordAuthentication .*/PasswordAuthentication "+passwordAuthentication+"/", "/etc/ssh/sshd_config").Run()
	if err != nil {
		fmt.Println(err)
	}

	// Start SSH service
	if err := exec.Command("sudo", "service", "ssh", "start").Run(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("OpenSSH server has been installed and configured.")
	return nil
}
