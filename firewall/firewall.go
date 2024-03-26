package firewall

import (
	"fmt"
	"os/exec"
	"sec-infra/utils"
)

func InfraStarter() error {
	fmt.Println("Setting up firewall infrastructure...")
	utils.InstallPackage("iptables")
	utils.InstallPackage("kmod")

	mods := []string{
		"modprobe iptables_nat",
		"modprobe iptables_mangle",
		"modprobe iptables_filter",
	}
	for _, mod := range mods {
		fmt.Println("Installing adtional module for :", mod[8:])
		utils.RunCommand(mod)
	}

	enableIpForward := "echo 1 > /proc/sys/net/ipv4/ip_forward"
	utils.RunCommand(enableIpForward)

	return nil
}

func ClearRules() error {
	fmt.Println("Cleaning previous rule chain tables...")

	 clearCommands:= []string {
		"iptables -t nat -F",
		"iptables -t mangle -F",
		"iptables -t filter -F",
		"iptables -X",
	}

	for _,clearCommand := range clearCommands {
		utils.RunCommand(clearCommand)		
	}
	return nil
}

func ZeroCounters() error {
	fmt.Println("Zeroing rule counters...")
	zeroCommands := []string {
		"iptables -t nat -Z",
		"iptables -t mangle -Z",
		"iptables -t filter -Z",
	}
	for _,zeroCommand := range zeroCommands {
		util.RunCommand(zeroCommand)
	}
	
	return nil
}
