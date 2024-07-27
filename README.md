# hyphenlu

A command-line tool for detecting hyphens :)

# About

This tool is designed to make it easy to analyze strings and identify the number of hyphens in each one. It's a simple but useful tool for anyone working with strings.

# Features

Detects hyphens in a list of string.
Provides detailed information about the file, including the minimum and maximum number of hyphens.
Allows filtering by hyphen count using the `-f` flag.
Supports parallel processing for large files.

# Usage

To use the Hyphenlu Tool, simply run the following command:

`go run hyphenlu.go [-h] [-f <hyphen count>] <filename>`

# Switches

`-h` : Show this help message and exit.
`-f` : Filter strings by hyphen count.
`<filename>` : The file containing the list of strings.

# Examples

1. `go run main.go example.txt` : Analyze the strings in example.txt and print the results. 
2. `go run main.go -f 2 example.txt` : Filter the strings in example.txt to only show those with `2` hyphens.

# Installation

To install the Hyphen Detector Tool, simply clone this repository and run `go build hyphenlu.go`. This will create an executable file that you can use to run the tool.
