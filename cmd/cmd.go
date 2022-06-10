package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jesperkha/todo/util"
	"github.com/jesperkha/todo/writer"
)

type configFile struct {
	Prefix        string   `json:"prefix"`
	Depth         int      `json:"depth"`
	IgnoreFiles   []string `json:"ignoreFiles"`
	IgnoreDirs    []string `json:"ignoreDirs"`
	RelativePaths bool     `json:"relativePaths"`
}

func newWriter(config configFile) *writer.TodoWriter {
	return writer.NewTodoWriter(config.Prefix, config.Depth, config.RelativePaths, config.IgnoreFiles, config.IgnoreDirs)
}

func Run() {
	// Parse subcommand/option
	subcommand := ""
	for idx, arg := range os.Args[1:] {
		option := arg
		switch option {
		case "--help":
			util.PrintAndExit(HELP_STRING)
		case "--version":
			util.PrintAndExit(VERSION)
		default:
			if idx == 0 && !strings.HasPrefix(arg, "-") {
				subcommand = arg
			}
		}
	}

	// Load and parse json config file
	var config configFile
	if file, err := os.ReadFile("config.json"); err == nil {
		if json.Unmarshal(file, &config) != nil {
			util.ErrAndExit(fmt.Errorf("invalid input in config.json"))
		}
	} else {
		// Defualt options
		config = configFile{Prefix: "Todo:", Depth: 5}
	}

	// Run correct subcommand
	var err error
	switch subcommand {
	case "":
		cmdMain(config, os.Args[1:])
	case "rm":
		err = cmdRemove(config, os.Args[2:])
	default:
		util.ErrAndExit(fmt.Errorf("unknown command %s", subcommand))
	}

	if err != nil {
		util.ErrAndExit(err)
	}
}

func cmdMain(config configFile, args []string) {
	directory := "."

	for _, arg := range args {
		argPrefix := arg[:2]
		value := arg[2:]
		switch argPrefix {
		case "-r":
			config.RelativePaths = true
		case "-d":
			directory = value
		case "-p":
			config.Prefix = value
		default:
			util.ErrAndExit(fmt.Errorf("unknown option '%s'", argPrefix))
		}
	}

	// Read and print
	w := newWriter(config)
	if err := w.Read(directory); err != nil {
		util.ErrAndExit(err)
	}

	w.PrintList()
}
