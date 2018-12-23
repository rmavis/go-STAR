# Arguments and Action Codes

The functions in `args.go` will read the user's input, interpret it, and build a slice of integers that, taken together, represents the user's intention.

This file exists to centralize the documentation and correlations between arguments (flags) and action codes (the slice of integers).


## Arguments (flags)

  -1, --one-line  Print compressed, one-line output.
  -2, --two-line  Print full, two-line output.
  -a, --asc       Sort records from low to high.
  -b, --browse    Show matching entries, take no action.
  -d, --desc      Sort records from high to low.
  -e, --edit      Edit an entry.
  -h, --help      Show this message.
  -i, --init      Create the ~/.config/star/store file.
  -l, --loose     Match loosely, rather than strictly.
  -m, --demo      Run the demo.
  -n, --new       Add a new entry.
  -p, --pipe      Pipe the selected record to an action.
  -s, --strict    Match strictly rather than loosely.
  -t, --tags      Show all tags.
  -v, --vals      Show all values.
  -x, --delete    Delete an entry.


## Action Codes

An action code is a slice of integers that encodes the user's intended action.

The first int indicates the command/operation:
1: select
2: create
3: help
4: view
5: initialize
6: demo

If the first int is a verb, the second int is an adverb. Only some verbs can be modified:
1 (select):
  0: default (read from config)
  1: pbcopy
  2: open
  3: edit
  4: delete
  5: browse
3 (help):
  1: commands
  2: readme
  3: customization
  4: examples
4 (view):
  1: values
  2: tags

And if the second int is an adverb, the third is an adjective:
1, 4 (select, browse):
  0: default (read from config)
  1: loose
  2: strict

The fourth indicates formatting.
  0: full list (two lines per record: value, then tags)
  1: compressed (one line per record)

The fifth indicates sort order.
  0: highest to lowest
  1: lowest to highest

A 0 in any place means that no relevant argument was given and that a sensible default should be used, if applicable.


## Argument Types

- Matching
  l (match loosely -- or)
  s (match strictly -- and)

- Action
  o (pipe value to `open`)
  c (pipe value to `pbcopy`)
  (((would be cool if user could pipe value to arbitrary program)))
  d (delete)
  e (edit)
  b (browse)

- Static messages
  c (help [commands])
  f (help [flags])
  h (help)
  r (readme)
  xh (extra help)
  xr (extra readme)
  x (examples)

- Dynamic messages (view?)
  a (all)
  t (tags)

- Creation
  i (init -- `touch` store file)
  n (add new entry)

- Compound
  m (demo)


## Usage Forms

- Find and operate on exiting entries:
  star -[ls][ocdeb] tag(s)
  Operation (Verb): search.
  Match mode (Adverb): determined by [ls] or default.
  Action: determined by [ocde] or default.

- See a help message
  star -[cfh][r][x]
  Operation (Verb): help message.
  Variation: determined by optional presence of `x`.

- See data dumps on existing entries
  star -[at]
  Operation (Verb): view.
  View material: determined by [at].

- Create a new entry
  star -n tag(s) value
  Operation (Verb): create.

- Touch the store file
  star -i
  Operation (Verb): initialize.

- Run the demo
  star -m -[ls][ocde] tag(s)
  Operation (Verb): demo.
