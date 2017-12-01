/*
  qstr.go
  
  version: 17.12.01
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
	return r
}

func SubStr(a string,pos,length int) string{
	runes := []rune(a)
	endpos:=pos+length
	safeSubstring := string(runes[pos:endpos])
	return safeSubstring
}

// Whoohoo, let's do it the BASIC way :P
func Mid(a string,pos,length int) string{
	return SubStr(a,pos-1,length)
}

func Left(a string,l int) string{
	return SubStr(a,0,l)
}

func Right(a string,l int) string{
	return SubStr(a,len(a)-l,l)
}

// returns -1 if not found at all, otherwise the position number
func FindLast(a string,s string) int{
	for f:=len(a)-len(s);f>0;f--{
		if Mid(a,f,len(s))==s{
			return f
		}
	}
	return -1
}

// Always normal slashes even in Windows file names!
func Slash(s string) string{
	return strings.Replace(s,"\\","/",-1)
}

// strips the extention of a filename.
// Please note, 
func StripExt(file string) string{
	f :=Slash(file)
	lp:=FindLast(f,".")
	ls:=FindLast(f,"/")
	if lp<=1 || ls>lp || lp==ls+1 {
		return f
	}
	return Left(f,lp-1)
}

func StripDir(file string) string{
	f :=Slash(file)
	ls:=FindLast(f,"/")
	if ls==-1{
		return f
	}
	return Right(f,len(f)-ls)

}


func StripAll(file string) string{
	return StripDir(StripExt(file))
}

func init(){
mkl.Lic    ("Tricky's Go Units - qstr.go","ZLib License")
mkl.Version("Tricky's Go Units - qstr.go","17.12.01")
}
