# atag

Extender of window's tags in **Acme**

## Introduction

**atag** extends all tags of existing and new windows of **Acme** with specified commands after `|`

## Using

Run **atag** with a list of commands you want to add in every window's tag.
If you want to add commands only for specific kind of files, you can specify a regular expression with a command list:
```
atag ".go:go build" ahist
```
In the case `go build` will be added to every window whose name is matched by `.go` and `ahist` will be added to every ***Acme***'s window.

Compound commands can be specified:
```
atag ".w:ahist 'make install'" "\"Edit s/one/two/g\""
```