/*
        where-is-all-the-data.go
	(c) 2017 Jeroen Petrus Broks.
	
	This Source Code Form is subject to the terms of the 
	Mozilla Public License, v. 2.0. If a copy of the MPL was not 
	distributed with this file, You can obtain one at 
	http://mozilla.org/MPL/2.0/.
        Version: 17.12.04
*/
package jcr6wad


/*
 * 
 * This file is a translation of the BlitzMax code
 * used to create this driver.
 * 
 * Below is the original license block
 *	JCR6_WAD.bmx
 *	(c) 2015 Jeroen Petrus Broks.
 *	
 *	This Source Code Form is subject to the terms of the 
 *	Mozilla Public License, v. 2.0. If a copy of the MPL was not 
 *	distributed with this file, You can obtain one at 
 *	http://mozilla.org/MPL/2.0/.
 *	Version: 15.09.23
 * 
 * This code translation to Go is done by the same guy ;P
 */
 

import(
	"trickyunits/jcr6/jcr6main"
	"os"
	"strings"
	"trickyunits/qff"
	"trickyunits/qstr"
	"trickyunits/mkl"
)

/*
 * Original license code in BMax. In Go this should be taken to the init() function
 * MKL_Version "JCR6 - JCR6_WAD.bmx","15.09.23"
 * MKL_Lic     "JCR6 - JCR6_WAD.bmx","Mozilla Public License 2.0"
 */
 

// "SupportLevel" means that the maps as they are set in DOOM, Heretic and Hexen are taken a a folder.
// If set to false it will just see them as files giving strange effects, however WAD files not
// set up for this approach (like in the case of Rise of the Triad) might have strange effects also
func fetchWAD (WAD string,SupportLevel bool) jcr6main.TJCR6Dir {
returner:= jcr6main.TJCR6Dir{}
ret:=map[string] jcr6main.TJCR6Entry{} 
returner.Entries = ret
returner.Comments=map[string] string{}
returner.CFGbool=map[string] bool{}
returner.CFGint=map[string] int32{}
returner.CFGstr=map[string] string{}
var e jcr6main.TJCR6Entry
BT,err:= os.Open(WAD)
Level:=""
LevelFiles := []string {"THINGS","LINEDEFS","SIDEDEFS","VERTEXES","SEGS","SSECTORS","NODES","SECTORS","REJECT","BLOCKMAP","BEHAVIOR"} //' All files used in a DOOM/HERETIC/HEXEN level, in which I must note that "BEHAVIOR" is only used in HEXEN.
if err!=nil{
	//'JCRD_DumpError = "JCR_FetchWAD(~q"+WAD+"~q): WAD file could not be read"
	jcr6main.JCR6_JamErr("WAD file could not be read!",WAD,"N/A","JCR6 WAD Driver - JCR_FetchWAD")
	return returner
}
//BT = LittleEndianStream(BT) ' WADs were all written for the MS-DOS platform, which used LittleEndian, so we must make sure that (even if the routine is used on PowerPC Macs) that LittleEndian is used
//'WAD files start with a header that can either be 'IWAD' for main wad files or 'PWAD' for patch WAD files. For JCR this makes no difference it all (it didn't even to WAD for that matter), but this is our only way to check if the WAD loaded is actually a WAD file.
Header := qff.RawReadString(BT,4)  
switch Header {
	case "IWAD":
		returner.Comments["Important notice -- IWAD"] = "The WAD file you are viewing is an IWAD,\nmeaning it belongs to a copyrighted project.\n\nAll content within it is very likely protected by copyright\neither by iD software or Apogee's Developers of Incredible Power or Raven Software.\n\nNothing can stop you from analysing this file and viewing its contents,\nbut don't extract and distribute any contents of this file\nwithout proper permission from the original copyright holder"
	case "PWAD":
		returner.Comments["Notice -- PWAD"] = "This WAD file is a PWAD or Patch-WAD.\nIt's not part of any official file of the games using the WAD system.\nPlease respect the original copyright holders copyrights though!"
	default:
		jcr6main.JCR6_JamErr( "JCR_FetchWAD(\""+WAD+"\"): Requested file is not a WAD file",WAD,"N/A","JCR_FetchWAD")
		return returner
} //End Select	
returner.CFGbool["__CaseSensitive"] = false
//'Next in the WAD files are 2 32bit int values telling how many files the WAD file contains and where in the WAD file the File Table is stored
FileCount := qff.ReadInt(BT)
DirOffset := qff.ReadInt(BT)
returner.FAToffset = int32(DirOffset) // Not that it matters, but hey :P
//DebugLog "This WAD contains "+FileCount+" entries starting at "+DirOffset
BT.Seek(int64(DirOffset),0)
//'And let's now read all the crap
for Ak:=1;Ak<=FileCount;Ak++{
	//'DebugLog "Reading entry #"+Ak
	e = jcr6main.TJCR6Entry{}
	//e.Vars = map[string] string {} //New StringMap ' Just has to be present to prevent crashes in viewer based software.
	e.Mainfile = WAD
	e.Offset = qff.ReadInt(BT)
	e.Size = qff.ReadInt(BT)
	e.Entry = strings.Replace(strings.Trim(qff.RawReadString(BT,8)," \n\x00\r\t"),"\x00","",-123)
	e.Compressedsize = e.Size
	e.Storage = "Store"     //' WAD does not support compression, so always deal this as "Stored"
	//'E.Encryption = 0  ' WAD does not support encryption, so always value 0
	if SupportLevel { //' If set the system will turn DOOM levels into a folder for better usage. When unset the system will just dump everything together with not the best results, but hey, who cares :)
		//'Print "File = "+E.FileName+" >> Level = ~q"+Level+"~q >> Len="+Len(E.FileName)+" >> 1 = "+Left(E.FileName,1)+" >> 3 = "+Mid(E.FileName,3,1)
		//'If Level="" 
		if (qstr.Left(e.Entry,3)=="MAP") {
			Level="MAP_"+e.Entry+"/"
		} else if ((len(e.Entry)==4 && qstr.Left(e.Entry,1)=="E" && qstr.Mid(e.Entry,3,1)=="M")) { 
			Level="MAP_"+e.Entry+"/"
		} else if Level!="" {
			Ok:=false
			for _,S:=range LevelFiles{
				if e.Entry==S {
					Ok=true
				}
				//'Print "Comparing "+E.FileName+" with "+S+"   >>>> "+Ok
				} //Next
			if Ok {
				e.Entry = Level+e.Entry 
			} else { 
				Level=""
			}
		}
	} //EndIf
	//Print "Adding: "+E.FileName	
	ret[strings.ToUpper(e.Entry)] = e //MapInsert Ret,Upper(E.FileName),E
} //Next
BT.Close() //CloseFile BT
//'Return Ret
return returner
} //End Function
//Public


