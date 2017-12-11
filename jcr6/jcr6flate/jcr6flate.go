/*
        jcr6flate.go
	(c) 2017 Jeroen Petrus Broks.
	
	This Source Code Form is subject to the terms of the 
	Mozilla Public License, v. 2.0. If a copy of the MPL was not 
	distributed with this file, You can obtain one at 
	http://mozilla.org/MPL/2.0/.
        Version: 17.12.12
*/
package jcr6flate

import (
	"trickyunits/jcr6/jcr6main"
	"trickyunits/mkl"
	"compress/flate"
	"bytes"
	//"fmt"
)

func init() {
mkl.Version("Tricky's Go Units - jcr6flate.go","17.12.12")
mkl.Lic    ("Tricky's Go Units - jcr6flate.go","Mozilla Public License 2.0")
	jcr6main.JCR6StorageDrivers["flate"] = &jcr6main.TJCR6StorageDriver{}
	jcr6main.JCR6StorageDrivers["flate"].Pack = func(b []byte)[]byte{
		var z bytes.Buffer
		bt,err := flate.NewWriter(&z,9)
		if err!=nil{
			jcr6main.JCR6Error = "FLATE.PACK: "+err.Error()
			return make([]byte,0)
		}
		
		//bt := flate.NewWriter(&z)
		bt.Write(b)
		bt.Close()
		return z.Bytes()
	}
	jcr6main.JCR6StorageDrivers["flate"].Unpack = func(b []byte,size int)[]byte{
		//var z bytes.Buffer = bytes.NewBuffer(b)
		var r []byte
		var b2 []byte = make([]byte,1)
		z:= bytes.NewBuffer(b)
		//bti, err := flate.NewReader(z)
		bti := flate.NewReader(z)
		/*
		if err!=nil{
			jcr6main.JCR6Error = "FLATE.UNPACK: "+err.Error()
			return r
		}
		if err!=nil{
			jcr6main.JCR6Error = "FLATE.UNPACK: "+err.Error()
		}
		*/ 
		//fi, _ := bti.Stat()
		r = make([]byte,size) //fi.Size())
		for i:=0;i<size;i++{
			_,err:=bti.Read(b2)
			r[i]=b2[0]
			// I know this looks pretty amateur, but reading everything 
			// at once causes the data to be truncated, and I simply 
			// cannot allow that to happen.
			// I'll try to investigate how this issue will go on
			// other compression methods once they are being fully
			// implemented.
			if err!=nil && err.Error()!="EOF" {
				jcr6main.JCR6Error = "FLATE.UNPACK.R: "+err.Error()
			}
		}
		bti.Close()
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
