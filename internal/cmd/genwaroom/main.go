package main

import (
	"flag"
	"fmt"
	"os"
)

const usage = `usage:
  genwaroom analyze -src <ogen pkg> -spec <openapi.yaml> -config <config> -o <model.yaml>
      Derive the semantic SDK model from the compiled ogen package.
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
		src    = fs.String("src", "./internal/ogen", "package pattern of the ogen-generated client")
		spec   = fs.String("spec", "api/openapi.yaml", "OpenAPI spec (operationId → tag mapping only)")
		config = fs.String("config", ".genwaroom.yaml", "config file (optional)")
		out    = fs.String("o", "api/sdk-model.yaml", "output model file")
	)
	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := loadConfig(*config)
	if err != nil {
		return err
	}
	tags, err := loadTags(*spec)
	if err != nil {
		return err
	}
	model, err := load(*src, tags, cfg)
	if err != nil {
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
