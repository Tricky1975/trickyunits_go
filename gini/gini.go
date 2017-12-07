/*
  gini.go
  
  version: 17.12.07
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

// GINI Is Not Ini
// This is a simple configuration form that I like to use. ;)
// This is just a simple parser and saver
package gini

/*
 * This system was originally set up in BlitzMax
 * This source may be for a part a translation, although some parts 
 * will very like be coded a new for better work in Go
 */
 
import(
	"fmt"
	"bytes"
	"strings"
	"trickyunits/qff"
	"trickyunits/qstr"
)
 
const allowedChars  = "qwertyuiopasdfghjklzxcvbnm[]{}1234567890-_+$!@%^&*()_+QWERTYUIOPASDFGHJKL|ZXCVBNM<>?/ '."



// The GINI type
type TGINI struct{
	vars map [string] string
	lists map [string] []string
	init bool
	
} 

func (g *TGINI) init1st(){
	if g.init {
		return
	}
	//fmt.Println("Init new GINI variable")
	//fmt.Printf("before %s\n",g.vars)
	g.init  = true
	g.vars  = make(map[string] string)
	g.lists = map[string] []string{}
	//fmt.Printf("after %s\n",g.vars)
}

// Define var
func (g *TGINI) D(s string,v string) {
	g.init1st()
	//g.vars  = make(map[string] string) // debug!
	g.vars[strings.ToUpper(s)] = v
}

// Read (call) var
func (g *TGINI) C(s string) string{
	g.init1st()
	if v,ok:=g.vars[strings.ToUpper(s)];ok {
		return v
	} else {
		return ""
	}
}

// Creates a list
func (g *TGINI) CL(a string, onlyifnotexist bool) {
	g.init1st()
	if _,ok:=g.lists[strings.ToUpper(a)];ok{
		if onlyifnotexist {
			return
		}
	}
	g.lists[strings.ToUpper(a)] = make([]string,0)
}

// Add value to a list. If not existent create it
func (g *TGINI) Add(nlist string,value string){
	g.CL(nlist,true)
	l:=strings.ToUpper(nlist)
	g.lists[l] = append(g.lists[l],value)
}

// Just returns the list. Creates it if it doesn't yet exist!
func (g *TGINI) List(nlist string) []string{
	g.CL(nlist,true)
	return g.lists[strings.ToUpper(nlist)]
}

// Parses the lines of a text-based GINI file into the GINI data
// Please note this method is for merging data purposes, if you don't
// want to merge, use the regular functions ;)

func (g *TGINI) ParseLines(l []string) {
	// this entire function has been translated from BlitzMax, however the [OLD] tag has been removed.
	g.init1st()
	lst:=make([]string,0)
	tag:=""
	tagsplit:=make([] string,0)
	//tagparam:=make([] string,0)
	tline:=""
	cmd:=""
	para:=""
	pos:=0
	line:=""
	listkeys:=make([]string,0)
	linenumber:=0 // Not present in BMax, but required in go, and makes it even easier of debugging too :P
	for linenumber,line=range l{
		if line!=""{			
			if qstr.Left(qstr.MyTrim(line),1)=="[" && qstr.Right(qstr.MyTrim(line),1)=="]" {
				wTag := qstr.Mid(qstr.MyTrim(line),2,len(qstr.MyTrim(line))-2)
				if strings.ToUpper(wTag)=="OLD"{
					fmt.Printf("ERROR! The [old] tag is NOT supported in this Go version of GINI (and in the original BlitzMax version it's deprecated) in line %d",linenumber)
					return
				}
				tagsplit=strings.Split(wTag,":")
				tag = strings.ToUpper(tagsplit[0])
				if strings.ToUpper(tagsplit[0])=="LIST" {
					if len(tagsplit[0])<2{
						fmt.Println("ERROR! Incorrectly defined list in line %d!",linenumber)
						return
					}
					lst = make([] string,0)
					listkeys=strings.Split(tagsplit[1],",")
					for _,K:=range  listkeys{
						//'ini.clist(UnIniString(K))
						//fmt.Printf("Creating list: %s\n",K)
						g.lists[strings.ToUpper(UnIniString(K))] = lst
					} //Next
					//'lst=ini.list(UnIniString(K))	
				}//EndIf
			} else {
				switch(tag) { //Select tag
				case "REM":
				/* This is the "OLD" tag code. This code is still in the original BlitzMax form and NOT translated.
				 * It's kept for archiving sake, in case the code may still be needed for whatever reason.				
				Case "OLD"
					tline = Trim(line)
					If Left(tline,2)<>"--" 
						tagsplit=tline.split(":")
						If Len(tagsplit)<2 
							Print "Invalid old definition: "+tline
						Else
							If Len(tagsplit)>2 
								For Local i=2 Until Len(tagsplit)
									tagsplit[1]:+":"+tagsplit[i]
									Next
								EndIf
							Print "WARNING! The [old] system has been deprecated and will be removed in future versions! -- "+line
							Select tagsplit[0]
								Case "Var"
									tagparam = tagsplit[1].split("=")
									If Len(tagparam)<2 
										Print "Invalid old var definition: "+Tline
									Else
										For Local ak=0 Until 256 
										    tagparam[1] = Replace(tagparam[1],"%"+Right(Hex(ak),2),Chr(ak))
										    Next
										ini.D(tagparam[0],tagparam[1])
										EndIf
								Case "Add"
									tagparam = tagsplit[1].split(",")
									If Len(tagparam)<2 
										Print "Invalid old var definition: "+Tline
									Else
										ini.Add(tagparam[0],Right(tagsplit[1],Len(tagsplit[1])-(Len(tagparam[0])+1)))
										EndIf
								Case "Dll"		
									tagparam = tagsplit[1].split(",")
									If Len(tagparam)<2 
										Print "Invalid old var definition: "+Tline
									Else
										ini.DuplicateList(tagparam[0],tagparam[1])
										EndIf
								End Select
							EndIf
						EndIf
				*/
				case "SYS","SYSTEM":
					tline = qstr.MyTrim(line)
					pos = strings.IndexAny(tline," ") //tline.find(" ")
					if pos<= -1 {
						pos = len(tline)
					}
					cmd  = strings.ToUpper(tline[:pos])
					para = tline[pos+1:]
					switch( cmd ){
						case "IMPORT","INCLUDE":
							pos = strings.IndexAny(para,"/") //para.find("/")<0
							/*
							?win32
							pos = pos And Chr(para[1])<>":"
							pos = pos And para.find("\")
							?
							*/
							/*
							if pos>0 {
								para=filepath(String(file))+"/"+para
							}
							*/
							/*
							?debug
							Print "Including: "+para
							?
							*/ 
							//g.ReadFile(para) //LoadIni para,ini
							fmt.Printf("Line %d -- WARNING\nNo support yet for file inclusion or importing\n")
						default:
							fmt.Printf("System command %s not understood: %s in line %d\n",cmd,tline,linenumber)
					} //End Select	 
				case "VARS":
					if strings.IndexAny(line,"=")<0 {
						fmt.Printf("Warning! Invalid var definition: %s in line %d\n",line,linenumber)
					} else {
						//tagsplit=line.split("=")
						temppos:=strings.IndexAny(line,"=")
						tagsplit=make([]string,2)
						tagsplit[0]=line[:temppos]
						tagsplit[1]=line[temppos+1:]
						g.D( UnIniString(tagsplit[0]),UnIniString(tagsplit[1]) )
					} //EndIf
				case "LIST":
					lst = append(lst,UnIniString(line)) //ListAddLast lst,uninistring(line)
					for _,K:=range  listkeys{
						g.lists[K]=lst
						//fmt.Printf("[%s] ",K)
					}
					//fmt.Printf("Adding string '%s' to list\n",line)
				case "CALL":
					fmt.Print("WARNING! I cannot execute line %d as the [CALL] block is not supported in Go\n",linenumber)
					/*If line.find(":")<0
						Print "Call: Syntax error: "+line
					Else
						tagsplit=line.split(":")
						inicall tagsplit[0],ini,UnIniString(tagsplit[1])
						EndIf
					*/
				default:
					fmt.Printf("ERROR! Unknown tag: %s (line %d)\n ",tag,linenumber)
					return	
				} //End Select	
			} //EndIf
		} //EndIf		
	} // Next
} //End Function


