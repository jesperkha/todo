package writer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/tabwriter"
)

type TodoWriter struct {
	writer     *tabwriter.Writer
	list       []Todo
	target     string
	maxDepth   int
	relative   bool
	fileIgnore map[string]struct{}
	dirIgnore  map[string]struct{}
}

type Todo struct {
	File       string
	Line       int
	Message    string
	ByteOffset int
}

func NewTodoWriter(target string, depth int, relativePaths bool, ignoredFiles []string, ignoredDirs []string) *TodoWriter {
	fileIgnore := map[string]struct{}{}
	dirIgnore := map[string]struct{}{}

	for _, f := range ignoredFiles {
		fileIgnore[f] = struct{}{}
	}
	for _, d := range ignoredDirs {
		dirIgnore[d] = struct{}{}
	}

	return &TodoWriter{
		writer:     tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0),
		list:       []Todo{},
		target:     target,
		maxDepth:   depth,
		relative:   relativePaths,
		fileIgnore: fileIgnore,
		dirIgnore:  dirIgnore,
	}
}

// Reads the given directory for todos and adds them to the writers list.
// The writer list is not cleared before reading.
func (td *TodoWriter) Read(dirname string) error {
	dir, err := os.ReadDir(dirname)
	if err != nil {
		return err
	}

	for _, entry := range dir {
		// Skip ignored directories
		if _, ok := td.dirIgnore[entry.Name()]; ok {
			continue
		}

		curdir := fmt.Sprintf("%s/%s", dirname, entry.Name())
		// Read subdirectory
		if entry.IsDir() {
			td.Read(curdir)
			continue
		}

		// Ingore certain file extensions
		if ext := strings.Split(entry.Name(), "."); len(ext) > 1 {
			if _, ok := td.fileIgnore[ext[1]]; ok {
				continue
			}
		}

		// Read file content as text
		content, err := ioutil.ReadFile(curdir)
		if err != nil {
			return err
		}

		line := 1
		maxPos := len(content) - len(td.target)
		for idx, b := range content {
			if b == '\n' {
				line++
			}

			if b != td.target[0] || idx > maxPos {
				continue
			}

			// Seek end of todo line and write todo data to list
			interval := content[idx : idx+len(td.target)]
			if string(interval) == td.target {
				endIdx := bytes.IndexByte(content[idx:], '\n')
				if endIdx == -1 {
					endIdx = len(content)
				}

				// Message after todo prefix
				message := string(content[idx+len(td.target) : endIdx+idx])
				message = strings.TrimSpace(message)

				filepath := dirname + "/" + entry.Name()
				if !td.relative {
					// Remove leading ./
					filepath = filepath[2:]
				}
				td.list = append(td.list, Todo{filepath, line, message, idx})
			}
		}
	}

	return nil
}

// Prints raw todos without any formatting
func (td *TodoWriter) PrintRaw() {
	for _, item := range td.list {
		fmt.Println(item.Message)
	}
}

// Prints the formatted list
func (td *TodoWriter) PrintList() {
	if len(td.list) == 0 {
		fmt.Println("no results...")
		return
	}

	fmt.Println()
	empty := "│ \t \t \t \t │\n"
	fmt.Fprint(td.writer, empty)

	longestPath := 0
	longestMsg := 0
	for idx, item := range td.list {
		filepath := fmt.Sprintf("%s:%d", item.File, item.Line)
		formatted := fmt.Sprintf("│ \t %d. \t %s \t %s \t │\n", idx+1, filepath, item.Message)
		td.writer.Write([]byte(formatted))

		if len(filepath) > longestPath {
			longestPath = len(filepath)
		}
		if len(item.Message) > longestMsg {
			longestMsg = len(item.Message)
		}
	}
	fmt.Fprint(td.writer, empty)

	length := longestMsg + longestPath

	header := "┌"
	footer := "└"
	for i := 0; i < length+14; i++ {
		header += "─"
		footer += "─"
	}

	headerRunes := []rune(header)
	title := fmt.Sprintf(" %s ", td.target)
	for i := 0; i < len(title); i++ {
		headerRunes[i+3] = rune(title[i])
	}
	header = string(headerRunes)

	fmt.Println(header + "┐")
	td.writer.Flush()
	fmt.Println(footer + "┘")
}

// Returns the writers current item list
func (td *TodoWriter) GetList() []Todo {
	return td.list
}
