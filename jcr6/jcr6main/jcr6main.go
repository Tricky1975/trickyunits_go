/*   -- Start License block
        jcr6main.go
	(c) 2017 Jeroen Petrus Broks.
	
	This Source Code Form is subject to the terms of the 
	Mozilla Public License, v. 2.0. If a copy of the MPL was not 
	distributed with this file, You can obtain one at 
	http://mozilla.org/MPL/2.0/.
        Version: 17.11.27
     -- End License block   */

package jcr6main

import (
	//"io/ioutil"
	"fmt"
	"os"
	"trickyunits/mkl"
	"trickyunits/qff"
)

var debugchat = true

func chat(s string) {
	if debugchat {
		fmt.Println(s)
	}
}

type TJCR6Entry struct {
	entry          string
	mainfile       string
	offset         int
	size           int
	compressedsize int
	storage        string
	author         string
	notes          string
	attrib         int
	data           map[string]string
}

type TJCR6Dir struct {
	entries  map[string]TJCR6Entry
	comments map[string]string
}

type TJCR6Driver struct {
	drvname   string
	recognize func(file string) bool
	dir       func(file string) TJCR6Dir
}

var JCR6Drivers = make(map[string]*TJCR6Driver)

type TJCR6StorageDriver struct {
	pack   func(b []byte) []byte
	unpack func(b []byte) []byte
}

var JCR6StorageDrivers = make(map[string]*TJCR6StorageDriver)
var JCR6Error string = ""

// Returns the name of the recognized file type.
// If none are recognized it will return NONE
func Recognize(file string) string {
	ret := "NONE"
	for k, v := range JCR6Drivers {
		chat("Is " + file + " of type " + k + "?")
		//fmt.Printf("key[%s] value[%s]\n", k, v)
		if v.recognize(file) {
			ret = k
		}
	}
	return ret
}

func init() {
mkl.Version("Tricky's Go Units - jcr6main.go","17.11.27")
mkl.Lic    ("Tricky's Go Units - jcr6main.go","Mozilla Public License 2.0")
	JCR6Drivers["JCR6"] = &TJCR6Driver{"JCR6", func(file string) bool {
		if !qff.Exists(file) {
			chat("File " + file + " does not exist so it cannot be JCR6!")
			return false
		}
		bt, e := os.Open(file)
		if e != nil {
			JCR6Error = e.Error()
			chat("File " + file + " gave the error: " + JCR6Error)
			return false
		}
		head := make([]byte, 5)
		ch1, e := bt.Read(head)
		bt.Close()
		if ch1 != 5 {
			if debugchat {
				fmt.Printf("File %s did not have 5 bytes but had %i instead", file, ch1)

			}
			return false
		}
		chead := make([]byte, 5)
		chead[0] = 74
		chead[1] = 67
		chead[2] = 82
		chead[3] = 54
		chead[4] = 26
		for i := 0; i < 5; i++ {
			if chead[i] != head[i] {
				if debugchat {
					fmt.Printf("Byte %i is %i but had to be %i", i, head[i], chead[i])
				}
				return false
			}
		}
		chat("All is fine for the JCR6 type")
		return true
	}, func(file string) TJCR6Dir {
		ret := TJCR6Dir{make(map[string]TJCR6Entry), make(map[string]string)}
		return ret
	}}

	JCR6StorageDrivers["Store"] = &TJCR6StorageDriver{func(b []byte) []byte {
		return b
	}, func(b []byte) []byte {
		return b
	}}

}
