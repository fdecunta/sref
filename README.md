# sref-db

A command-line tool for managing academic references using DOIs and the CrossRef API.

## Overview

`sref-db` is a simple reference manager that stores bibliographic information in JSON format. It queries the CrossRef API to fetch metadata for academic papers and stores them locally for easy access and citation formatting.

## Installation

```bash
make install
```

This will:
- Build the `sref-db` binary
- Create configuration directory at `~/.config/sref/`
- Install the binary to `/usr/local/bin/`

Before using the tool, add your email to `~/.config/sref/email.conf`. This is part of the "polite" guidelines for the API:

```bash
echo "your.email@example.com" > ~/.config/sref/email.conf
```

## Usage

### Add a reference

```bash
sref-db add -doi "10.1007/s11104-024-06671-1"
```

With logging:

```bash
sref-db add -doi "10.1007/s11104-024-06671-1" -o log.txt
```

Maybe want to add a lot of DOIs from a file:

```bash
cat dois.txt | xargs -I '{}' sref-db add -doi '{}' -o log.txt
```

### Read references

Print all references:

```bash
sref-db read
```

Print a specific reference:

```bash
sref-db read -doi "10.1007/s11104-024-06671-1"
```

### Delete a reference

```bash
sref-db del -doi "10.1007/s11104-024-06671-1"
```

### Custom database file

By default, references are stored in `~/.config/sref/references.json`. Use the `-f` flag to specify a different file:

```bash
sref-db add -f custom.json -doi "10.1007/s11104-024-06671-1"
```

## Scripts

The `scripts/` directory contains helper utilities:

- `json2apa` - Bash script for APA formatting with jq. Not real APA, but kind of

## Uninstall

```bash
make uninstall
```

This removes the binary and configuration directory.

## Requirements

- Go 1.23.4 or later

## License

This project uses the [crossrefapi](https://github.com/caltechlibrary/crossrefapi) library from Caltech Library.
