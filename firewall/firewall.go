package firewall

import (
	"fmt"
	"os/exec"
)

// SetupFirewall sets up basic firewall rules and configures SSH within the container.
func SetupFirewall() error {
	// Check if iptables is installed
	if !isIptablesInstalled() {
		fmt.Println("iptables is not installed. Attempting to install...")
		err := installIptables()
		if err != nil {
			return fmt.Errorf("error installing iptables: %v", err)
		}
	}

	// Define the network interface (replace "eth0" with your actual interface)
	networkInterface := "eth0"

	// Define firewall rules
	rules := []string{
		// Allow incoming SSH connections
		fmt.Sprintf("iptables -A INPUT -i %s -p tcp --dport 22 -m conntrack --ctstate NEW,ESTABLISHED -j ACCEPT", networkInterface),
		// Allow incoming HTTP connections
		fmt.Sprintf("iptables -A INPUT -i %s -p tcp --dport 80 -m conntrack --ctstate NEW,ESTABLISHED -j ACCEPT", networkInterface),
		// Drop all other incoming connections
		"iptables -A INPUT -j DROP",
	}

	// Apply each firewall rule
	for _, rule := range rules {
		err := exec.Command("bash", "-c", rule).Run()
		if err != nil {
			return fmt.Errorf("error applying firewall rule: %v", err)
		}
	}

	fmt.Println("Firewall rules configured successfully")

	// Install and configure SSH within the container
	err := installSSH()
	if err != nil {
		return fmt.Errorf("error installing and configuring SSH: %v", err)
	}

	return nil
}

// isIptablesInstalled checks if iptables is installed.
func isIptablesInstalled() bool {
	_, err := exec.LookPath("iptables")
	return err == nil
}

// installIptables installs iptables.
func installIptables() error {
	cmd := exec.Command("apt", "install", "-y", "iptables")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// installSSH installs and configures SSH within the container.
func installSSH() error {
	// Install required packages inside the container
	err := exec.Command("apt-get", "update").Run()
	if err != nil {
		return err
	}

	err = exec.Command("apt-get", "install", "-y", "sudo", "openssh-server").Run()
	if err != nil {
		return err
	}

	// Set SSH configuration
	sshPort := "22"
	permitRootLogin := "no"
	passwordAuthentication := "yes"

	err = exec.Command("sed", "-i", "s/^Port .*/Port "+sshPort+"/", "/etc/ssh/sshd_config").Run()
	if err != nil {
		return err
	}

	err = exec.Command("sed", "-i", "s/^PermitRootLogin .*/PermitRootLogin "+permitRootLogin+"/", "/etc/ssh/sshd_config").Run()
	if err != nil {
		return err
	}

	err = exec.Command("sed", "-i", "s/^PasswordAuthentication .*/PasswordAuthentication "+passwordAuthentication+"/", "/etc/ssh/sshd_config").Run()
	if err != nil {
		return err
	}

	// Start SSH service
	err = exec.Command("sudo", "service", "ssh", "start").Run()
	if err != nil {
		return err
	}

	fmt.Println("OpenSSH server has been installed and configured.")

	return nil
}
