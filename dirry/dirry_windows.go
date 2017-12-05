/*
  dirry_windows.go
  
  version: 17.12.05
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
mkl.Version("Tricky's Go Units - dirry_windows.go","17.12.05")
mkl.Lic    ("Tricky's Go Units - dirry_windows.go","ZLib License")

predata=map[string] string{}
predata["AppSupport"] = os.Getenv("APPDATA")
predata["Documents"]  = os.Getenv("HOMEPATH") + "/Documents"
predata["AppDir"]     = filepath.Dir(os.Args[0])
predata["AppFile"]    = filepath.Dir(os.Args[0])
predata["LinuxDot"]   = ""
predata["Home"]       = os.Getenv("HOMEPATH")
predata["User"]       = os.Getenv("USER")

}
