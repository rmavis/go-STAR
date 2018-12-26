package main

import (
	"fmt"
)


func printUsageInformation() {
	msg := `USAGE

  INITIALIZING
    $ star -i [path/to/store/file]

    This command will ensure that your config and store files exist.
    The config file will be '~/.config/star/config.yaml'. And if you
    provide a path for the store file, that will be used. If not, the
    default is '~/.config/star/store'.


  CREATING
    $ star -n value[ tag...]

    This command will create a new entry in the store file. An entry
    consists of a value, any number of tags, and metadata (timestamps
    for the dates the value was created and last accessed and a count
    of the number of times the entry has been accessed).


  SEARCHING & ACTING
    $ star [flags] term[ term...]

    FLAGS
      -1, --one-line  Print output compressed to one line.
      -2, --two-line  Print output on two lines (value, tags).
      -a, --asc       Print records in ascending order.
      -b, --browse    Browse (do not select and pipe value to external tool).
      -d, --desc      Print output in descending order.
      -e, --edit      Edit the specified entries in your $EDITOR.
      -h, ---help     Print this help message.
      -i, --init      Initialize.
      -l, --loose     Match loosely.
      -n, --new       Create an entry.
      -p, --pipe      Pipe the value of the selected record to external tool.
      -s, --strict    Match strictly.
      -x, --delete    Delete the selected record(s).

    Searching is the default action. If no flags are given, the match
    mode (strict or loose) and action to take (external tool to pipe
    the value to) will be read from '~/.config/star/config.yaml'.
    Keys read from the config file are:
      store_file: ~/path/to/store/file
      filter_mode: (strict|loose)
      editor: /path/to/editor
      print_lines: (1|2)
      sort_order: (asc|desc)
      pipe_to: /path/to/tool

    If values are missing, these defaults will be used:
      store_file: ~/.config/star/store
      filter_mode: loose
      editor: $EDITOR or /usr/bin/vi
      print_lines: 2
      sort_order: desc
      pipe_to: {none}

    If no "pipe_to" action is present, then records will be printed
    to stdout.
`

	fmt.Println(msg)
}
