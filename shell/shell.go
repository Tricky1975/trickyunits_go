/*
  shell.go
  
  version: 17.12.09
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
package shell

import "os/exec"
import "os"
import "fmt"
import "trickyunits/mkl"
import "trickyunits/ansistring"
import "runtime"


var Platform = runtime.GOOS


// Will use the system's command shell in stead of a direct execution
// Especially for some unix built external commands this approach can 
// better. Please note, the system's underlying approach is now
// in order, so differences between systems can be expected
// Especially when using Windows as that target does, unlike Mac and
// Linux, not use a unix-approach!
func Shell(command string){
	shit:=[]string{}
	//prog:=""
	for _,p:=range shelldata{
		shit = append(shit,p)
	}
	shit = append(shit,command)
	cmd:= &exec.Cmd{
		Path: shit[0],
		Args: shit,
	}
	if lp,err:=exec.LookPath(cmd.Path); err!= nil {
		fmt.Println(ansistring.SCol("ERROR!",ansistring.A_Red,ansistring.A_Blink)+"\n"+ansistring.SCol(err.Error(),ansistring.A_Yellow,0))
		os.Exit(50)
	} else {
		cmd.Path = lp
	}
	o,err := cmd.Output()
	outputstring:=fmt.Sprintf("%s",o)
	fmt.Println(outputstring)
	if err!=nil{
		fmt.Println(ansistring.SCol("EXECUTION ERROR!",ansistring.A_Red,ansistring.A_Blink)+"\n"+ansistring.SCol(err.Error(),ansistring.A_Yellow,0))
		os.Exit(51)
	}
}


func init(){
mkl.Version("Tricky's Go Units - shell.go","17.12.09")
mkl.Lic    ("Tricky's Go Units - shell.go","ZLib License")
}

