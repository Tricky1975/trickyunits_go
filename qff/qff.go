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

/*
   I need to note, I ONLY deal in LittleEndian. I always did, even when I was on a PPC based Mac, where BigEndian was the standard.
   This most of all since I've always worked for multiple platforms, including Windows (where LittleEndian is the standard).
   I just didn't want conflicts.
   If you really want to use my package yourself and you need BigEndian, well by all means lemme know, it's only
   5 minutes or so to implement this. Pushing that to github will very likely take more time :-P
*/

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

func ReadInt64(f io.Reader) int64 {
	var ret int64 = 0
	err := binary.Read(f, binary.LittleEndian, &ret)
	qerr.QERR(err)
	return ret
}

func RawReadString(f io.Reader, l int32) string {
	ret := make([]byte, l)
	f.Read(ret)
	return qstr.BA2S(ret)
}

func ReadString(f io.Reader) string {
	l := ReadInt32(f)
	return RawReadString(f, l)
}

func ReadByte(r io.Reader) byte {
	buf := make([]byte, 1)
	_, err := r.Read(buf)
	qerr.QERR(err)
	return buf[0]
}

func Seek(r io.Reader, offs int64) {
	r.Seek(offs, 0)
}

func init() {
mkl.Version("Tricky's Go Units - qff.go","17.11.27")
mkl.Lic    ("Tricky's Go Units - qff.go","ZLib License")
}
