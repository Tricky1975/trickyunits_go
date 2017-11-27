/*
        jcr6main.go
	(c) 2017 Jeroen Petrus Broks.
	
	This Source Code Form is subject to the terms of the 
	Mozilla Public License, v. 2.0. If a copy of the MPL was not 
	distributed with this file, You can obtain one at 
	http://mozilla.org/MPL/2.0/.
        Version: 17.11.27
*/
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
	"io"
	"os"
	"trickyunits/mkl"
	"trickyunits/qerr"
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
	entries   map[string]TJCR6Entry
	comments  map[string]string
	cfgint    map[string]int32
	cfgbool   map[string]bool
	cfgstr    map[string]string
	fatoffset int32
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

func Dir(file string) TJCR6Dir {
	t := Recognize(file)
	return JCR6Drivers[t].dir(file)
}

/*
func JOpen(d TJCR6Dir, entry string) io.Reader {
	chat("Opening: " + entry)

}

*/

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
		JCR6Error = ""
		ret := TJCR6Dir{make(map[string]TJCR6Entry), make(map[string]string), make(map[string]int32), make(map[string]bool), make(map[string]string), 0}
		bt, e := os.Open(file)
		qerr.QERR(e)
		if qff.RawReadString(bt, 5) != "JCR6\x1a" {
			panic("YIKES!!! A NONE JCR6 FILE!!!! HELP! HELP! I'M DYING!!!")
		}
		ret.fatoffset = qff.ReadInt32(bt)
		if ret.fatoffset <= 0 {
			JCR6Error = "Invalid FAT offset. Maybe you are trying to read a JCR6 file that has never been properly finalized"
			bt.Close()
			return ret
		}
		var TTag byte
		var Tag string
		TTag = qff.ReadByte(bt)
		for TTag != 255 {
			Tag = qff.ReadString(bt)
			switch TTag {
			case 1:
				ret.cfgstr[Tag] = qff.ReadString(bt)
			case 2:
				ret.cfgbool[Tag] = qff.ReadByte(bt) == 1
			case 3:
				ret.cfgint[Tag] = qff.ReadInt32(bt)
			default:
				JCR6Error = "Invalid config tag"
				bt.Close()
				return ret
			}
		}
		if ret.cfgbool["_CaseSensitive"] {
			JCR6Error = "Case Sensitive dir support was already deprecated and removed from JCR6 before it went to the Go language. It's only obvious that support for this was never implemented in Go in the first place."
			bt.Close()
			return ret
		}

		bt.Close()
		return ret
	}}

	JCR6StorageDrivers["Store"] = &TJCR6StorageDriver{func(b []byte) []byte {
		return b
	}, func(b []byte) []byte {
		return b
	}}

}
