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

	// Install squid if not installed yet
	if err := utils.InstallPackage("squid"); err != nil {
		fmt.Println("failed to install squid: ", err)
		return err
	}

	// Install apache2 utils if not installed yet for proxy users pass storage
	if err := utils.InstallPackage("apache2-utils"); err != nil {
		fmt.Println("failed to install apache2-utils: ", err)
		return err
	}

	if err := generateSquidConfig(); err != nil {
		return err
	}

	return nil
}

// GenerateSquidConfig generates Squid configuration
func generateSquidConfig() error {
	fmt.Println("Checking if  Squid configuration file is done...")
	if _, err := os.Stat("/etc/squid/squid.conf"); err == nil {
		fmt.Println("Squid configuration file already exists. Skipping generation.")
		return nil
	}

	fmt.Println("Opening Squid configuration file...")
	configFile, err := os.OpenFile("/etc/squid/squid.conf", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("failed to open squid.conf file: ", err)
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
		fmt.Println("failed to write basic configuration in squid config file: ", err)
		return err
	}

	// restart Squid using shell command
	cmd := exec.Command("/etc/init.d/squid", "restart")
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to restart squid: ", err)
		return err
	}

	return nil
}

func AddBlockedSite(url string) error {
	// Sanitize URL to remove special characters
	sanitizedURL := strings.ReplaceAll(url, "/", "_")

	fmt.Println("Opening Squid configuration file...")
	configFile, err := os.OpenFile("/etc/squid/squid.conf", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("failed to open squid.conf file: ", err)
		return err
	}
	defer configFile.Close()

	fmt.Println("Attempting to write access control list for the url: ", url)

	_, err = configFile.WriteString(fmt.Sprintf("\nacl block_%s dstdomain %s\n", sanitizedURL, url))
	if err != nil {
		fmt.Println("failed while blocking dstdomain: ", err)
		return err
	}
	_, err = configFile.WriteString(fmt.Sprintf("http_access deny block_%s\n", sanitizedURL))
	if err != nil {
		fmt.Println("failed while blocking http access: ", err)
		return err
	}

	// restart Squid using shell command
	cmd := exec.Command("/etc/init.d/squid", "restart")
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to restart squid: ", err)
		return err
	}

	return nil
}
