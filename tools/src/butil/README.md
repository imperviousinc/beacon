# butil

This is just a dummy utility that can be installed with `go install`. It clones the actual repo, builds and runs the real `butil`. It also makes sure the real beacon utility is always up-to-date. You can modify `realbutil` and this tool will take care of compiling the new code.

### Usage

Running it with default options:

```
$ mkdir beacon && cd beacon
$ butil clone
```

The clone command first calls `git clone` internally, so you can pass other options to git like repository url and `--branch` name of a branch other than main example:

```
$ butil clone https://github.com/imperviousinc/beacon --branch <name>
```
