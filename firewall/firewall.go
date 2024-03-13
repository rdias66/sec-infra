package firewall

import (
	"fmt"
	"os/exec"
)

// SetupFirewall sets up basic firewall rules and configures SSH within the container.
func SetupFirewall() error {
	fmt.Println("Setting up firewall...")

	// Install iptables if not already installed
	if err := installPackage("iptables"); err != nil {
		return fmt.Errorf("failed to install iptables: %v", err)
	}

	// Install nftables if not already installed
	if err := installPackage("nftables"); err != nil {
		return fmt.Errorf("failed to install nftables: %v", err)
	}

	// Define the network interface
	networkInterface := "eth0"

	// Define firewall rules
	rules := []string{
		"sudo nft add rule ip filter input iifname " + networkInterface + " tcp dport 22 ct state new,established accept",
		"sudo nft add rule ip filter input iifname " + networkInterface + " tcp dport 80 ct state new,established accept",
		"sudo nft add rule ip filter input drop",
	}

	// Apply each firewall rule
	for _, rule := range rules {
		fmt.Println("Applying rule:", rule)
		if err := runCommand(rule); err != nil {
			return fmt.Errorf("failed to apply firewall rule '%s': %v", rule, err)
		}
	}

	// Install and configure SSH
	if err := installSSH(); err != nil {
		return fmt.Errorf("failed to install and configure SSH: %v", err)
	}

	fmt.Println("Firewall setup completed successfully.")
	return nil
}

// installPackage installs a package using apt-get.
func installPackage(pkg string) error {
	fmt.Printf("Installing package %s...\n", pkg)
	cmd := exec.Command("apt-get", "install", "-y", pkg)
	if err := cmd.Run(); err != nil {
		return err
	}
	fmt.Printf("Package %s installed.\n", pkg)
	return nil
}

// runCommand runs a command in the shell.
func runCommand(command string) error {
	cmd := exec.Command("bash", "-c", command)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// installSSH installs and configures SSH.
func installSSH() error {
	fmt.Println("Installing and configuring SSH...")
	// Install SSH server
	if err := installPackage("openssh-server"); err != nil {
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

	// Start SSH service
	cmd := exec.Command("service", "ssh", "start")
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Println("SSH installed and configured successfully.")
	return nil
}
