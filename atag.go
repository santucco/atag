

/*2:*/


//line atag.w:8


//line license:1
// This file is part of atag
//
// Copyright (c) 2020, 2023 Alexander Sychev. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * The name of author may not be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//line atag.w:11

package main

import(


/*3:*/


//line atag.w:23

"fmt"
"os"



/*:3*/



/*5:*/


//line atag.w:48

"strings"
"regexp"
"unicode"



/*:5*/



/*9:*/


//line atag.w:105

"github.com/santucco/goacme"



/*:9*/


//line atag.w:15

)

var(


/*6:*/


//line atag.w:54

common[]string
rgx map[*regexp.Regexp][]string= make(map[*regexp.Regexp][]string)



/*:6*/


//line atag.w:19

)



/*:2*/



/*4:*/


//line atag.w:31

func main(){
if len(os.Args)==1{
fmt.Fprintf(os.Stderr,"Tag extender\nExtends tags of Acme with specified commands\n")
fmt.Fprintf(os.Stderr,"Usage: %s [<regular expression>:]<commands> ...\nwhere:\n",os.Args[0])
fmt.Fprintf(os.Stderr,"\t<regular expression> - a regular expression applied to window's name\n")
fmt.Fprintf(os.Stderr,"\t<commands> - a list of commands is added in every Acme's window\n")
fmt.Fprintf(os.Stderr,"\t\t\tor in windows matched by a specified <regular expression>\n")
return
}


/*8:*/


//line atag.w:88

for _,v:=range os.Args[1:]{
sv:=args(v)
f:=strings.Split(sv[0],":")
if len(f)==1||len(f[0])==0{
common= append(common,v)
}else if r,err:=regexp.Compile(f[0]);err!=nil{
fmt.Fprintf(os.Stderr,"cannot compile regexp %q: %s\n",f[0],err)
}else{
rgx[r]= args(f[1])
if len(sv)> 1{
rgx[r]= append(rgx[r],sv[1:]...)
}
}
}



/*:8*/


//line atag.w:41

sync:=make(chan bool)


/*11:*/


//line atag.w:126

go func(){
<-sync
ids,err:=goacme.WindowsInfo()
if err!=nil{
fmt.Fprintf(os.Stderr,"cannot get a list of the opened windows of Acme: %v\n",err)
return
}
for _,v:=range ids{
id:=v.Id
name:=""
if len(v.Tag)> 0{
name= v.Tag[0]
}


/*12:*/


//line atag.w:145

var tag[]string
for r,v:=range rgx{
if r.Match([]byte(name)){
tag= append(tag,v...)
}
}
tag= append(tag,common...)
if err:=writeTag(id,tag);err!=nil{
fmt.Fprint(os.Stderr,err)
}



/*:12*/


//line atag.w:140

}
}()



/*:11*/


//line atag.w:43



/*10:*/


//line atag.w:109

log,err:=goacme.OpenLog()
if err!=nil{
fmt.Fprint(os.Stderr,err)
return
}
defer log.Close()
close(sync)
for ev,err:=log.Read();err==nil;ev,err= log.Read(){
if ev.Type==goacme.NewWin{
id:=ev.Id
name:=ev.Name


/*12:*/


//line atag.w:145

var tag[]string
for r,v:=range rgx{
if r.Match([]byte(name)){
tag= append(tag,v...)
}
}
tag= append(tag,common...)
if err:=writeTag(id,tag);err!=nil{
fmt.Fprint(os.Stderr,err)
}



/*:12*/


//line atag.w:121

}
}



/*:10*/


//line atag.w:44

}



/*:4*/



/*7:*/


//line atag.w:59

func args(s string)[]string{
openeds:=false
openedd:=false
escaped:=false
ff:=func(r rune)bool{
if!openeds&&!openedd&&!escaped&&unicode.IsSpace(r){
return true
}
if r=='\\'{
escaped= !escaped
return false
}
if r=='\''&&!escaped{
openeds= !openeds
}

if r=='"'&&!escaped{
openedd= !openedd
}
escaped= false
return false
}
return strings.FieldsFunc(s,ff)
}





/*:7*/



/*13:*/


//line atag.w:158

func writeTag(id int,list[]string)error{


/*14:*/


//line atag.w:170

if len(list)==0{
return nil
}



/*:14*/


//line atag.w:160



/*15:*/


//line atag.w:176

w,err:=goacme.Open(id)
if err!=nil{
return fmt.Errorf("cannot open a window with id %d: %s\n",id,err)
}
defer w.Close()



/*:15*/


//line atag.w:161



/*16:*/


//line atag.w:184

f,err:=w.File("tag")
if err!=nil{
return fmt.Errorf("cannot get tag file of the window with id %d: %s\n",id,err)
}
var b[200]byte
n,err:=f.Read(b[:])
if err!=nil{
return fmt.Errorf("cannot read the tag of the window with id %d: %s\n",id,err)
}
s:=string(b[:n])



/*:16*/


//line atag.w:162



/*17:*/


//line atag.w:197

if n= strings.Index(s,"|");n==-1{
n= 0
}else{
n++
}
s= s[n:]



/*:17*/


//line atag.w:163



/*18:*/


//line atag.w:206

{
for _,v:=range list{
s= strings.ReplaceAll(s,v,"")
s= strings.ReplaceAll(s,strings.Trim(v,"\"'"),"")
}
list= append(list,strings.Fields(s)...)
s= " "+strings.Join(list," ")
}



/*:18*/


//line atag.w:164



/*19:*/


//line atag.w:217

if err:=w.WriteCtl("cleartag");err!=nil{
return fmt.Errorf("cannot clear the tag of the window with id %d: %s\n",id,err)
}
if _,err:=f.Write([]byte(s));err!=nil{
return fmt.Errorf("cannot write the tag of the window with id %d: %s\n",id,err)
}

/*:19*/


//line atag.w:165

return nil
}



/*:13*/


