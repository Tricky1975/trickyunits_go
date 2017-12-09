/*
  jcr6lzw.go
  
  version: 17.12.09
  Copyright (C) 2017 Jeroen P. Broks
  This software is provided 'as-is', without any express or implied
  warranty.  In no event will the authors be held liable for any damages
  arising from the use of this software.
  Permission is granted to anyone to use this software for any purpose,
  including commercial applications, and to alter it and redistribute it
  freely, subject to the following restrictions:
  1. The origin of this software must not be misrepresented; you must not
     claim that you wrote the original software. If you use this software
     in a product, an acknowledgment in the product documentation would be
     appreciated but is not required.
  2. Altered source versions must be plainly marked as such, and must not be
     misrepresented as being the original software.
  3. This notice may not be removed or altered from any source distribution.
*/
package jcr6lzw

/* In order to run you must have the lzw library installed
 * go get github.com/itchio/lzw
 * 
 * When this is not done, lzw cannot be supported
 */ 


import (
	"trickyunits/jcr6/jcr6main"
	"trickyunits/mkl"
	"compress/lzw"
	"bytes"
	//"fmt"
)

// These values are for now conventional for this compression method.
// Should other values be needed then add some extra data to the name
// in order to prevent conflicts.
var litw = 8
var order = lzw.LSB

func init() {
mkl.Version("Tricky's Go Units - jcr6lzw.go","17.12.09")
mkl.Lic    ("Tricky's Go Units - jcr6lzw.go","ZLib License")
	jcr6main.JCR6StorageDrivers["lzw"] = &jcr6main.TJCR6StorageDriver{}
	jcr6main.JCR6StorageDrivers["lzw"].Pack = func(b []byte)[]byte{
		var z bytes.Buffer
		/*
		bt,err := zlib.NewWriter(&z)
		if err!=nil{
			JCR6Error = "ZLIB.PACK: "+err.Error()
			return make([]byte)
		}
		*/
		bt := lzw.NewWriter(&z,order,litw)
		bt.Write(b)
		bt.Close()
		return z.Bytes()
	}
	jcr6main.JCR6StorageDrivers["lzw"].Unpack = func(b []byte,size int)[]byte{
		//var z bytes.Buffer = bytes.NewBuffer(b)
		var r []byte
		var b2 []byte = make([]byte,1)
		z:= bytes.NewBuffer(b)
		bti := lzw.NewReader(z,order,litw)
		var err error
		defer bti.Close()
		/*
		if err!=nil{
			jcr6main.JCR6Error = "LZMA.UNPACK: "+err.Error()
			return r
		}
		if err!=nil{
			jcr6main.JCR6Error = "LZMA.UNPACK: "+err.Error()
		}
		*/
		//fi, _ := bti.Stat()
		r = make([]byte,size) //fi.Size())
		for i:=0;i<size;i++{
			_,err=bti.Read(b2)
			r[i]=b2[0]
			// I know this looks pretty amateur, but reading everything 
			// at once causes the data to be truncated, and I simply 
			// cannot allow that to happen.
			// I'll try to investigate how this issue will go on
			// other compression methods once they are being fully
			// implemented.
		}
		if err!=nil && err.Error()!="EOF" {
			jcr6main.JCR6Error = "LZW.UNPACK.R: "+err.Error()
		}
		//li:=-100
		 /* debug
		for i,v := range(r){
			if v==0 {
				if i-li>1{
					fmt.Println("===")
				}
				fmt.Printf("ZERO on %d / %d\n",i,len(r))
				li=i
			}
		}
		// */
		return r
	}
	
}
