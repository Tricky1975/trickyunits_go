/*
        jcr6realdir.go
	(c) 2017 Jeroen Petrus Broks.
	
	This Source Code Form is subject to the terms of the 
	Mozilla Public License, v. 2.0. If a copy of the MPL was not 
	distributed with this file, You can obtain one at 
	http://mozilla.org/MPL/2.0/.
        Version: 17.12.02
*/



// This driver will allow JCR6 to read a directory as if it were a JCR6 resource file
// You may wonder why you need this. DEBUGGING would be the most valid answer.
// Sometimes it can simply be too much trouble to keep rebuilding the JCR6 resource.
// With this driver you don't have to as you can just read the directory.
// Please note, as this library does not have any function exports it must always be
// prefixed with an underscore when you import it or Go will throw an error!
package jcr6realdir




import "trickyunits/jcr6/jcr6main"
import "trickyunits/qff"
import "trickyunits/qstr"
import "trickyunits/tree"
import "trickyunits/mkl"
import "strings"


func init(){
jcr6main.JCR6Drivers["Real Dir"] = &jcr6main.TJCR6Driver {}



jcr6main.JCR6Drivers["Real Dir"].Recognize = func(file string) bool {
	return qff.IsDir(file)
}

jcr6main.JCR6Drivers["Real Dir"].Dir = func(file string) jcr6main.TJCR6Dir {
	ret := jcr6main.TJCR6Dir{}
	ret.CFGbool = map[string]bool{}
	ret.CFGint = map[string]int32{}
	ret.CFGstr = map[string]string{}
	ret.FATstorage = "OS"
	ret.Entries = map[string]jcr6main.TJCR6Entry{}
	ret.Vars = map[string]string{}
	ret.Comments = map[string]string{}
	ret.Comments ["Real Dir"] = "This is actually not a real JCR6 file\nIt's just the directory "+file+" converted into a JCR6 resource."
	d:=tree.GetTree(file,false)
	dp := qstr.Slash(file)
	if qstr.Right(dp,1)!="/" {
		dp += "/"
	}
	for _,f := range d {
		s:=qff.FileSize(dp+f)
		newentry := jcr6main.TJCR6Entry{}
		newentry.Entry = f
		newentry.Mainfile = dp + f
		newentry.Storage = "Store"
		newentry.Offset  = 0
		newentry.Compressedsize = s
		newentry.Size = s
		centry := strings.ToUpper(f)
		ret.Entries[centry]=newentry
	}
	return ret
}

mkl.Version("Tricky's Go Units - jcr6realdir.go","17.12.02")
mkl.Lic    ("Tricky's Go Units - jcr6realdir.go","Mozilla Public License 2.0")

}
