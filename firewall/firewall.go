package firewall

import (
	"fmt"
	"os/exec"
)

// SetupFirewall sets up basic firewall rules and configures SSH within the container.
func SetupFirewall() error {
	// Check if iptables is installed
	fmt.Println("Setup firewall func accessed")
	if !isIptablesInstalled() {
		fmt.Println("iptables is not installed. Attempting to install...")
		err := installIptables()
		if err != nil {
			return fmt.Errorf("error installing iptables: %v", err)
		}
	}
	if isIptablesInstalled() {
		fmt.Println("Iptables is already installed.")
	}

	// Define the network interface (replace "eth0" with your actual interface)
	networkInterface := "eth0"

	//Add nf_tables command line binary path so the commands can work
	err := exec.Command("/usr/sbin/nft", "add", "rule", "ip", "filter", "input", "iifname", "eth0", "tcp", "dport", "22", "ct", "state", "new,established", "accept").Run()
	if err != nil {
		// Print the error if the command fails
		fmt.Printf("Error %v while trying to add nf_tables binary to path", err)
	}

	// Define firewall rules
	rules := []string{
		// Allow incoming SSH connections
		fmt.Sprintf("nft add rule ip filter input iifname %s tcp dport 22 ct state new,established accept", networkInterface),
		// Allow incoming HTTP connections
		fmt.Sprintf("nft add rule ip filter input iifname %s tcp dport 80 ct state new,established accept", networkInterface),
		// Drop all other incoming connections
		"nft add rule ip filter input drop",
	}

	// Apply each firewall rule
	for _, rule := range rules {
		// Print the rule before executing
		fmt.Println("Executing rule:", rule)

		// Execute the command
		cmd := exec.Command("bash", "-c", rule)
		err := cmd.Run()
		if err != nil {
			// Print the error if the command fails
			fmt.Printf("Error executing rule '%s': %v\n", rule, err)
			// Return an error if desired
			return fmt.Errorf("error applying firewall rule: %v", err)
		}
	}

	fmt.Println("Firewall rules configured successfully")

	//Install and configure SSH within the container
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
	fmt.Println("Install SSH func accessed")
	cmd := exec.Command("apt-get", "install", "-y", "iptables")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// installSSH installs and configures SSH within the container.
func installSSH() error {
	// Install required packages inside the container
	fmt.Println("Install SSH func accessed")
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
