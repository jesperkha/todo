# todo

CLI application for reading markers in codebases (typically todos).

<br>

## Use:

By default, it will run with configuration set in `config.json`
```console
$ todo

┌── Todo ──────────────────────────────────────────────────────────────┐
│                                                                      │
│   1.   main.go:20        remove debug code after testing             │
│   2.   other.go:6        include variable image sizes in export      │
│   3.   test/test.go:68   add test cases for resizing and bordering   │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

Flags:
- `-r`: use relative paths (overrides config)
- `-d`: specify directory
- `-p`: specify prefix (overrides config)

Options:
- `--help`: display help info
- `--version`: print version

Example:
```console
$ todo -r -pTODO -dsome_dir
```

<br>

## Commands

### `rm`

Remove a todo from a file (removes the line the todo is located at). The list index is the number on the left of the list item when printing (starts at 1).

`todo rm <list_index>`

<br>

## Config

You can add a `config.json` file in the same directory as the todo binary. If no config file is found then default options will be used. These are the config options you can put in you config file:

```js
// These are also the default values
{
    "prefix": "Todo:",
    // Max search depth
    "depth": 5,
    "relativePaths": false,
    // Ignored file extensions (exe, png etc)
    "ignoreFiles": [],
    // Ignored directories (.vscode, .git etc)
    "ignoreDirs": []
}
```
