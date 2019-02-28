// License Information:
//   Jeroen's eXtremely Simplistic Repeatative Character Compression Algorithm.go
//   JCR6 -- Jeroen's eXtremely Simplistic Repeatative Character Compression Algorithm
//   version: 19.02.28
//   Copyright (C) 2019 Jeroen P. Broks
//   This software is provided 'as-is', without any express or implied
//   warranty.  In no event will the authors be held liable for any damages
//   arising from the use of this software.
//   Permission is granted to anyone to use this software for any purpose,
//   including commercial applications, and to alter it and redistribute it
//   freely, subject to the following restrictions:
//   1. The origin of this software must not be misrepresented; you must not
//      claim that you wrote the original software. If you use this software
//      in a product, an acknowledgment in the product documentation would be
//      appreciated but is not required.
//   2. Altered source versions must be plainly marked as such, and must not be
//      misrepresented as being the original software.
//   3. This notice may not be removed or altered from any source distribution.
// End License Information
package jxsrcca  // Jeroen's eXtremely Simplistic Repeatative Character Compression Algorithm

/*

   I am thinking to do my fill in creating a DOS game. I really want
   to relive the old days again. JCR6 can play a role in this, however
   as DOS is a bit short on memory and stuff, I wanted a compression 
   algorithm that is very very simplistic, but which I can still 
   include in a DOS game written in Pascal without taking up any 
   extra RAM at all. This method will have to do. The same compression
   algorithm has also been used in my Easy Pack program. It basically
   uses 2 bytes for every byte repeating pattern. So aaa = 3a, that
   sort of way. This will not do for much formats, but when it handles
   DOS games, this will certainly do.

*/



import (
	"trickyunits/jcr6/jcr6main"
	"trickyunits/mkl"
	//"bytes"
	"fmt"
)

var allowchat = false;

func chat(msg string){
	if allowchat { fmt.Print("DEBUG: "+msg) }
}

func init() {	
	mkl.Version("Tricky's Go Units - Jeroen's eXtremely Simplistic Repeatative Character Compression Algorithm.go","19.02.28")
	mkl.Lic    ("Tricky's Go Units - Jeroen's eXtremely Simplistic Repeatative Character Compression Algorithm.go","ZLib License")
	//fmt.Println("JXSRCCA installed!")
	jxsrcca := &jcr6main.TJCR6StorageDriver{}
	jcr6main.JCR6StorageDrivers["jxsrcca"]=jxsrcca;
	
	jxsrcca.Pack = func(b []byte)[]byte{
		ret:= make([]byte,1)
		got:=false
		char:=byte(0)
		reap:=byte(0)
		chat("Let's go!")
		for _,cbyte:=range b{
			chat(fmt.Sprintf("A:%3d/%02X",cbyte,cbyte))
			if (char!=cbyte && got) || reap>250 { 
				ret = append(ret,char)
				ret = append(ret,reap)
				char=cbyte
				reap=1
				chat(fmt.Sprintf("Flushing byte %d x %d (total %9d byte)",char,reap,len(ret)))
			} else if (!got) {
				char=cbyte;
				reap=0;
				got=true;
			} else {
				reap++
			}
		}
		ret = append(ret,char)
		ret = append(ret,reap)
		return ret;
	}
	
	jxsrcca.Unpack = func(b []byte,size int)[]byte{
		ret:= make([]byte,1)
		for p:=1;p<len(b)-1;p+=2{
			ch := byte(b[p])
			rp := byte(b[p+1])
			for i:=byte(0);i<rp;i++ { ret = append(ret,ch) }
		}
		return ret
	}
}
