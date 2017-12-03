/*
  dirry_darwin.go
  
  version: 17.12.03
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
package dirry

import (
	"trickyunits/mkl"
	"path/filepath"
	"os"
)

var predata map[string] string

func init(){
mkl.Version("Tricky's Go Units - dirry_darwin.go","17.12.03")
mkl.Lic    ("Tricky's Go Units - dirry_darwin.go","ZLib License")

predata=map[string] string{}
af,e:=filepath.Abs(os.Args[0])
if e!=nil{
}
predata["AppSupport"] = os.Getenv("HOME") + "/Library/Application Support"
predata["Documents"]  = os.Getenv("HOME") + "/Documents"
predata["AppDir"]     = filepath.Dir(af)
predata["AppFile"]    = af
predata["LinuxDot"]   = ""
predata["Home"]       = os.Getenv("HOME")
predata["User"]       = os.Getenv("USER")


}
