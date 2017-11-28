/*
        jcr6main.go
	(c) 2017 Jeroen Petrus Broks.
	
	This Source Code Form is subject to the terms of the 
	Mozilla Public License, v. 2.0. If a copy of the MPL was not 
	distributed with this file, You can obtain one at 
	http://mozilla.org/MPL/2.0/.
        Version: 17.11.28
*/

package jcr6main

import (
	//"io/ioutil"
	"fmt"
	_ "io"
	"os"
	"strings"
	"trickyunits/mkl"
	"trickyunits/qerr"
	"trickyunits/qff"
	"bytes"
)

var debugchat = true

func chat(s string) {
	if debugchat {
		fmt.Println(s)
	}
}

func chats(f string, a ...interface{}) {
	chat(fmt.Sprintf(f, a))
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
	datastring     map[string]string
	dataint        map[string]int
	databool       map[string]bool
}

type TJCR6Dir struct {
	entries    map[string]TJCR6Entry
	comments   map[string]string
	cfgint     map[string]int32
	cfgbool    map[string]bool
	cfgstr     map[string]string
	fatoffset  int32
	fatsize    int
	fatcsize   int
	fatstorage string
}

type TJCR6Driver struct {
	drvname   string
	recognize func(file string) bool
	dir       func(file string) TJCR6Dir
}

var JCR6Drivers = make(map[string]*TJCR6Driver)

