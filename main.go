package main

import (
	"fmt"
	"sec-infra/firewall"
	"sec-infra/proxy"
	"sec-infra/ssh"
)

func main() {

	welcomeMessage := "Welcome to the sec-infra Golang CLI, this is an app developed to optimize and facilitate the security infrastructure of a Linux Server."
	fmt.Println(welcomeMessage)

	proxy.SetupSquidProxy()

	proxy.AddBlockedSite("https://www.instagram.com")

}
