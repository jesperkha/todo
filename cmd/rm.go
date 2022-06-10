package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/jesperkha/todo/util"
	"github.com/jesperkha/todo/writer"
)

var (
	itemIndexOutOfRange = errors.New("item index is out of range")
)

func newWriter(config configFile) *writer.TodoWriter {
	return writer.NewTodoWriter(config.Prefix, config.Depth, config.RelativePaths, config.IgnoreFiles, config.IgnoreDirs)
}

// Removes specified list item from file
func cmdRemove(config configFile, args []string) error {
	item, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		util.ErrAndExit(fmt.Errorf("invalid item index %s", args[0]))
	}

	w := newWriter(config)
	w.Read(".")
	if int(item)-1 >= len(w.GetList()) || int(item) == 0 {
		return itemIndexOutOfRange
	}

	todo := w.GetList()[item-1]
	file, err := os.ReadFile(todo.File)
	if err != nil {
		util.ErrAndExit(err)
	}

	// Todo: remove line

	return nil
}
