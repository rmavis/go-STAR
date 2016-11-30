# STAR: Simple Text Archiving and Retrieving

`star` is a simple tool for saving and retrieving bits of text. You could use it to save bookmarks, hard-to-remember commands, complex emoticons, stuff like that.


## Usage

I use `star` for all sorts of things.

    $ star music listen
    1) https://www.youtube.com/watch?v=Aloryfd5ipw
       analord, aphex twin, listen, music, music, afx
    2) http://modernarecords.com/releases
       label, listen, music, contemporary classical
    3) https://soundcloud.com/terryurban/sets/fka-biggie
       listen, music, terry urban, fka biggie
    4) http://prefuse73.bandcamp.com/album/travels-in-constants-vol-25
       music, prefuse 73, temporary residence, travels in constants, listen

    $ star book readme
    1) http://www.squeakland.org/resources/books/readingList.jsp
       alan kay, books, readme, syllabus
    2) http://www.gigamonkeys.com/book/
       Practical Common Lisp, readme, lisp
    3) https://tachyonpublications.com/product/central-station/
       books, readme, scifi
    4) https://en.wikipedia.org/wiki/The_World_as_Will_and_Representation
       books, philosophy, readme, schopenhauer
    5) http://www.goodreads.com/book/show/18877915-the-lonely-astronaut-on-christmas-eve
       blink-182, readme, space

    $ star emoticon
    1) ¯\_(ツ)_/¯
       emoticon, shrug
    2) (╯°□°）╯︵ ┻━┻
       emoticon, fliptable
    3) ⊙︿⊙
       bugeye, emoticon
    4) (╯︵╰,)
       emoticon, sadcry
    5) (*￣з￣)
       emoticon, kissing

Never [let your friends down][xkcd-tar] again!

    $ star unix tar
    1) tar -czvf DIR
       unix, tar, zip
    2) tar -tzvf DIR.tar.gz
       unix, tar, list
    3) tar -xzvf DIR.tar.gz
       unix, tar, unzip

That's the "Retrieving" part of `star`. Depending on your config (or command line flags) you can then pipe the string on the numbered line to `pbcopy`, `open`, or any other external tool you like.

The "Archiving" part is done like this:

    $ star -n value tag tag tag

So this:

    $ star -n "tar -czvf DIR" unix tar zip

will create a new record in the store file. `star` stores records in a plain text file, by default at `~/.config/star/store` but that's configurable via the config file at `~/.config/star/config.yaml`.

You can edit records via a temp file in your `$EDITOR`:

    $ star -e unix tar
    1) tar -czvf DIR
       unix, tar, zip
    2) tar -tzvf DIR.tar.gz
       unix, tar, list
    3) tar -xzvf DIR.tar.gz
       unix, tar, unzip
    Edit these records: 1 3

Deleting records is also easy:

    $ star -d todo
    1) shave
       todo, today
    2) shower
       todo, today
    3) wash dishes
       todo, today
    Delete these records: all

So essentially `star` saves, interfaces with, and acts on text snippets that it stores in a plain text file.



## Installation

This version of `star` is written in Go. In `bin/` there are compiled executables for every combination of `[linux darwin windows]` and `[amd64 386]` in the format `star_[os]-[arch]`. So to install you could just download the executable you need, put it somewhere useful, and create a symlink somewhere in your `$PATH`:

    $ cd wherever/you/keep/programs
    $ curl -O https://github.com/rmavis/go-STAR/raw/master/bin/star_[os]-[arch]
    $ ln -s [executable] /usr/local/bin/star

Or you can clone the repo and build it yourself. If you have a working Go setup, it's very easy:

    $ cd $GOPATH/src/where/ever
    $ git clone https://github.com/rmavis/go-STAR.git
    $ go get gopkg.in/yaml.v2  # This is star's only dependency.
    $ go install

There's also [a version written in Ruby][star-ruby], if you're into that. They're similar but I recommend this one---it's faster and has better browsing options.





[star-ruby]: https://github.com/rmavis/star
[xkcd-tar]: https://xkcd.com/1168/
