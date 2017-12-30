/*
        jcr6write.go
	(c) 2017 Jeroen Petrus Broks.
	
	This Source Code Form is subject to the terms of the 
	Mozilla Public License, v. 2.0. If a copy of the MPL was not 
	distributed with this file, You can obtain one at 
	http://mozilla.org/MPL/2.0/.
        Version: 17.12.30
*/
package jcr6main

import (
	"log"
	"fmt"
//	"io"
	"os"
	"strings"
	"runtime"
	"trickyunits/qff"
	"trickyunits/mkl"
//	"trickyunits/dirry"
)

func mklwrite(){
mkl.Version("Tricky's Go Units - jcr6write.go","17.12.30")
mkl.Lic    ("Tricky's Go Units - jcr6write.go","Mozilla Public License 2.0")
}



func packdata(input []byte,algorithm string) ([]byte,string) {
	// Try all loaded compression algorithms known
	// The one with the best result will be returned
	if algorithm=="Brute" || algorithm=="BRUTE" {
		retbuf:=input
		retalg:="Store"
		for dname,driver:=range JCR6StorageDrivers {
			tempdat:=driver.Pack(input)
			if len(tempdat)<len(retbuf) {
				retbuf = tempdat
				retalg = dname
			}
		}
		return retbuf,retalg
	} else if algorithm=="Store" || algorithm=="" {
		return input,"Store"
	} else {
		if strings.ToLower(algorithm)!=algorithm {
			log.Print("WARNING! Capital letters found in algorthm name ",algorithm,"\nI'll try to pack, but be aware of possible conflicts in future versions of JCR6")
		}
		// Pack the data with the wanted driver
		if driver,ok:=JCR6StorageDrivers[algorithm];ok {
			ret:=driver.Pack(input)
			// if the packing did indeed reduce the file then keep the 
			// packed data			
			if len(ret)<len(input) {
				return ret,algorithm
			// If the packing did not reduce data or even made the 
			// file bigger then return the original data as "Store"
			} else {
				return input,"Store"
			}
		}
	}
	// Nothing imprinted here should be possible, if it still happens 
	// a warning will show and the data will be returned as "Store"
	log.Print("WARNING! Something went wrong during the packing sequence! This can only be due to a serious bug! I'll return the data as 'Store'")
	return input,"Store"
	
}


type timport struct {
	file string
	sig  string
	kind string
}


// This type is used to create JCR6 files.
// Please note the writer only supports JCR6
type JCR6Create struct{
	Entries    map[string]TJCR6Entry
	Comments   map[string]string
	Vars       map[string]string
	CFGint     map[string]int32
	CFGbool    map[string]bool
	CFGstring  map[string]string
	FATstorage string
	First      bool
	bt         *os.File
	mainfile   string
	imports    []timport
	oof        int64
}

// Configure an integer number in the config.
// Field names prefixed with "__" are reserved
// Please note, this function ONLY works PRIOR to adding the first file
func (jc *JCR6Create) ConfigInt(name string, value int32) {
	if jc.First { return }
	jc.CFGint[name] = value
}

// Configure a boolean value in the config.
// Field names prefixed with "__" are reserved
// Please note, this function ONLY works PRIOR to adding the first file
func (jc *JCR6Create) ConfigBool(name string, value bool) {
	if jc.First { return }
	jc.CFGbool[name] = value
}

// Configure a string in the config.
// Field names prefixed with "__" are reserved
// Please note, this function ONLY works PRIOR to adding the first file
func (jc *JCR6Create) ConfigString(name string, value string) {
	if jc.First { return }
	jc.CFGstring[name] = value
}


