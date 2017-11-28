/*
        jcr6zlib.go
	(c) 2017 Jeroen Petrus Broks.
	
	This Source Code Form is subject to the terms of the 
	Mozilla Public License, v. 2.0. If a copy of the MPL was not 
	distributed with this file, You can obtain one at 
	http://mozilla.org/MPL/2.0/.
        Version: 17.11.28
*/
package jcr6zlib

import (
	"trickyunits/jcr6/jcr6main"
	"trickyunits/mkl"
	"compress/zlib"
	"bytes"
)

func init() {
mkl.Version("Tricky's Go Units - jcr6zlib.go","17.11.28")
mkl.Lic    ("Tricky's Go Units - jcr6zlib.go","Mozilla Public License 2.0")
	jcr6main.JCR6StorageDrivers["zlib"] = &jcr6main.TJCR6StorageDriver{}
	jcr6main.JCR6StorageDrivers["zlib"].Pack = func(b []byte)[]byte{
		var z bytes.Buffer
		/*
		bt,err := zlib.NewWriter(&z)
		if err!=nil{
			JCR6Error = "ZLIB.PACK: "+err.Error()
			return make([]byte)
		}
		*/
		bt := zlib.NewWriter(&z)
		bt.Write(b)
		bt.Close()
		return z.Bytes()
	}
	jcr6main.JCR6StorageDrivers["zlib"].Unpack = func(b []byte,size int)[]byte{
		//var z bytes.Buffer = bytes.NewBuffer(b)
		var r []byte
		z:= bytes.NewBuffer(b)
		bti, err := zlib.NewReader(z)
		if err!=nil{
			jcr6main.JCR6Error = "ZLIB.UNPACK: "+err.Error()
			return r
		}
		if err!=nil{
			jcr6main.JCR6Error = "ZLIB.UNPACK: "+err.Error()
		}
		//fi, _ := bti.Stat()
		r = make([]byte,size) //fi.Size())
		_,err=bti.Read(r)
		bti.Close()
		return r
	}
	
}
