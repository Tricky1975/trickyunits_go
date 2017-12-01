/*
  ansistring.go
  
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
package ansistring

/* This source is merely a TRANSLATION of code
 * originally written in BlitzMax.
 * I Translated this to Go and made a few adaptions to make
 * it easier to use in Go
 * 
 * Please note that Go may be more professional on BlitzMax in many 
 * fronts, but BlitzMax wins a very important battle. It supports
 * optional function parameters where Go doesn't. 
 * Nearly every coding language supports optional parameters, even
 * the most amateur ones, and Go doesn't... Stability... pah!
 * 
 * BlitzMax is since it's based on BASIC case insensitive.
 * Go is since it's based on C/C++ case sensitive.
 * This can cause conflicts in this translation, but I hope not.
 * 
 * The bbdoc: and about: blocks are for BlitzMax's automated 
 * documentation builder. I'm not quite sure how GoDoc will respond to 
 * this, but I'll keep it here for documentation's sake ;)
 * 
 */

// Strict

//Import tricky_units.MKL_Version
import(
	"trickyunits/mkl"
	"runtime"
	"fmt"
	)
	
	
var ANSI_Use bool = true
func init(){

mkl.Version("Tricky's Go Units - ansistring.go","17.12.01")
mkl.Lic    ("Tricky's Go Units - ansistring.go","ZLib License")

	/*Rem
	bbdoc: When True ANSI String is used. When False, all ANSI functions Return a normal String
	about: On Windows this is by default false, on Linux and Mac this is by default true.
	*/

	if runtime.GOOS=="windows"{
		ANSI_Use = false 
	}
	// I know I should have created a file suffixed with _windows.go for this, however, it's impossible to tell in which order all files are being compiled, so this was the most safe road to go.
	// I also couldn't find an "anything but windows" kind of compilation.
}


const A_Norm      = 0
const A_Bright    = 1
const A_Dark      = 2
const A_Italic    = 3
const A_Underline = 4
const A_Blink     = 5

const A_Black     = 0
const A_Red       = 1
const A_Green     = 2
const A_Yellow    = 3
const A_Blue      = 4
const A_Magenta   = 5
const A_Cyan      = 6
const A_White     = 7


/*Rem
bbdoc: Basic 3 digit ANSI string
returns: The asked string
*/
func String(d1,d2,d3,s string) string{
	if ANSI_Use{
		return fmt.Sprintf("\x1b[%d;%d;%dm%s\x1b[0m",d1,d2,d3,s) //"\x1b["+d1+";"+d2+";"+d3+"m"+s+Chr(27)+"[0m"
	} else {
		return s
	}
} //func

	
/*Rem
bbdoc: Basic color string
about: You can use A_Black, A_Red, A_Green, A_Yellow, A_Blue,A_Magenta,A_Cyan or A_White for color values and A_Norm, A_Bright, A_Dark, A_Underline, A_Blink for flags
returns: The worked out string
*/	
func Col( S string, c1 int, c2 int, flags int) string{
	return ANSI_String(flags,c1+30,C2+40,s)
} //func


/*Rem
bbdoc: Print with only one color
returns: The worked out string
*/
func SCol(s string,col int,flags int)string{
	if ANSI_Use{
		return fmt.Sprintf("\x1b[%d,%dm$m\x1b[0m",flags,col,s)   // "\x1b["+flags+";"+Int(col+30)+"m"+s+Chr(27)+"[0m"
	} else {
		return s
	}
} //func