func (jc *JCR6Create) saveconfig(){
	jc.First=true
	for k,v:=range jc.CFGstring{
		qff.WriteByte(jc.bt,1)
		qff.WriteString(jc.bt,k)
		qff.WriteString(jc.bt,v)
	}
	for k,v:=range jc.CFGbool{
		qff.WriteByte(jc.bt,2)
		qff.WriteString(jc.bt,k)
		if v { qff.WriteByte(jc.bt,1) } else {qff.WriteByte(jc.bt,0)}
	}
	for k,v:=range jc.CFGint{
		qff.WriteByte(jc.bt,3)
		qff.WriteString(jc.bt,k)
		qff.WriteInt32(jc.bt,v)
	}
	qff.WriteByte(jc.bt,255)
}
// This routine can be used to put an actual block of data into a JCR6
// file as a JCR6 entry. Handy for direct data packing.
// Returns two values. First one is the compressed size of the file, the second is the used storage algorithm.
func (jc *JCR6Create) AddData(data []byte,entryname,algorithm string,filemode int32,timestamp int64,author, notes string) (int32,string){
	// When this is the first data block to be saved, let's first save 
	//the configuration
	if (!jc.First){
		jc.saveconfig()
	}
	packed,storage:=packdata(data,algorithm)
	offs,_:=jc.bt.Seek(0,1)
	ent:=TJCR6Entry{}
	ent.Datastring = map[string]string{}
	ent.Dataint = map[string]int{}
	ent.Databool = map[string]bool{}
	ent.Datastring["__Entry"]=entryname
	ent.Datastring["__TimeStamp"]=fmt.Sprint(timestamp)
	ent.Dataint["__Offset"]=int(offs)
	ent.Dataint["__UnixPermissions"]=int(filemode)
	ent.Datastring["__Storage"]=storage
	ent.Dataint["__Size"]=int(len(data))
	ent.Dataint["__CSize"]=int(len(packed))
	ent.Datastring["__Author"]=author
	ent.Datastring["__Notes"]=notes
	jc.Entries[strings.ToUpper(entryname)] = ent
	written,err:=jc.bt.Write(packed)
	if err!=nil{
		JCR6_JamErr(err.Error(),jc.mainfile,entryname,"AddData")
	}
	return int32(written),storage
}


// Adds a file into the JCR6
// Returns size,compressedsize,storage algorithm
func (jc *JCR6Create) AddFile(originalfile,entryname,algorithm,author,notes string) (int32,int32,string){
	data:=qff.GetFile(originalfile)
	size:=int32(len(data))
	mode:=int32(0777)
	if runtime.GOOS != "windows"{ mode=int32(qff.FileMode(originalfile)) }
	time:=qff.TimeStamp(originalfile)
	csize,alg:=jc.AddData(data,entryname,algorithm,mode,time,author,notes)
	return size,csize,alg
}


// Adds a comment into the JCR6
// Comments are only shown when using JCR6 tools
// When JCR6 files are implemented for data in software they have in normal circumstances no value at all. They are strictly there for documentation purpose!
func (jc *JCR6Create) AddComment(name, comment string){
	jc.Comments[name]=comment
}

