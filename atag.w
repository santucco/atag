\input header

@** Introduction.
This is an implementation of \.{atag} command for \.{Acme}. It adds specified commands to a tag of every \.{Acme}'s window
or only in windows, matched by a regular expression.

@** Implementation.
@c

@i license

package main

import (
	@<Imports@>
)@#

var (
	@<Global variables@>
)@#

@
@<Imports@>=
"fmt"
"os"

@ At first, if no commands are specified, let's print the usage info and exit.
Then an enumeration of opened windows is processed in a separated goroutine.
Then pooling of \.{Acme}'s log is started.
Start of the enumeration is syncronized with the start of pulling \.{Acme}'s log.
@c
func main () {
	if len(os.Args)==1 {
		fmt.Fprintf(os.Stderr, "Tag extender\nExtends tags of Acme with specified commands\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [<regular expression>:]<commands> ...\nwhere:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\t<regular expression> - a regular expression applied to window's name\n")
		fmt.Fprintf(os.Stderr, "\t<commands> - a list of commands is added in every Acme's window\n")
		fmt.Fprintf(os.Stderr, "\t\t\tor in windows matched by a specified <regular expression>\n")
		return
	}
	@<Parsing of a command line@>
	sync:=make(chan bool)
	@<Enumerate the opened windows@>
	@<Start polling of \.{Acme}'s log@>
}

@
@<Imports@>=
"strings"
"regexp"
"unicode"

@
@<Global variables@>=
common []string
rgx    map[*regexp.Regexp][]string=make(map[*regexp.Regexp][]string)

@ Let's define an extendend func for fields parsing with escaped symbols and nested quoted strings
@c
func args(s string) []string {
	openeds:=false
	openedd:=false
	escaped:=false
	ff:=func(r rune) bool {
		if !openeds && !openedd && !escaped && unicode.IsSpace(r) {
			return true
		}
		if r=='\\' {
			escaped=!escaped
			return false
		}
		if r=='\'' && !escaped {
			openeds=!openeds
		}

		if r=='"' && !escaped {
			openedd=!openedd
		}
		escaped=false
		return false
	}
	return strings.FieldsFunc(s, ff)
}



@
@<Parsing of a command line@>=
for _,v:=range os.Args[1:] {
	f:=strings.Split(v, ":")
	if len(f)==1 {
		common=append(common, v)
	} else if r,err:=regexp.Compile(f[0]); err!=nil {
		fmt.Fprintf(os.Stderr, "cannot compile regexp %q: %s\n", f[0], err)
	} else {
		rgx[r]=args(f[1])
	}
}

@
@<Imports@>=
"github.com/santucco/goacme"

@
@<Start polling of \.{Acme}'s log@>=
log, err:=goacme.OpenLog()
if err!=nil {
	fmt.Fprint(os.Stderr, err)
	return
}
defer log.Close()
close(sync)
for ev, err:=log.Read(); err==nil; ev, err=log.Read() {
	if ev.Type==goacme.NewWin {
		id:=ev.Id
		name:=ev.Name
		@<Write specified commands to a tag of the new window with |id| after pipe simbol@>
	}
}

@
@<Enumerate the opened windows@>=
go func() {
	<-sync
	ids, err:=goacme.WindowsInfo()
	if err!=nil {
		fmt.Fprintf(os.Stderr, "cannot get a list of the opened windows of Acme: %v\n", err)
		return
	}
	for _, v:=range ids {
		id:=v.Id
		name:=""
		if len(v.Tag)>0 {
			name=v.Tag[0]
		}
		@<Write specified commands to a tag of the new window with |id| after pipe simbol@>
	}
}()

@
@<Write specified commands to a tag of the new window with |id| after pipe simbol@>=
var tag []string
for r,v:=range rgx {
	if r.Match([]byte(name)) {
		tag=append(tag, v...)
	}
}
tag=append(tag, common...)
if err:=writeTag(id, tag); err!=nil {
	fmt.Fprint(os.Stderr, err)
}

@ Let's describe a writing of tag like a function
@c
func writeTag(id int, list []string) error {
	@<Check if |list| is empty@>
	@<Open a window |w| by |id|@>
	@<Read the tag into |s|@>
	@<Remove the tag content before the pipe symbol@>
	@<Compose a new tag@>
	@<Clear the tag and write the new tag@>
	return nil
}

@
@<Check if |list| is empty@>=
if len(list)==0 {
	return nil
}

@
@<Open a window |w| by |id|@>=
w, err:=goacme.Open(id)
if err!=nil {
	return fmt.Errorf("cannot open a window with id %d: %s\n", id, err)
}
defer w.Close()

@
@<Read the tag into |s|@>=
f, err:=w.File("tag")
if err!=nil {
	return fmt.Errorf("cannot get tag file of the window with id %d: %s\n", id, err)
}
var b [200]byte
n, err:=f.Read(b[:])
if err!=nil {
	return fmt.Errorf("cannot read the tag of the window with id %d: %s\n", id, err)
}
s:=string(b[:n])

@
@<Remove the tag content before the pipe symbol@>=
if n=strings.Index(s, "|"); n==-1 {
	n=0
} else {
	n++
}
s=s[n:]

@ We remove duplicates from added command
@<Compose a new tag@>=
{
	for _, v:=range list {
		s=strings.ReplaceAll(s, v, "")
		s=strings.ReplaceAll(s, strings.Trim(v, "\"'"), "")
	}
	list=append(list, strings.Fields(s)...)
	s=" "+strings.Join(list, " ")
}

@
@<Clear the tag and write the new tag@>=
if err:=w.WriteCtl("cleartag"); err!=nil {
	return fmt.Errorf("cannot clear the tag of the window with id %d: %s\n", id, err)
}
if _, err:=f.Write([]byte(s)); err!=nil {
	return fmt.Errorf("cannot write the tag of the window with id %d: %s\n", id, err)
}
