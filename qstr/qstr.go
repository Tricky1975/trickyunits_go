/*
  qstr.go
  
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
package qstr

import(
    "trickyunits/mkl"
    "strconv"
    )
    

// This function was set up by PeterSO on StackOverflow
// https://stackoverflow.com/questions/14230145/what-is-the-best-way-to-convert-byte-array-to-string
func CToGoString(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}

// Quicker way :P
func BA2S(c []byte) string {
	return CToGoString(c[:])
}

func Val(s string) int {
	r,e:=strconv.Atoi(s)
	if e!=nil {
		r=0
	}
	return 0
}

func init(){
mkl.Lic    ("Tricky's Go Units - qstr.go","ZLib License")
mkl.Version("Tricky's Go Units - qstr.go","17.11.29")
}
