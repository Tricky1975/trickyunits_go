/*
  qff.go

  version: 17.11.27
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
package qff

import (
	"encoding/binary"
	"io"
	"os"
	"trickyunits/mkl"
	"trickyunits/qerr"
	"trickyunits/qstr"
)

func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func ReadInt32(f io.Reader) int32 {
	var ret int32 = 0
	err := binary.Read(f, binary.LittleEndian, &ret)
	qerr.QERR(err)
	return ret
}

func RawReadString(f io.Reader, l int32) string {
	ret := make(byte[],l)
	f.Read(ret)
	return qstr.BA2S(ret)
}

func ReadString(f io.reader) string {
	l := ReadInt32(f)
	return RawReadString(f, l)
}

func init() {
	mkl.Version("Tricky's Go Units - qff.go", "17.11.27")
	mkl.Lic("Tricky's Go Units - qff.go", "ZLib License")
}