type TJCR6StorageDriver struct {
	Pack   func(b []byte)          []byte
	Unpack func(b []byte,size int) []byte
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

func Entries(J TJCR6Dir) string {
	ret := ""
	for _, v := range J.entries {
		if ret != "" {
			ret += "\n"
		}
		ret += v.entry
	}
	return ret
}

var JCR6Crash bool = true

func jcr6err(em string, p ...interface{}){
	fem:=fmt.Sprintf(em,p)
	fmt.Println("JCR6 Error")
	fmt.Println(fem)
	if JCR6Crash {
		os.Exit(1)
	} else {
		JCR6Error = fem
	}
}

func JCR_B(j TJCR6Dir,entry string) []byte {
	en := strings.ToUpper(entry)
	//var e TJCR6Entry
	if _,ok:= j.entries[en]; !ok{
		jcr6err("Entry %s was not found in the requested resource.",entry)
	}
	e  := j.entries[en]
	pb := make([]byte,e.compressedsize); 
	bt,err := os.Open(e.mainfile)
	if err!=nil {
		jcr6err("Error while opening resource file: %s",e.mainfile)
		return make([]byte,2)
	}
	bt.Seek(int64(e.offset),0)
	bt.Read(pb)
	var ub []byte
	if stdrv,ok:=JCR6StorageDrivers[e.storage];ok{
		ub = stdrv.Unpack(pb,e.size)
	} else {
		jcr6err("Tried to read %s from %s, but the storage algorithm %s does not exist!",entry,e.mainfile,e.storage)
	}
	return ub
}

func JCR_String(j TJCR6Dir,entry string) string {
	return string(JCR_B(j),entry)
}

func init() {
mkl.Version("Tricky's Go Units - jcr6main.go","17.11.28")
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
		ret := TJCR6Dir{} //make(map[string]TJCR6Entry), make(map[string]string), make(map[string]int32), make(map[string]bool), make(map[string]string), 0, 0, 0, "Store"}
		ret.cfgbool = map[string]bool{}
		ret.cfgint = map[string]int32{}
		ret.cfgstr = map[string]string{}
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
		chat(fmt.Sprintf("FAT Offest %i", ret.fatoffset))
		var TTag byte
		var Tag string
		TTag = qff.ReadByte(bt)
		for TTag != 255 {
			Tag = qff.ReadString(bt)
			chat(fmt.Sprintf("cfgtag %i/%s", TTag, Tag))
			switch TTag {
			case 1:
				ret.cfgstr[Tag] = qff.ReadString(bt)
			case 2:
				ret.cfgbool[Tag] = qff.ReadByte(bt) == 1
			case 3:
				ret.cfgint[Tag] = qff.ReadInt32(bt)
			case 255:
			default:
				JCR6Error = "Invalid config tag"
				bt.Close()
				return ret
			}
			TTag = qff.ReadByte(bt)
		}
		if ret.cfgbool["_CaseSensitive"] {
			JCR6Error = "Case Sensitive dir support was already deprecated and removed from JCR6 before it went to the Go language. It's only obvious that support for this was never implemented in Go in the first place."
			bt.Close()
			return ret
		}
		chat("Reading FAT")
		chats("Seeking at: %d", ret.fatoffset)
		qff.Seek(*bt, int(ret.fatoffset))
		chats("Positioned at: %i of %d", qff.Pos(*bt), qff.Size(*bt))
		theend := false
		chats("The End: %s", theend)
		chats("EOF:     %s", qff.EOF(*bt))
		ret.fatsize = qff.ReadInt(bt)
		ret.fatcsize = qff.ReadInt(bt)
		ret.fatstorage = qff.ReadString(bt)
		ret.entries = map[string]TJCR6Entry{}
		chats("FAT Compressed Size: %d",ret.fatcsize)
		chats("FAT True Size:       %d",ret.fatsize)
		chats("FAT Comp. Algorithm: %s",ret.fatstorage)
		
		fatcbytes:=make([]byte,ret.fatcsize)
		var fatbytes []byte
		bt.Read(fatcbytes)
		bt.Close()
		if _,ok:=JCR6StorageDrivers[ret.fatstorage];!ok{
			JCR6Error = fmt.Sprintf("There is no driver found for the %s compression algorithm, so I cannot unpack the file table",ret.fatstorage)
			return ret
		}
		fatbytes=JCR6StorageDrivers[ret.fatstorage].Unpack(fatcbytes,ret.fatsize)
		if len(fatbytes)!=ret.fatsize{
			fmt.Printf("WARNING!!!\nSize after unpacking does NOT match the size written inside the JCR6 file.\nSize is %d and it must be %d\nErrors can be expected!\n",len(fatbytes),ret.fatsize)
		}
		//fatbuffer:=bytes.NewBuffer(fatbytes)
		btf := bytes.NewReader(fatbytes)
		qff.DEOF=false
		for (!qff.DEOF) && (!theend) {
			mtag := qff.ReadByte(btf)
			ppp,_ :=btf.Seek(0,1)
			chats("FAT POSITION %d",ppp)
			chat(fmt.Sprintf("FAT MAIN TAG %d", mtag))
			switch mtag {
			case 0xff:
				theend = true
			case 0x01:
				tag := strings.ToUpper(qff.ReadString(btf))
				chats("FAT TAG %s", tag)
				switch tag {
				case "FILE":
					newentry := TJCR6Entry{}
					newentry.mainfile = file
					newentry.datastring = map[string]string{}
					newentry.dataint = map[string]int{}
					newentry.databool = map[string]bool{}
					ftag := qff.ReadByte(btf)
					for ftag != 255 {
						chats("FILE TAG %d", ftag)
						switch ftag {
						case 1:
							k := qff.ReadString(btf)
							chats("string key %s", k)
							v := qff.ReadString(btf)
							chats("string value %s", v)
							newentry.datastring[k] = v
						case 2:
							kb := qff.ReadString(btf)
							vb := qff.ReadByte(btf) > 0
							chats("boolean key %s", kb)
							chats("boolean value %s",vb)
							newentry.databool[kb] = vb
						case 3:
							ki := qff.ReadString(btf)
							vi := qff.ReadInt32(btf)
							chats("integer key %s",ki)
							chats("integer value %d",vi)
							newentry.dataint[ki] = int(vi)
						case 255:

						default:
							JCR6Error = "Illegal tag"
							bt.Close()
							return ret
						}
						newentry.entry = newentry.datastring["__Entry"]
						newentry.size = newentry.dataint["__Size"]
						newentry.compressedsize = newentry.dataint["__CSize"]
						newentry.offset = newentry.dataint["__Offset"]
						newentry.storage = newentry.datastring["__Storage"]
						newentry.author = newentry.datastring["__Author"]
						newentry.notes = newentry.datastring["__notes"]
						centry := strings.ToUpper(newentry.entry)
						ret.entries[centry] = newentry
						ftag = qff.ReadByte(btf)

					}
				}
			default:
				JCR6Error = fmt.Sprintf("Unknown main tag %d", mtag)
				bt.Close()
				return ret
			}
		}

		bt.Close()
		return ret
	}}

	JCR6StorageDrivers["Store"] = &TJCR6StorageDriver{func(b []byte) []byte {
		return b
	}, func(b []byte,size int) []byte {
		return b
	}}

}
