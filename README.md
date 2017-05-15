bsh - An Alternate BOSH CLI
===========================

Lots of people rag on the `bosh` Ruby CLI utility because it's
Ruby and installing the Ruby runtime and all the supporting gems
can be a PITA on some jumpboxen.

There is a [go CLI][gocli] for BOSH, dubbed "BOSH 2".  It's a
single static binary, that makes for easy deployment, but it's not
a rewrite &mdash; it's a completely different utility, with
different aims, behaviors, and quirks.

`bsh` is my attempt to rewrite the Ruby CLI in Go (for the
static binary distribution) without losing the existing semantics
and ease-of-use that make the original BOSH CLI such a joy to work
with.


[gocli]: https://github.com/cloudfoundry/bosh-cli
