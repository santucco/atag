

/*2:*/


//line atag.w:8


//line license:1

// This file is part of atag version 0.1
//
// Copyright (c) 2020 Alexander Sychev. All rights reserved.
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


//line atag.w:19

"fmt"
"os"



/*:3*/



/*5:*/


//line atag.w:39

"github.com/santucco/goacme"



/*:5*/



/*9:*/


//line atag.w:80

"strings"



/*:9*/


//line atag.w:15

)



/*:2*/



/*4:*/


//line atag.w:27

func main(){
if len(os.Args)==1{
fmt.Fprintf(os.Stderr,"Tag extender\nExtends tags of Acme with specified commands\nUsage: %s <commands>\n",os.Args[0])
return
}
sync:=make(chan bool)


/*7:*/


//line atag.w:59

go func(){
<-sync
ids,err:=goacme.WindowsInfo()
if err!=nil{
fmt.Fprintf(os.Stderr,"cannot get a list of the opened windows of Acme: %v\n",err)
return
}
for _,v:=range ids{
id:=v.Id


/*8:*/


//line atag.w:74

if err:=writeTag(id,os.Args[1:]);err!=nil{
fmt.Fprint(os.Stderr,err)
}



/*:8*/


//line atag.w:69

}
}()



/*:7*/


//line atag.w:34



/*6:*/


//line atag.w:43

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


/*8:*/


//line atag.w:74

if err:=writeTag(id,os.Args[1:]);err!=nil{
fmt.Fprint(os.Stderr,err)
}



/*:8*/


//line atag.w:54

}
}



/*:6*/


//line atag.w:35

}



/*:4*/



/*10:*/


//line atag.w:84

func writeTag(id int,list[]string)error{


/*11:*/


//line atag.w:96

if len(list)==0{
return nil
}



/*:11*/


//line atag.w:86



/*12:*/


//line atag.w:102

w,err:=goacme.Open(id)
if err!=nil{
return fmt.Errorf("cannot open a window with id %d: %s\n",id,err)
}
defer w.Close()



/*:12*/


//line atag.w:87



/*13:*/


//line atag.w:110

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



/*:13*/


//line atag.w:88



/*14:*/


//line atag.w:123

if n= strings.Index(s,"|");n==-1{
n= 0
}else{
n++
}
s= s[n:]



/*:14*/


//line atag.w:89



/*15:*/


//line atag.w:132

s= " "+strings.Join(list," ")+s



/*:15*/


//line atag.w:90



/*16:*/


//line atag.w:136

if err:=w.WriteCtl("cleartag");err!=nil{
return fmt.Errorf("cannot clear the tag of the window with id %d: %s\n",id,err)
}
if _,err:=f.Write([]byte(s));err!=nil{
return fmt.Errorf("cannot write the tag of the window with id %d: %s\n",id,err)
}

/*:16*/


//line atag.w:91

return nil
}



/*:10*/


