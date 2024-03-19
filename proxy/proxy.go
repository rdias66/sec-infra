package proxy

import (
	"fmt"
	"os"
	"os/exec"
	"sec-infra/utils"
	"strings"
)

// SetupSquidProxy sets up the Squid proxy server
func SetupSquidProxy() error {
	fmt.Println("Setting up proxy server with Squid...")

	// Install Squid if not already installed
	if err := utils.InstallPackage("squid"); err != nil {
		fmt.Println("Failed to install Squid: ", err)
		return err
	}

	// Install apache2-utils if not already installed (used for proxy user password storage)
	if err := utils.InstallPackage("apache2-utils"); err != nil {
		fmt.Println("Failed to install apache2-utils: ", err)
		return err
	}

	// Generate Squid configuration
	if err := generateSquidConfig(); err != nil {
		return err
	}

	return nil
}

// generateSquidConfig generates Squid configuration
func generateSquidConfig() error {
	fmt.Println("Opening Squid configuration file...")
	configFile, err := os.OpenFile("/etc/squid/squid.conf", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Failed to open squid.conf file: ", err)
		return err
	}
	defer configFile.Close()

	fmt.Println("Writing basic Squid configuration into the file...")
	_, err = configFile.WriteString(`
# Basic Squid configuration
http_port 3128
auth_param basic program /usr/lib/squid/basic_ncsa_auth /etc/squid/passwd
auth_param basic realm Squid proxy
acl authenticated_users proxy_auth REQUIRED
http_access allow authenticated_users
`)
	if err != nil {
		fmt.Println("Failed to write basic configuration in squid config file: ", err)
		return err
	}

	// Restart Squid
	if err := restartSquid(); err != nil {
		return err
	}

	return nil
}

// restartSquid restarts the Squid service
func restartSquid() error {
	fmt.Println("Restarting Squid service...")
	cmd := exec.Command("/etc/init.d/squid", "restart")
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to restart Squid: ", err)
		return err
	}
	return nil
}

// AddBlockedSite adds a site to the list of blocked sites in Squid configuration
func AddBlockedSite(url string) error {
	// Sanitize URL to remove special characters
	sanitizedURL := strings.ReplaceAll(url, "/", "_")

	fmt.Println("Opening Squid configuration file...")
	configFile, err := os.OpenFile("/etc/squid/squid.conf", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open squid.conf file: ", err)
		return err
	}
	defer configFile.Close()

	fmt.Println("Attempting to block access to the URL: ", url)

	// Write ACL for the blocked URL
	_, err = configFile.WriteString(fmt.Sprintf("\nacl block_%s dstdomain %s\n", sanitizedURL, url))
	if err != nil {
		fmt.Println("Failed while blocking dstdomain: ", err)
		return err
	}
	_, err = configFile.WriteString(fmt.Sprintf("http_access deny block_%s\n", sanitizedURL))
	if err != nil {
		fmt.Println("Failed while blocking HTTP access: ", err)
		return err
	}

	// Restart Squid
	if err := restartSquid(); err != nil {
		return err
	}

	return nil
}
