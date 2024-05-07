package main

import (
	"fmt"
	"sec-infra/firewall"
	"sec-infra/proxy"
	"sec-infra/ssh"
)

func main() {
	proxy.SetupSquidProxy()
	proxy.AddBlockedSite("https://www.instagram.com")

}
