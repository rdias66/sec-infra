package proxy

import (
	"fmt"
	"os"
	"os/exec"
	"sec-infra/utils"
)

// SquidConfig represents Squid proxy server configuration parameters
type SquidConfig struct {
	BlockedSites  []string
	Authenticated bool
	Username      string
	Password      string
}

func SetupSquidProxy() error {
	fmt.Println("Setting up proxy server with Squid...")

	//Install squid if not installed yet
	if err := utils.InstallPackage("squid"); err != nil {
		fmt.Println("failed to install squid: ", err)
	}

	if err := generateSquidConfig(); err != nil {
		fmt.Println("failed to generate Squid config: ", err)
	}
	return nil
}

// GenerateSquidConfig generates Squid configuration based on the provided parameters
func generateSquidConfig() error {
	fmt.Println("Opening Squid configuration file...")
	configFile, err := os.OpenFile("/etc/squid/squid.conf", os.O_APPEND|os.O_WRONLY, 0644)
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
	fmt.Println("Restarting Squid...")
	cmd := exec.Command("systemctl", "restart", "squid")
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to restart squid: ", err)
		return err
	}

	return nil
}

// AddUserToSquidConfig adds a user to Squid authentication configuration
func AddUserToSquidConfig(username, password string) error {
	fmt.Println("Attempting to add user:", username, " to proxy server...")
	cmd := exec.Command("htpasswd", "-b", "-c", "/etc/squid/passwd", username, password)
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to add new user to proxy server: ", err)
		return err
	}
	return nil
}

func AddBlockedSite(url string) error {
	fmt.Println("Opening Squid configuration file...")
	configFile, err := os.OpenFile("/etc/squid/squid.conf", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("failed to open squid.conf file: ", err)
		return err
	}
	defer configFile.Close()

	fmt.Println("Attempting to write access control list for the url: ", url)

	_, err = configFile.WriteString(fmt.Sprintf("\nacl block_%s dstdomain %s\n", url, url))
	if err != nil {
		fmt.Println("failed while blocking dstdomain: ", err)
		return err
	}
	_, err = configFile.WriteString(fmt.Sprintf("http_access deny block_%s\n", url))
	if err != nil {
		fmt.Println("failed while blocking http access: ", err)
		return err
	}

	fmt.Println("Attempting to restard squid... ")
	cmd := exec.Command("systemctl", "restart", "squid")
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to restart squid: ", err)
		return err
	}

	return nil
}
