/*
This package is just a collection of test cases
*/
package main

import (
	"github.com/codegangsta/cli"
	"os"
)

var (
	// flags
	fWant4 bool
	fWant6 bool
	fAllProbes bool
	fAsn string
	fCountry string
	fVerbose bool
)

// main is the starting point (and everything)
func main() {
	app := cli.NewApp()
	app.Name = "atlas"
	app.Usage = "RIPE Atlas cli interface"
	app.Author = "Ollivier Robert <roberto@keltia.net>"
	app.Version = "0.0.1"

	// General flags
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "v",
			Usage: "more verbose",
			Destination: &fVerbose,
		},
	}

	// Fill-in the various commands
	app.Commands = []cli.Command{
		{
			Name: "probes",
			Aliases: []string{
				"p",
				"pb",
			},
			Usage:       "probe-related keywords",
			Description: "All the commands for probes",
			Subcommands: []cli.Command{
				{
					Name:        "list",
					Aliases:     []string{"ls"},
					Usage:       "lists all probes",
					Description: "displays all probes",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "country,c",
							Usage: "filter on country",
							Value: "fr",
							Destination: &fCountry,
						},
						cli.StringFlag{
							Name: "asn",
							Usage: "filter on asn",
							Value: "",
							Destination: &fAsn,
						},
						cli.BoolFlag{
							Name: "A",
							Usage: "all probes even inactive ones",
							Destination: &fAllProbes,
						},
					},
					Action:      probesList,
				},
				{
					Name:        "info",
					Usage:       "info for one probe",
					Description: "gives info for one probe",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "country,c",
							Usage: "filter on country",
							Value: "fr",
							Destination: &fCountry,
						},
						cli.StringFlag{
							Name: "asn",
							Usage: "filter on asn",
							Value: "",
							Destination: &fAsn,
						},
					},
					Action:      probeInfo,
				},
			},
		},
		{
			Name:        "ip",
			Usage:       "returns current ip",
			Description: "shorthand for getting current ip",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "ipv6",
					Usage: "displays only IPv6",
					Destination: &fWant6,
				},
				cli.BoolFlag{
					Name:  "ipv4",
					Usage: "displays only IPv4",
					Destination: &fWant4,
				},
			},
			Action: cmdIP,
		},
	}
	app.Run(os.Args)

}
