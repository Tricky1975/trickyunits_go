/*
  gini.go
  
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
	"strings"
)
 
// The GINI type
type TGINI struct{
	vars map [string] string
	lists map [string] []string
	init bool
	
} 

func (g TGINI).init1st(){
	if g.init {
		return
	}
	g.init  = true
	g.vars  = map [string] string{}
	g.lists = map [string] []string
}

// Define var
func (g TGINI).D(s string,v string) {
	g.init1st()
	g.vars[strings.ToUpper(s)] = v
}

// Read (call) var
func (g.TGINI).C(s string) string
	g.init1st()
	if v,ok:=g.vars[strings.ToUpper(s)];!ok {
		return v
	} else {
		return ""
	}
}

// Creates a list
func (g.TGINI).CL(a string, onlyifnotexist bool) {
	if _,ok=g.list[strings.ToUpper(a)];ok{
		if onlyifnotexist {
			return
		}
	}
	g.lists[strings.ToUpper(a)] = make([]string,0)
}

// Add value to a list. If not existent create it
func (g.TGINI).Add(nlist string,value string){
	g.CL(nlist,true)
	l:=strings.ToUpper(nlist)
	g.lists[l] = append(g.lists[l],value)
}

// Just returns the list. Creates it if it doesn't yet exist!
func (g.TGINI).List(nlist string) []string{
	g.CL(nlist,true)
	return g.lists(strings.ToUpper(nlist)]
}
