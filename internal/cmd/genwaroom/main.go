package main

import (
	"flag"
	"fmt"
	"os"
)

const usage = `usage:
  genwaroom analyze -spec <openapi.yaml> -ogen-config <.ogen.yaml> -ogen-target <dir> -config <config> -o <model.yaml>
      Run ogen's generator: write the low-level client to -ogen-target and
      derive the semantic SDK model from ogen's intermediate representation.
  genwaroom render -model <model.yaml> -out <dir>
      Generate the public waroom package from the model.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "analyze":
		err = analyzeCmd(os.Args[2:])
	case "render":
		err = renderCmd(os.Args[2:])
	default:
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "genwaroom:", err)
		os.Exit(1)
	}
}

func analyzeCmd(args []string) error {
	fs := flag.NewFlagSet("analyze", flag.ExitOnError)
	var (
		spec       = fs.String("spec", "api/openapi.yaml", "OpenAPI spec (ogen input)")
		ogenConfig = fs.String("ogen-config", ".ogen.yaml", "ogen generator config")
		ogenTarget = fs.String("ogen-target", "internal/ogen", "output directory for the ogen client")
		config     = fs.String("config", ".genwaroom.yaml", "config file (optional)")
		out        = fs.String("o", "api/sdk-model.yaml", "output model file")
	)
	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := loadConfig(*config)
	if err != nil {
		return err
	}
	model, g, err := analyze(*spec, *ogenConfig, cfg)
	if err != nil {
		return err
	}
	if err := writeOgen(g, *ogenTarget); err != nil {
		return err
	}
	return writeModel(*out, model)
}

func renderCmd(args []string) error {
	fs := flag.NewFlagSet("render", flag.ExitOnError)
	var (
		modelPath = fs.String("model", "api/sdk-model.yaml", "input model file")
		out       = fs.String("out", "./waroom", "output directory")
	)
	if err := fs.Parse(args); err != nil {
		return err
	}

	model, err := readModel(*modelPath)
	if err != nil {
		return err
	}
	return render(model, *out)
}
