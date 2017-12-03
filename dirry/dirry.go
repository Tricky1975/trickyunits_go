/*
  dirry.go
  
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


import(
	"trickyunits/mkl"
	"strings"
	//"fmt"
	)
	
var doneinit bool = false

// You can use this to assign values you want dirry to automatically
// replace. You can even overwrite the original system settings,
// however be careful when doing that. 
// If you fear future conflicts with future versions having more
// things set by default (not likely gonna happen, but still)
// put in some kind of trademark to be 100% this will never happen.
// Please note the things you set take priority over the initial settings
var DirryMap map[string] string = map[string] string {}
	
func initdirry(){
	doneinit=true
	for k,v:=range predata{
		if _,ok:=DirryMap[k];!ok{
			DirryMap[k] = v
			// fmt.Printf("Defined %s as %s\n",k,v)
		}
	}
}


// Replace the dirs in accordace with the OS.
func Dirry(s string) string{
	if !doneinit {
		initdirry() // This way I can be 100% all platform specific code has been taken care off!
	}
	ret := s
	for k,v := range DirryMap {
		// fmt.Printf("Replacing %s with %s\n",k,v)
		ret = strings.Replace(ret,"$"+k+"$",v,-1)
		// fmt.Printf("Ret is now '%s'\n",ret)
	}
	return ret
}


func init(){
mkl.Version("Tricky's Go Units - dirry.go","17.12.03")
mkl.Lic    ("Tricky's Go Units - dirry.go","ZLib License")
}
