package main

import (
	"encoding/json"
	"fmt"
	"os"
)

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
		err = json.Unmarshal(configFile, &mainConfig)
		if err != nil {
			fmt.Println("invalid input in config.json")
			return
		}
	} else {
		// Defualt options
		mainConfig = config{Prefix: "Todo:", Depth: 5}
	}

	var (
		relative = mainConfig.RelativePaths
		prefix   = mainConfig.Prefix
		dir      = "."
	)

	args := os.Args[1:]
	for _, arg := range args {
		argPrefix := arg[:2]
		value := arg[2:]

		switch argPrefix {
		case "-r":
			relative = true
		case "-d":
			dir = value
		case "-p":
			prefix = value
		default:
			fmt.Printf("unknown option '%s'\n", argPrefix)
			return
		}
	}

	w := NewTodoWriter(prefix, mainConfig.Depth, relative, mainConfig.IgnoreFiles, mainConfig.IgnoreDirs)
	if err := w.Read(dir); err != nil {
		fmt.Println(err)
		return
	}

	w.PrintList()
}
