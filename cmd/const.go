package cmd

var VERSION = "1.2.0"

var HELP_STRING = `Use:
	todo <flag> <option>

Examples:
	todo -pBUG: -r
	todo --version

Flags:
	-d		Specify directory to look into
	-p 		Choose prefix (defualts to one in config.json)
	-r		Use relative paths (can also be set in config.json)
	-t		Print todos as raw text without formatting

Options:
	--help		What you are reading
	--version	Print program version

https://github.com/jesperkha/todo`