func (g *TGINI) ReadFromBytes(b []byte){
	// This is a new approach for GINI.
	// The BlitzMax variabt doesn't even support it.
	g.init1st()
	bt:=bytes.NewReader(b)
	head:=qff.RawReadString(bt,5)
	if head!="GINI\x1a" {
		fmt.Println("The buffer read is not a GINI binary")
		return
	}
	for {
		tag:=qff.ReadByte(bt)
		switch(tag){
			case   1:
				k:=qff.ReadString(bt)
				v:=qff.ReadString(bt)
				g.D(k,v)
			case   2:
				cklst:=qff.ReadString(bt)
				g.CL(cklst,false)
			case   3:
				kl:=qff.ReadString(bt)
				kv:=qff.ReadString(bt)
				g.Add(kl,kv)
			case   4:
				list2link:=qff.ReadString(bt)
				list2link2:=qff.ReadString(bt)
				g.lists[list2link]=g.lists[list2link2]
			case 255:
				return
			default:
				fmt.Printf("ERROR! Unknown tag: %d",tag)
				return
			
		}
	} // for
} // func

// The functions below have also been translated from BlitzMax

// It tries to get unwanted characters out, but it's never been fully trustworthy.
// Any ideas to get this fully working are welcome!
func IniString(A string) string {// XAllow been removed
	i:=0
	//Local ret$[] = ["",A]
	ret:=""
	allowed := true
	for i=0;i<len(A);i++{
		allowed = allowed && strings.IndexAny(allowedChars,string(A[i]))>(-1) //(allowedchars+XAllow).find(Chr(A[i]))>=0
		//'If Not allowed Print "I will not allow: "+Chr(A[i])+"/"+A
		ret+=fmt.Sprintf("#(%d)",A[i])
	} //Next
	if allowed {
		return A
	} else {
		return ret
	}
	// Return ret[allowed]	
} //End Function

// Undo the inistring
func UnIniString(A string) string {
	ret:=A
	for i:=0;i<256;i++{
		ret = strings.Replace(ret,fmt.Sprintf("#(%d)",i),string(i),-900)
		//ret = string.Replace(ret,"#u("+i+")",string(i))
	} //Next
	return ret	
} //End Function

func ReadFromLines (lines []string) TGINI{
	ret:=TGINI{}
	ret.init1st()
	ret.ParseLines(lines)
	return ret
}

func ReadFromBytes(thebytes []byte) TGINI{
	ret:=TGINI{}
	ret.init1st()
	ret.ReadFromBytes(thebytes)
	return ret
}


// This function can read a GINI file.
// Either being a 'text' file or 'compiled' file doesn't matter
// this routine can autodetect that.
func ReadFromFile(file string) TGINI{
	var ret TGINI
	if !qff.Exists(file) {
		fmt.Printf("GINI file %f doesn't exist",file)
		return ret
	}
	b:=qff.GetFile(file)
	if string(b[:5])=="GINI\x1a" {
		ret = ReadFromBytes(b)
	} else {
		s:=string(b)
		sl := strings.Split(s,"\n")
		ret = ReadFromLines(sl)
	}
	return ret
}
