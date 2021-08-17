package main

import (
	"fmt"
	"log"
	"os"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/logging"
	"github.com/scrapli/scrapligo/transport"
	"github.com/srl-labs/srlinux-scrapli"
)

func main() {
	// uncomment to enable logging
	logging.SetDebugLogger(log.Print)

	contName := "srlinux"
	tlsProfileName := "demo"

	d, err := srlinux.NewSRLinuxDriver(
		contName,
		base.WithAuthStrictKey(false),
		base.WithAuthBypass(true),
	)
	if err != nil {
		fmt.Printf("failed to create driver; error: %+v\n", err)
		return
	}
	defer d.Close()

	transport, _ := d.Transport.(*transport.System)
	transport.ExecCmd = "docker"
	transport.OpenCmd = []string{"exec", "-u", "root", "-it", contName, "sr_cli"}
	fmt.Println("Configuring TLS certificate...")
	err = srlinux.AddSelfSignedServerTLSProfile(d, tlsProfileName, false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	d.AcquirePriv("configuration")

	fmt.Println("Enabling gNMI...")
	gnmiCfg := []string{
		"set / system gnmi-server admin-state enable",
		"set / system gnmi-server timeout 7200",
		"set / system gnmi-server rate-limit 60",
		"set / system gnmi-server session-limit 20",
		"set / system gnmi-server network-instance mgmt admin-state enable",
		"set / system gnmi-server network-instance mgmt use-authentication true",
		"set / system gnmi-server network-instance mgmt port 57400",
		fmt.Sprintf("set / system gnmi-server network-instance mgmt tls-profile %s", tlsProfileName),
		"commit save",
	}

	_, err = d.SendConfigs(gnmiCfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
