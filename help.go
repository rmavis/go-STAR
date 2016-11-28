package main

import (
	"fmt"
)





func printUsageInformation() {
	msg := `USAGE

  CREATING
     $ star -n value tag[, tag...]


  SEARCHING & ACTING
    $ star [flags] term[, term...]

    FLAGS
      -b, --browse  Browse (do not pipe value to external tool).
      -c, --copy    Pipe the value of the specified entries to 'pbcopy'.
      -d, --delete  Delete the specified entries.
      -e, --edit    Edit the specified entries in your $EDITOR.
      -l, --loose   Match loosely.
      -o, --open    Pipe the value of the specified entries to 'open'.
      -s, --strict  Match strictly.

    Searching is the default action. If no flags are given, the match
    mode (strict or loose) and action to take (external tool to pipe
    the value to) will be read from '~/.config/star/config.yaml'.
    Keys read from the config file are:
      file_name: ~/path/to/store/file
      filter_mode: (strict|loose)
      pipe_to: (copy|open|browse|/path/to/tool)
      editor: /path/to/editor

    If values are missing, these defaults will be used:
      file_name: ~/.config/star/store
      filter_mode: loose
      pipe_to: /usr/bin/pbcopy
      editor: $EDITOR or /usr/bin/vi

    The only exception to the search command pattern is for the
    browse action. If no search terms are given, then every entry
    will be printed. It's usually best to pipe the output from the
    browse action to your $PAGER, e.g.:
      $ star -b what ever | less
`

	fmt.Println(msg)
}
