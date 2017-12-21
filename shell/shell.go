/*
  shell.go
  
  version: 17.12.21
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
import "trickyunits/qff"
import "runtime"
import "path/filepath"


var Platform = runtime.GOOS
var ShellError = ""


// Will use the system's command shell in stead of a direct execution
// Especially for some unix built external commands this approach can 
// better. Please note, the system's underlying approach is now
// in order, so differences between systems can be expected
// Especially when using Windows as that target does, unlike Mac and
// Linux, not use a unix-approach!
func Shell(command string){
	shit:=[]string{}
	ShellError=""
	//prog:=""
	for _,p:=range shelldata{
		shit = append(shit,p)
	}
	shit = append(shit,command)
	cmd:= &exec.Cmd{
		Path: shit[0],
		Args: shit,
	}
	if qff.Exists("./"+cmd.Path) {
		cmd.Path = "./"+cmd.Path
	}
	if lp,err:=exec.LookPath(cmd.Path); err!= nil {
		fmt.Println(ansistring.SCol("ERROR!",ansistring.A_Red,ansistring.A_Blink)+"\n"+ansistring.SCol(err.Error(),ansistring.A_Yellow,0))
		ShellError=err.Error()
		os.Exit(50)
	} else {
		cmd.Path = lp
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
	/*
	o,err := cmd.Output()
	outputstring:=fmt.Sprintf("%s",o)
	fmt.Println(outputstring)
	if err!=nil{
		fmt.Println(ansistring.SCol("EXECUTION ERROR!",ansistring.A_Red,ansistring.A_Blink)+"\n"+ansistring.SCol(err.Error(),ansistring.A_Yellow,0))
		os.Exit(51)
	}
	* */
}


func ArrayCommand(name string, argsarray []string) *exec.Cmd{
	ShellError=""
	shit:=[]string{}
	shit = append(shit,name)
	for _,a:=range argsarray { shit = append(shit,a) }
	cmd:= &exec.Cmd{
		Path: name,
		Args: shit,
	}
	if filepath.Base(name) == name {
		if lp, err := exec.LookPath(name); err != nil {
			ShellError = err.Error()
		} else {
			cmd.Path = lp
		}
	}
	return cmd
}
  


func init(){
mkl.Version("Tricky's Go Units - shell.go","17.12.21")
mkl.Lic    ("Tricky's Go Units - shell.go","ZLib License")
}

