package main

import (
	"fmt"
	"github.com/keltia/ripe-atlas"
	"github.com/urfave/cli"
	"log"
)

// init injects our "ip" related commands/options.
func init() {
	// Fill-in the various commands
	cliCommands = append(cliCommands, cli.Command{
		Name:        "ping",
		Usage:       "ping selected address",
		Description: "send echo/reply to an IP",
		Action:      cmdPing,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "T, tags",
				Usage:       "add tags to measurement",
				Destination: &fMTags,
			},
		},
	})
}

// shortcuts

func preparePing(target string) (req *atlas.MeasurementRequest) {
	opts := map[string]string{
		"Type":        "ping",
		"Description": fmt.Sprintf("Ping - %s", target),
		"Target":      target,
	}

	// Check global parameters
	opts = checkGlobalFlags(opts)

	// Add a tag?
	if fMTags != "" {
		opts["Tags"] = fMTags
	}

	// Try to configure -4/-6 depending on the argument to DTRT
	AF := prepareFamily(target)
	if AF == "" {
		AF = wantAF
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

	if fVerbose {
		displayOptions(opts)
	}

	return
}

// cmdIP is a short for displaying the IPs for one probe
func cmdPing(c *cli.Context) error {
	args := c.Args()
	if args == nil || len(args) != 1 {
		log.Fatal("Error: you must specify a hostname/IP")
	}

	addr := args[0]

	req := preparePing(addr)
	debug("req=%#v", req)

	m, err := client.Ping(req)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	//str := res.Result.Display()
	displayMeasurementID(*m)
	return nil
}
