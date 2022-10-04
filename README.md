# The 16 Bit Datacenter

This repository contains the source code and demos for a talk I gave at [CHINOG 10](https://chinog.org/chi-nog-10)

As soon as I figure out how Github Pages works, you can also view the slides and demos there.

## Slides
The presentation slides are generated using [backslide](https://github.com/sinedied/backslide)
If you have backslide installed, you can generate the static HTML for the slidedeck by running
`make slides`.  The slides will be stored in `dist/presentation.html`

## generate.go
There is a simple `generate.go` file that can be built to generate configs for
an Arista 7280SR2A in the manner described in the talk.  You can build this binary
by running `make generate` if you have GNU Make and Golang installed.

Once built, you can run `bin/generate -router-id XXX.XXX.XXX.XXX` to write the config to stdout

## Containerlab Demo
The lab folder contains the YAML and configurations to build a [containerlab](https://containerlab.dev/)
demo of Arista cEOS and Ubuntu [FRR](https://frrouting.org/) nodes.  You will need to provide your own cEOS
image -- update [lab/bitwiselclab.yaml line 5](lab/bitwise.clab.yaml#L5) with your local image.

Building the lab is straightforward

```bash
$ cd lab
$ containerlab deploy -c
```

Refer to the [containerlab documentation](https://containerlab.dev/quickstart/) for more information about
accessing and using the lab after it is deployed.
