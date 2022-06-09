package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var VERSION = "1.1.0"

var HELP_STRING = `Use:
	todo <flag> <option>

Examples:
	todo -pBUG: -r
	todo --version

Flags:
	-d		Specify directory to look into
	-p 		Choose prefix (defualts to one in config.json)
	-r		Use relative paths (can also be set in config.json)

Options:
	--help		What you are reading
	--version	Print program version

https://github.com/jesperkha/todo`

func printexit(f string, args ...any) {
	fmt.Println(fmt.Sprintf(f, args...))
	os.Exit(0)
}

type config struct {
	Prefix        string   `json:"prefix"`
	Depth         int      `json:"depth"`
	IgnoreFiles   []string `json:"ignoreFiles"`
	IgnoreDirs    []string `json:"ignoreDirs"`
	RelativePaths bool     `json:"relativePaths"`
}

func main() {
	var mainConfig config

	if configFile, err := os.ReadFile("config.json"); err == nil {
		if json.Unmarshal(configFile, &mainConfig) != nil {
			printexit("invalid input in config.json")
		}
	} else {
		// Defualt options
		mainConfig = config{Prefix: "Todo:", Depth: 5}
	}

	// Variable config set by flags
	relative := mainConfig.RelativePaths
	prefix := mainConfig.Prefix
	directory := "."

	for _, arg := range os.Args[1:] {
		option := arg
		switch option {
		case "--help":
			printexit(HELP_STRING)
		case "--version":
			printexit(VERSION)
		}

		argPrefix := arg[:2]
		value := arg[2:]
		switch argPrefix {
		case "-r":
			relative = true
		case "-d":
			directory = value
		case "-p":
			prefix = value
		default:
			printexit("unknown option '%s'", argPrefix)
		}
	}

	w := NewTodoWriter(prefix, mainConfig.Depth, relative, mainConfig.IgnoreFiles, mainConfig.IgnoreDirs)
	if err := w.Read(directory); err != nil {
		printexit("%s", err)
	}

	w.PrintList()
}
