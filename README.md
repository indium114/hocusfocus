# Hocusfocus

*Hocusfocus* is a terminal-based productivity tracker written in Go!,

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/StikyPiston/hocusfocus)

## Installation

```shell
brew install stikypiston/formulae/hocusfocus
```

## Usage

Run the `hocusfocus` command to enter the main interface

There are **three** session types to choose from: `Work`, `Study`, and `Waste`.

Use the `arrow keys` to navigate up and down, and use the `enter` key to select an option.

Press `enter` on one of the session types to start tracking it. When you are done, press `stop` to stop, or press another session's name to switch to that session.

Running `hocusfocus` with the `stats` argument will print how long you've spent in total in each session.

Running `hocusfocus` with the `currentsession` argument will print the current session to the terminal.