/*
Type DRV_WAD Extends DRV_JCRDIR
	Method Name$()
	Return "Where's All the Data? (WAD)"
	End Method

	Method Recognize(fil$)
	If FileType(fil)<>1 Return False
	Local bt:TStream = ReadFile(fil)
	If Not bt Return False
	Local head$ = ReadString(bt,4)
	CloseFile bt
	Return head="IWAD" Or head="PWAD"
	End Method
	
	Method Dir:TJCRDir(fil$)
	Return JCR_FetchWAD(fil$)
	End Method

    End Type

New DRV_WAD
*/

func init(){
	a:= &jcr6main.TJCR6Driver{ 	Drvname : "Where's All the Data", Recognize : func(wad string) bool {
			if qff.IsDir(wad) { 
				return false
			}
			if !qff.Exists(wad) {
				return false
			}
			bt,err:=os.Open(wad)
			defer bt.Close()
			if err!=nil {
				return false
			}
			head:=qff.RawReadString(bt,4)
			return head=="IWAD" || head=="PWAD"
		}, Dir : func(wad string) jcr6main.TJCR6Dir {
			return fetchWAD(wad,true)
		}}
	jcr6main.JCR6Drivers["WAD"] = a
mkl.Version("Tricky's Go Units - where-is-all-the-data.go","17.12.04")
mkl.Lic    ("Tricky's Go Units - where-is-all-the-data.go","Mozilla Public License 2.0")
}
