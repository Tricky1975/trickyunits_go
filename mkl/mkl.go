/*
  mkl.go
  
  version: 17.11.29
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

package mkl
import (
	"strconv"
	"strings"
)

var mkl_versions = make(map[string]string)
var mkl_licenses = make(map[string]string)

func Version(n string,v string){
    mkl_versions[n] = v
}

func Lic(n string,l string){
    mkl_licenses[n] = l
}

func ListAll() string {
   ret:=""
   for k, v := range mkl_versions { 
       //fmt.Printf("key[%s] value[%s]\n", k, v)
       ret += k + " ... " + v + " "
       ret += mkl_licenses[k]
       ret += "\n"
   }
   return ret
}


func Newest() string{
	ret:=""
	high:=0
	for _, v := range mkl_versions { 
		a:=strings.Split(v,".")
		i,_:=strconv.Atoi(a[0]+a[1]+a[2])
		if i>high {
			high=i
			ret=v
		}
	}
	return ret
}


/* --
mkl.Version("Tricky's Go Units - mkl.go","17.11.29")
mkl.Lic    ("Tricky's Go Units - mkl.go","ZLib License")
-- */

func init(){
  Version("Tricky's Go Units - mkl.go","17.11.29")
  Lic    ("Tricky's Go Units - mkl.go","ZLib License")
}
