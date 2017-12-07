/*
  quick linked list.go
  
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
package qll


/*
 * The reason I came up with this system is because I often need 
 * multiple variables to be able to communicate with the same pointer
 * however as "append" always creates a new pointer some links get
 * broken. This way, I can force them to the same pointer
 * 
 * For some reason Go likes the (needlessly) complicated approach on this... 
 * So in the end this didn't work as intended!
 * Still I will keep this in case i find out how to PROPERLY deal with 
 * this.
 * 
 */

import(
	"trickyunits/mkl"
	"fmt"
	)
	
const debugchat = true
func chat(dtext string,s ...interface{}) {
	if debugchat {
		fmt.Printf(dtext+"\n",s)
	}
}
	
type StringList struct {
	list []string
}

func (sl *StringList) Add(s string) {
	sl.list = append(sl.list,s)
	l:=len(sl.list)
	chat("Added '%s' to stringlist",s)
	chat("I know have %d items",l)
}

func (sl *StringList) RemoveIndexes(idx ...int) {
	nl:=make([]string,0)
	for i,v:=range sl.list{
		k:=true
		for _,j:=range(idx){
			k = k && j!=i
		}
		if k {
			nl = append(nl,v)
		}
	}
	sl.list = nl
}

func (sl *StringList) RemoveStrings(str ...string) {
	nl:=make([]string,0)
	for _,v:=range sl.list{
		k:=true
		for _,j:=range(str){
			k = k && j!=v
		}
		if k {
			nl = append(nl,v)
		}
	}
	sl.list = nl
}

func (sl *StringList) Items() []string {
	chat("Returning stringlist with %d items in it",len(sl.list))
	return sl.list
}

func (sl *StringList) Item(i int) string {
	if i<0 || i>=len(sl.list) {
		return "ERROR!"
	} else {
		return sl.list[i]
	}
}

func (sl StringList) Count() int {
	return int(len(sl.list))
}


func CreateStringList() StringList {
	r := StringList{}
	r.list = make([]string,0)
	return r
}

func StringListAddLast(a *StringList,v string) {
	a.list = append(a.list,v)
}

func init() {
mkl.Lic    ("Tricky's Go Units - quick linked list.go","ZLib License")
mkl.Lic    ("Tricky's Go Units - quick linked list.go","ZLib License")
}

