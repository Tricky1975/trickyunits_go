/*
        jcr6main.go
	(c) 2017 Jeroen Petrus Broks.
	
	This Source Code Form is subject to the terms of the 
	Mozilla Public License, v. 2.0. If a copy of the MPL was not 
	distributed with this file, You can obtain one at 
	http://mozilla.org/MPL/2.0/.
        Version: 17.12.04
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
	"trickyunits/qstr"
	"trickyunits/dirry"
	"bytes"
	"sort"
	"path/filepath"
)

var debugchat = false
var impdebug = false

func chat(s string) {
	if debugchat {
		fmt.Println(s)
	}
}

func chats(f string, a ...interface{}) {
	chat(fmt.Sprintf(f, a))
}


// Used to store the information of a JCR6 entry
type TJCR6Entry struct {
	Entry          string
	Mainfile       string
	Offset         int
	Size           int
	Compressedsize int
	Storage        string
	Author         string
	Notes          string
	Attrib         int
	Datastring     map[string]string
	Dataint        map[string]int
	Databool       map[string]bool
}

// Used to store the directory inside a JCR6 resource (all patches included)
type TJCR6Dir struct {
	Entries    map[string]TJCR6Entry
	Comments   map[string]string
	Vars       map[string]string
	CFGint     map[string]int32
	CFGbool    map[string]bool
	CFGstr     map[string]string
	FAToffset  int32
	FATsize    int
	FATcsize   int
	FATstorage string
}

// Used to create a driver to make JCR6 recognize other files.
// Only for users who KNOW what they are doring!
type TJCR6Driver struct {
	Drvname   string
	Recognize func(file string) bool
	Dir       func(file string) TJCR6Dir
}

// Used to store all drivers
// Only for users who KNOW what they are doring!
var JCR6Drivers = make(map[string]*TJCR6Driver)

// Used to store all compression methods.
// "Store" is there by default
// Only for users who KNOW what they are doring!
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
		if v.Recognize(file) {
			ret = k
		}
	}
	return ret
}

// Returns the directory of a JCR6 file or a file recognised as such
func Dir(file string) TJCR6Dir {
	t := Recognize(file)
	return JCR6Drivers[t].Dir(file)
}


// Returns a string with all entries inside a JCR6 file 
// The order can be pretty random.
func Entries(J TJCR6Dir) string {
	ret := ""
	for _, v := range J.Entries {
		if ret != "" {
			ret += "\n"
		}
		if v.Entry!=""{
			ret += v.Entry
		}
	}
	// fmt.Printf("returning \"%s\"",ret) // << Debug
	return ret
}

// Returns a list of all entries inside a JCR file as an array.
// The files are sorted by alphabet
func EntryList(J TJCR6Dir) []string{
	r:= strings.Split(Entries(J),"\n")
	sort.Strings(r)
	return r
}

// Entry information (for advanced users).
func Entry(J TJCR6Dir,entry string) TJCR6Entry{
	var ret TJCR6Entry
	var ok bool
	JCR6Error = ""
	if ret,ok=J.Entries[strings.ToUpper(entry)];!ok{
		JCR6Error = "Non-existent entry: "+entry
	}
	return ret
}

// Will crash out your program with exit code 1 if an error occurs.
// If you do not want that to happen, set this to false.
// jcr6main.JCR6Crash = false in your own code will do.
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

// Retreives all content of a JCR6 entry and unpacks it by the
// required algorithm (if the driver for that algorithm is loaded
// by your program that is. :P
func JCR_B(j TJCR6Dir,entry string) []byte {
	en := strings.ToUpper(entry)
	//var e TJCR6Entry
	if _,ok:= j.Entries[en]; !ok{
		jcr6err("Entry %s was not found in the requested resource.",entry)
	}
	e  := j.Entries[en]
	pb := make([]byte,e.Compressedsize); 
	bt,err := os.Open(e.Mainfile)
	if err!=nil {
		jcr6err("Error while opening resource file: %s",e.Mainfile)
		return make([]byte,2)
	}
	bt.Seek(int64(e.Offset),0)
	bt.Read(pb)
	var ub []byte
	if stdrv,ok:=JCR6StorageDrivers[e.Storage];ok{
		ub = stdrv.Unpack(pb,e.Size)
	} else {
		jcr6err("Tried to read %s from %s, but the storage algorithm %s does not exist!",entry,e.Mainfile,e.Storage)
	}
	return ub
}

func Patch(jo *TJCR6Dir,ji TJCR6Dir){
	for k,e := range ji.Entries{
		jo.Entries[k]=e
	}
}

func PatchFile(jo *TJCR6Dir,patch string) bool{
	ji:=Dir(patch)
	if JCR6Error=="" { 
		Patch(jo,ji) 
		return true
	} else {
		if impdebug {
			fmt.Println("ERROR IN PATCHING:\n\t= "+JCR6Error)
		}
		return false
	}
	
}

// Basically the same as JCR_B, but now returns all data as one big string
func JCR_String(j TJCR6Dir,entry string) string {
	return string(JCR_B(j,entry))
}


// Gives the content of a text files line by line.
// Please note, this function has only been set up for the unix "\n" 
// based text files and therefore it will very likely be faulty on Windows.
func JCR_ListEntry(j TJCR6Dir,entry string) []string {
	r:=strings.Split(JCR_String(j,entry),"\n")
	return r
}

func init() {
mkl.Version("Tricky's Go Units - jcr6main.go","17.12.04")
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
		ret.CFGbool = map[string]bool{}
		ret.CFGint = map[string]int32{}
		ret.CFGstr = map[string]string{}
		ret.Vars = map[string]string{}
		ret.Comments = map[string]string{}
		bt, e := os.Open(file)
		qerr.QERR(e)
		if qff.RawReadString(bt, 5) != "JCR6\x1a" {
			panic("YIKES!!! A NONE JCR6 FILE!!!! HELP! HELP! I'M DYING!!!")
		}
		ret.FAToffset = qff.ReadInt32(bt)
		if ret.FAToffset <= 0 {
			JCR6Error = "Invalid FAT offset. Maybe you are trying to read a JCR6 file that has never been properly finalized"
			bt.Close()
			return ret
		}
		chat(fmt.Sprintf("FAT Offest %i", ret.FAToffset))
		var TTag byte
		var Tag string
		TTag = qff.ReadByte(bt)
		for TTag != 255 {
			Tag = qff.ReadString(bt)
			chat(fmt.Sprintf("CFGtag %i/%s", TTag, Tag))
			switch TTag {
			case 1:
				ret.CFGstr[Tag] = qff.ReadString(bt)
			case 2:
				ret.CFGbool[Tag] = qff.ReadByte(bt) == 1
			case 3:
				ret.CFGint[Tag] = qff.ReadInt32(bt)
			case 255:
			default:
				JCR6Error = "Invalid config tag"
				bt.Close()
				return ret
			}
			TTag = qff.ReadByte(bt)
		}
		if ret.CFGbool["_CaseSensitive"] {
			JCR6Error = "Case Sensitive dir support was already deprecated and removed from JCR6 before it went to the Go language. It's only obvious that support for this was never implemented in Go in the first place."
			bt.Close()
			return ret
		}
		chat("Reading FAT")
		chats("Seeking at: %d", ret.FAToffset)
		qff.Seek(*bt, int(ret.FAToffset))
		chats("Positioned at: %i of %d", qff.Pos(*bt), qff.Size(*bt))
		theend := false
		chats("The End: %s", theend)
		chats("EOF:     %s", qff.EOF(*bt))
		ret.FATsize = qff.ReadInt(bt)
		ret.FATcsize = qff.ReadInt(bt)
		ret.FATstorage = qff.ReadString(bt)
		ret.Entries = map[string]TJCR6Entry{}
		chats("FAT Compressed Size: %d",ret.FATcsize)
		chats("FAT True Size:       %d",ret.FATsize)
		chats("FAT Comp. Algorithm: %s",ret.FATstorage)
		
		fatcbytes:=make([]byte,ret.FATcsize)
		var fatbytes []byte
		bt.Read(fatcbytes)
		bt.Close()
		if _,ok:=JCR6StorageDrivers[ret.FATstorage];!ok{
			JCR6Error = fmt.Sprintf("There is no driver found for the %s compression algorithm, so I cannot unpack the file table",ret.FATstorage)
			return ret
		}
		fatbytes=JCR6StorageDrivers[ret.FATstorage].Unpack(fatcbytes,ret.FATsize)
		if len(fatbytes)!=ret.FATsize{
			fmt.Printf("WARNING!!!\nSize after unpacking does NOT match the size written inside the JCR6 file.\nSize is %d and it must be %d\nErrors can be expected!\n",len(fatbytes),ret.FATsize)
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
					newentry.Mainfile = file
					newentry.Datastring = map[string]string{}
					newentry.Dataint = map[string]int{}
					newentry.Databool = map[string]bool{}
					ftag := qff.ReadByte(btf)
					for ftag != 255 {
						chats("FILE TAG %d", ftag)
						switch ftag {
						case 1:
							k := qff.ReadString(btf)
							chats("string key %s", k)
							v := qff.ReadString(btf)
							chats("string value %s", v)
							newentry.Datastring[k] = v
						case 2:
							kb := qff.ReadString(btf)
							vb := qff.ReadByte(btf) > 0
							chats("boolean key %s", kb)
							chats("boolean value %s",vb)
							newentry.Databool[kb] = vb
						case 3:
							ki := qff.ReadString(btf)
							vi := qff.ReadInt32(btf)
							chats("integer key %s",ki)
							chats("integer value %d",vi)
							newentry.Dataint[ki] = int(vi)
						case 255:

						default:
							p,_:=btf.Seek(0,1)
							JCR6Error = fmt.Sprintf("Illegal tag in FILE part %d on fatpos %d",ftag,p)
							bt.Close()
							return ret
						}
						ftag = qff.ReadByte(btf)
					}
					newentry.Entry = newentry.Datastring["__Entry"]
					newentry.Size = newentry.Dataint["__Size"]
					newentry.Compressedsize = newentry.Dataint["__CSize"]
					newentry.Offset = newentry.Dataint["__Offset"]
					newentry.Storage = newentry.Datastring["__Storage"]
					newentry.Author = newentry.Datastring["__Author"]
					newentry.Notes = newentry.Datastring["__notes"]
					centry := strings.ToUpper(newentry.Entry)
					//fmt.Println("Adding entry: ",centry) // <- Debug
					ret.Entries[centry] = newentry
				case "COMMENT":
					commentname:=qff.ReadString(btf)
					ret.Comments[commentname]=qff.ReadString(btf)
				case "IMPORT","REQUIRE":
					if impdebug {
						fmt.Printf("%s request from %s\n",tag,file)
					}
					// Now we're playing with power. Tha ability of 
					// JCR6 to automatically patch other files into 
					// one resource
					deptag := qff.ReadByte(btf)
					var depk,depv string
					depm := map[string] string {}
					for deptag!=255 {
						depk = qff.ReadString(btf)
						depv = qff.ReadString(btf)
						depm[depk] = depv
						deptag = qff.ReadByte(btf)
					}
					depfile  := depm["File"]
					//depsig   := depm["Signature"]
					deppatha := depm["AllowPath"]=="TRUE"
					depcall  := ""
					var depgetpaths [2][] string 
					owndir   := filepath.Dir(file)
					deppath  := 0
					if impdebug{
						fmt.Printf("= Wanted file: %s\n",depfile)
						fmt.Printf("= Allow Path:  %d\n",deppatha)
						fmt.Printf("= ValConv:     %d\n",deppath)
						fmt.Printf("= Prio entnum  %d\n",len(ret.Entries))
					}
					if deppatha {
						deppath=1
					}
					if owndir != "" {
						owndir += "/"
					}
					depgetpaths[0] = append(depgetpaths[0],owndir)
					depgetpaths[1] = append(depgetpaths[1],owndir)
					depgetpaths[1] = append(depgetpaths[1],dirry.Dirry("$AppData$/JCR6/Dependencies/") )
					if qstr.Left(depfile,1)!="/" && qstr.Left(depfile,2)!=":"{
						for _,depdir:=range depgetpaths[deppath]{
							if (depcall=="") && qff.Exists(depdir+depfile) {
								depcall=depdir+depfile
							} else if depcall=="" && impdebug {
								if !qff.Exists(depdir+depfile) {
									fmt.Printf("It seems %s doesn't exist!!\n",depdir+depfile)
								}
							}
						}	
					} else {
						if qff.Exists(depfile) {
							depcall=depfile
						}
					}
					if depcall!="" {
						if (!PatchFile(&ret,depcall)) && tag=="REQUIRE"{
							jcr6err("Required JCR6 addon file ("+depcall+") could not imported!~n~nImporter reported:~n"+JCR6Error) //,fil,"N/A","JCR 6 Driver: Dir()")
							return ret
						} else if tag=="REQUIRE"{
							jcr6err("Required JCR6 addon file ("+depcall+") could not found!") //,fil,"N/A","JCR 6 Driver: Dir()")
						}
					} else if impdebug {
						fmt.Printf("Importing %s failed!",depfile)
						fmt.Printf("Request:    %s",tag)
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
