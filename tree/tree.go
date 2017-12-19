/*
  tree.go
  
  version: 17.12.14
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
package tree


import (
	"fmt"
	"os"
	"strings"
//	"path"
	"path/filepath"
	"trickyunits/qstr"
)

func GetTree(rootpath string,hidden bool) []string {

	list := make([]string, 0, 10)
	rp := qstr.Slash(rootpath)
	if qstr.Right(rp,1)!="/" {
	   rp += "/"
   }

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if qstr.Left(filepath.Base(path),1)!="." || hidden {
			ok:=true
			tname:=qstr.Slash(qstr.Right(path,len(path)-len(rp)))
			for _,d := range strings.Split(tname,"/") { ok=ok && (qstr.Left(d,1)!="." || hidden) }
			if ok {
				list = append(list, tname)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("tree.GetTree(\"%s\",%d): walk error [%v]\n", rootpath,hidden,err)
	}
	return list
}

