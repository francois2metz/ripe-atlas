package main

import (
	"fmt"
	"github.com/keltia/ripe-atlas"
	"github.com/urfave/cli"
	"log"
	"strings"
)

// init injects our "sslcert" related commands/options.
func init() {
	// Fill-in the various commands
	cliCommands = append(cliCommands, cli.Command{
		Name: "sslcert",
		Aliases: []string{
			"tlscert",
			"tls",
		},
		Usage:       "get TLS certificate from host/IP",
		Description: "connect and fetch TLS certificate from host/IP",
		Action:      cmdTLSCert,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "T, tags",
				Usage:       "add tags to measurement",
				Destination: &fMTags,
			},
		},
	})
}

// prepareTraceroute build the request with our parameters
func prepareTLSCert(target string, port int) (req *atlas.MeasurementRequest) {
	opts := map[string]string{
		"Type":        "sslcert",
		"Description": fmt.Sprintf("SSLCert - %s", target),
		"Hostname":    target,
		"Target":      target,
		"Port":        fmt.Sprintf("%d", port),
	}

	// Check global parameters
	opts = checkGlobalFlags(opts)

	// Try to configure -4/-6 depending on the argument to DTRT
	AF := prepareFamily(target)
	if AF == "" {
		AF = wantAF
	}

	// Add a tag?
	if fMTags != "" {
		opts["Tags"] = fMTags
	}

	req = client.NewMeasurement()
	if AF == WantBoth {

		opts["AF"] = "4"
		req.AddDefinition(opts)

		opts["AF"] = "6"
		req.AddDefinition(opts)
	} else {
		opts["AF"] = AF
		req.AddDefinition(opts)
	}

	return
}

func cmdTLSCert(c *cli.Context) (err error) {
	args := c.Args()
	if len(args) == 0 {
		log.Fatal("Error: you must specify a hostname/site!")
	}

	// We expect target to be using <site>[:port]
	target := args[0]
	if !strings.HasPrefix(target, "http") {
		target = fmt.Sprintf("https://%s/", target)
	}

	_, site, _, port := analyzeTarget(target)

	req := prepareTLSCert(site, port)
	debug("req=%#v", req)
	//str := res.Result.Display()

	tls, err := client.SSLCert(req)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	displayMeasurementID(*tls)

	return
}
