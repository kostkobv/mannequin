package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kostkobv/mannequin"
	"github.com/kostkobv/mannequin/feat"
	"github.com/kostkobv/mannequin/feat/deploy"
	"github.com/kostkobv/mannequin/feat/deploy/latest"
	"github.com/kostkobv/mannequin/feat/implode"
	"github.com/kostkobv/mannequin/feat/initproject"
	"github.com/kostkobv/mannequin/feat/version"
)

var out = os.Stdout

func main() {
	// check if initialised first.
	if err := mannequin.CheckInited(); err != nil {
		fmt.Fprintf(out, "Couldn't execute the command: %s\n", err)
		fmt.Fprintln(out, "Mannequin is not initialised yet.")
		fmt.Fprintln(out, "Do you want to initialise the client now? (Y/n):")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(out, "Something went wrong: %s\n", err)
			os.Exit(2)
			return
		}
		text = strings.TrimSpace(text)
		if text != "y" && text != "Y" && text != "" {
			fmt.Fprintln(out, "Not initializing.")
			os.Exit(1)
			return
		}

		if err := mannequin.Initialize(out); err != nil {
			fmt.Fprintf(out, "Couldn't initialize: %s\n", err)
			os.Exit(2)
			return
		}
	}

	cfgPath, err := mannequin.ConfigPath()
	if err != nil {
		fmt.Fprintf(out, "Couldn't get the configuration path: %s\n", err)
		os.Exit(2)
		return
	}
	cfgFile, err := os.OpenFile(cfgPath, os.O_RDWR, 0775)
	if err != nil {
		fmt.Fprintf(out, "Couldn't open the configration file: %s\n", err)
		os.Exit(2)
		return
	}
	defer cfgFile.Close()
	cfg, err := mannequin.NewConfigFromFile(cfgFile)
	if err != nil {
		fmt.Fprintf(out, "Couldn't read the configration: %s\n", err)
		os.Exit(2)
		return
	}
	mnqn, err := mannequin.New(os.Stdout, &cfg)
	if err != nil {
		fmt.Fprintf(out, "Couldn't parse the configuration: %s\n", err)
		os.Exit(2)
		return
	}

	deployctl, err := deploy.New(latest.New())
	if err != nil {
		fmt.Fprintf(out, "Couldn't initialize deploy features: %s\n", err)
		os.Exit(2)
		return
	}

	// register available features.
	mnqnctl, err := feat.NewMnqnctlFeats(
		deployctl,
		initproject.New(),
		implode.New(),
		version.New(),
	)
	if err != nil {
		fmt.Fprintf(out, "Couldn't initialize features: %s\n", err)
		os.Exit(2)
		return
	}

	// execute the command.
	if err := mnqnctl.Do(mnqn, params(os.Args[1:])...); err != nil {
		fmt.Fprintf(out, "Couldn't execute: %s", err)
		os.Exit(1)
	}
}

func params(ps []string) []string {
	var res []string
	for _, p := range ps {
		if strings.HasPrefix(p, "-") {
			continue
		}

		res = append(res, p)
	}

	return res
}
