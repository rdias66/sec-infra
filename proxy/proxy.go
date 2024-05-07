package proxy

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"sec-infra/utils"
	"strings"
)

func SetupSquidProxy() error {
	fmt.Println("Setting up proxy server with Squid...")

	if err := utils.InstallPackage("squid"); err != nil {
		fmt.Println("Failed to install Squid: ", err)
		return err
	}


	if err := utils.InstallPackage("apache2-utils"); err != nil {
		fmt.Println("Failed to install apache2-utils: ", err)
		return err
	}

	
	if err := generateSquidConfig(); err != nil {
		return err
	}
	return nil
}


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

	
	if err := restartSquid(); err != nil {
		return err
	}

	return nil
}

func restartSquid() error {
	fmt.Println("Restarting Squid service...")
	cmd := exec.Command("/etc/init.d/squid", "restart")
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to restart Squid: ", err)
		return err
	}
	return nil
}

func AddBlockedSite(url string) error {
<<<<<<< HEAD

=======
	
>>>>>>> 3c020d2653c45b801cc745fa953f38f804629be4
	sanitizedURL, err := sanitizeURL(url)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Opening Squid configuration file...")
	configFile, err := os.OpenFile("/etc/squid/squid.conf", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open squid.conf file: ", err)
		return err
	}
	defer configFile.Close()

	fmt.Println("Attempting to block access to the URL: ", url)

<<<<<<< HEAD

=======
	
>>>>>>> 3c020d2653c45b801cc745fa953f38f804629be4
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

<<<<<<< HEAD

=======
	
>>>>>>> 3c020d2653c45b801cc745fa953f38f804629be4
	if err := restartSquid(); err != nil {
		return err
	}

	return nil
}

func sanitizeURL(rawURL string) (string, error) {
<<<<<<< HEAD

=======
	
>>>>>>> 3c020d2653c45b801cc745fa953f38f804629be4
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

<<<<<<< HEAD
=======
	
>>>>>>> 3c020d2653c45b801cc745fa953f38f804629be4
	sanitizedScheme := strings.ToLower(parsedURL.Scheme)
	sanitizedHost := strings.ToLower(parsedURL.Host)
	sanitizedPath := parsedURL.Path

<<<<<<< HEAD
=======
	
>>>>>>> 3c020d2653c45b801cc745fa953f38f804629be4
	sanitizedURL := sanitizedScheme + "://" + sanitizedHost + sanitizedPath

	return sanitizedURL, nil
}
