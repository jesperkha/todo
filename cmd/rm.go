package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/jesperkha/todo/util"
)

var (
	itemIndexOutOfRange = errors.New("item index is out of range")
)

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

	// Remove line and stitch file back together
	lines := bytes.Split(file, []byte{'\n'})
	stitched := append(lines[:todo.Line-1], lines[todo.Line:]...)

	joined := bytes.Join(stitched, []byte{'\n'})
	if err := os.WriteFile(todo.File, joined, os.ModeAppend); err != nil {
		return err
	}

	return nil
}
