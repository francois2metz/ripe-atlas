// cmd_measurements.go

package main

import (
	"fmt"
	"github.com/keltia/ripe-atlas"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

// init injects our measurement-related commands
func init() {
	cliCommands = append(cliCommands, cli.Command{
		Name: "measurements",
		Aliases: []string{
			"measures",
			"m",
		},
		Usage:       "measurements-related keywords",
		Description: "All the commands for measurements",
		Subcommands: []cli.Command{
			{
				Name:        "list",
				Aliases:     []string{"ls"},
				Usage:       "lists all measurements",
				Description: "displays all measurements",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:        "country, c",
						Usage:       "filter on country",
						Destination: &fCountry,
					},
					cli.StringFlag{
						Name:        "asn",
						Usage:       "filter on asn",
						Value:       "",
						Destination: &fAsn,
					},
					cli.BoolFlag{
						Name:        "A, all",
						Usage:       "all measurements even inactive ones",
						Destination: &fAllMeasurements,
					},
					cli.BoolFlag{
						Name:        "is-anchor",
						Usage:       "select anchor measurements",
						Destination: &fWantAnchor,
					},
					cli.BoolFlag{
						Name:        "mine",
						Usage:       "limit to my own measurements",
						Destination: &fWantMine,
					},
					cli.StringFlag{
						Name:        "t, type",
						Usage:       "filter on measurement type",
						Value:       "",
						Destination: &fMeasureType,
					},
				},
				Action: measurementsList,
			},
			{
				Name:        "info",
				Usage:       "info for one measurement",
				Description: "gives info for one measurement",
				Action:      measurementInfo,
			},
			{
				Name:        "results",
				Aliases:     []string{"r", "res"},
				Usage:       "results for one measurement",
				Description: "returns results for one measurement",
				Action:      measurementResults,
			},
			{
				Name:        "delete",
				Aliases:     []string{"rm", "del", "destroy"},
				Usage:       "info for one measurement",
				Description: "stops one measurement (or all)",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:        "A",
						Usage:       "select all measurements",
						Destination: &fAllMeasurements,
					},
				},
				Action: measurementDelete,
			},
		},
	})
}

// displayMeasurement returns a string with one measurement
func displayMeasurement(m *atlas.Measurement, verbose bool) (res string) {
	if verbose {
		res = fmt.Sprintf("%v\n", m)
	} else {
		res = fmt.Sprintf("ID: %d type: %s description: %s Target: %s\n", m.ID, m.Type, m.Description, m.Target)
	}
	return
}

// displayAllMeasurements returns a string with all measurements
func displayAllMeasurements(ml *[]atlas.Measurement, verbose bool) (res string) {
	res = ""
	for _, m := range *ml {
		res += displayMeasurement(&m, verbose)
	}
	return
}

// displayResult returns a string with <obj>.Result
func displayResult(body []byte, verbose bool) (res string) {
	return string(body)
}

// measurementsList returns a list of measurements according to parameters
func measurementsList(c *cli.Context) error {
	opts := make(map[string]string)

	if fCountry != "" {
		opts["country_code"] = fCountry
	}

	if fAsn != "" {
		opts["asn"] = fAsn
	}

	if fMeasureType != "" {
		opts["type"] = fMeasureType
	}

	if fWantMine {
		opts["mine"] = "true"
	}

	// Check global parameters
	opts = checkGlobalFlags(opts)

	if fVerbose {
		displayOptions(opts)
	}

	q, err := atlas.GetMeasurements(opts)
	if err != nil {
		log.Printf("GetMeasurements err: %v - q:%v", err, q)
		os.Exit(1)
	}
	log.Printf("Got %d measurements with %v\n", len(q), opts)
	fmt.Print(displayAllMeasurements(&q, fVerbose))

	return nil
}

// measurementInfo is for one measurement only
func measurementInfo(c *cli.Context) error {
	args := c.Args()
	if args[0] == "" {
		log.Fatalf("Error: you must specify a measurement ID!")
	}

	id, _ := strconv.ParseInt(args[0], 10, 32)

	p, err := atlas.GetMeasurement(int(id))
	if err != nil {
		fmt.Printf("err: %v", err)
		os.Exit(1)
	}
	fmt.Print(displayMeasurement(p, fVerbose))

	return nil
}

// measurementResults returns the result part of the measurement
func measurementResults(c *cli.Context) error {
	args := c.Args()
	if args[0] == "" {
		log.Fatalf("Error: you must specify a measurement ID!")
	}

	id, _ := strconv.ParseInt(args[0], 10, 32)

	m, err := atlas.GetMeasurement(int(id))
	if err != nil {
		fmt.Printf("err: %v", err)
		os.Exit(1)
	}

	// m.Result is an URI pointing to results, fetch it
	if m.Result == "" {
		fmt.Println("Empty result")
	}

	resp, err := http.Get(m.Result)
	if err != nil {
		fmt.Errorf("bad net/http answer for %s: %v", m.Result, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("error reading body for %s: %v", m.Result, err)
	}

	fmt.Print(displayResult(body, fVerbose))
	return nil
}

func measurementCreate(c *cli.Context) error {
	return nil
}

func measurementDelete(c *cli.Context) (err error) {
	var id int64

	if fAllMeasurements {
		opts := make(map[string]string)

		// Check global parameters
		opts = checkGlobalFlags(opts)

		list, err := atlas.GetMeasurements(opts)
		if err != nil {
			fmt.Errorf("Delete all failed: %v", err)
		} else {
			for _, m := range list {
				err = atlas.DeleteMeasurement(m.ID)
				if err != nil {
					err = fmt.Errorf("Error: can not delete measurement %d", m.ID)
				}
			}
			fmt.Printf("All measurements stopped.")
		}
	} else {
		args := c.Args()
		if args[0] == "" {
			log.Fatalf("Error: you must specify a measurement ID!")
		}

		id, _ = strconv.ParseInt(args[0], 10, 32)
		err = atlas.DeleteMeasurement(int(id))
	}

	return
}