// Closes the JCR file and writes its file table
// The data written into the JCR6 file during this process is VITAL. 
// Without it a JCR6 file is 100% useless! (trying to read it will result into an error).
// This function does need a swap file!
func (jc *JCR6Create) Close(){
	if (!jc.First){
		jc.saveconfig()
	}
	offs,_:=jc.bt.Seek(0,1)
	//workdir:=dirry.Dirry("$AppSupport$/$LinuxDot$JCR6G/Create")
	workbas:=jc.mainfile
	i:=1
	for qff.Exists(workbas+"."+fmt.Sprint(i)+".tmp")  { i++ }
	workfat:=workbas+"."+fmt.Sprint(i)+".tmp"
	bt,err:=os.Create(workfat)
	defer bt.Close()
	if err!=nil{
		JCR6_JamErr(err.Error(),jc.mainfile,"<< FILE TABLE >>","<<JCR6CREATE>>.Close()")
		return
	}
	// Dependency call requests
	for _,dependency := range(jc.imports){
		qff.WriteByte(bt,1)
		qff.WriteString(bt,dependency.kind) // This is where either the tag IMPORT or REQUIRE will be written!
		qff.WriteByte(bt,1)
		qff.WriteString(bt,"File")
		qff.WriteString(bt,dependency.file)
		qff.WriteByte(bt,1)
		qff.WriteString(bt,"Signature")
		qff.WriteString(bt,dependency.sig)
		qff.WriteByte(bt,255)
	}
	// Comments
	for k,v := range(jc.Comments){
		qff.WriteByte(bt,1)
		qff.WriteString(bt,"COMMENT")
		qff.WriteString(bt,k)
		qff.WriteString(bt,v)
	}
	// And now, the most important part the entry data
	for _,ent := range(jc.Entries){
		qff.WriteByte(bt,1)
		qff.WriteString(bt,"FILE")
		for k,v:=range ent.Datastring{
			qff.WriteByte(bt,1)
			qff.WriteString(bt,k)
			qff.WriteString(bt,v)
		}
		for k,v:=range ent.Databool{
			qff.WriteByte(bt,2)
			qff.WriteString(bt,k)
			if v { qff.WriteByte(bt,1) } else {qff.WriteByte(bt,0)}
		}
		for k,v:=range ent.Dataint{
			qff.WriteByte(bt,3)
			qff.WriteString(bt,k)
			qff.WriteInt32 (bt,int32(v))
		}
		qff.WriteByte(bt,255)
	}
	qff.WriteByte(bt,255)
	bt.Close()
	fat:=qff.GetFile(workfat)
	packedfat,fatstore:=packdata(fat,jc.FATstorage)
	os.Remove(workfat)
	qff.WriteInt32 (jc.bt,int32(len(fat)))
	qff.WriteInt32 (jc.bt,int32(len(packedfat)))
	qff.WriteString(jc.bt,fatstore)
	jc.bt.Write(packedfat)
	jc.bt.Seek(jc.oof,0)
	//oof,_:=jc.bt.Seek(0,1);
	//fmt.Printf("oof = %d/%d\n",oof,jc.oof) // debug line. oof should be 5 ALWAYS!
	
	qff.WriteInt32(jc.bt,int32(offs))
	jc.bt.Close()
}

func (jc *JCR6Create) addimport(kind,file,sig string){
	jc.imports = append(jc.imports,timport{kind:kind,file:file,sig:sig})
}


func (jc *JCR6Create) AddImport(file,sig string){
	jc.addimport("IMPORT",file,sig)
}

func (jc *JCR6Create) AddRequire(file,sig string){
	jc.addimport("REQUIRE",file,sig)
}
func (jc *JCR6Create) AliasFile(original,target string){
	centryname:=strings.ToUpper(original)
	ctarget:=strings.ToUpper(target)
	if ent,ok:=jc.Entries[strings.ToUpper(centryname)] ; ok {
		newalias:=TJCR6Entry{}
		newalias.Datastring = map[string]string{}
		newalias.Dataint = map[string]int{}
		newalias.Databool = map[string]bool{}
		for k,v:=range ent.Datastring { newalias.Datastring[k] = v }
		for k,v:=range ent.Dataint    { newalias.Dataint   [k] = v }
		for k,v:=range ent.Databool   { newalias.Databool  [k] = v }
		newalias.Datastring["__Entry"]=target
		jc.Entries[ctarget] = newalias
	} else {
		JCR6_JamErr("<jcrcreate>.AliasFile(\""+original+"\",\""+target+"\"): Original not found!",jc.mainfile,"N/A","AliasFile")
	}
}


func JCR_Create(file, FATstorage string) JCR6Create {
	ret:=JCR6Create{
		Entries:map[string]TJCR6Entry{},
		Comments:map[string]string{},
		Vars:map[string]string{},
		CFGstring:map[string]string{},
		CFGbool:map[string]bool{},
		CFGint:map[string]int32{},
		FATstorage:FATstorage,
		First:false,
		mainfile:file,
		imports:[]timport{},
	}
	f, err := os.Create(file)
	if err!=nil {
		JCR6_JamErr(err.Error(),file,"N/A","JCR_Create")
	} else {
		ret.bt=f
		qff.RawWriteString(ret.bt,"JCR6\x1a")
		oof,_ := ret.bt.Seek( 0,1 ) 
		ret.oof = oof
		//fmt.Printf(" oof = %d\n",ret.oof) // debug line... oof should be 5 always.
		qff.WriteInt32(ret.bt,0) // This value will later contain the FAT offset, but that value is not yet known!
		//toof,_ :=ret.bt.Seek(0,1)
		//fmt.Printf("toof = %d\n",toof) // debug line... toof should be 9 always.
		
	}
	return ret
}
