/*
  qff.go
  
  version: 18.07.24
  Copyright (C) 2017, 2018 Jeroen P. Broks
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
	"crypto/md5"
	"strings"
	"io"
	"io/ioutil"
	"os"
	"fmt"
	"log"
	"errors"
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

var DEOF bool = false

func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func ReadInt32(f io.Reader) int32 {
	var ret int32 = 0
	err := binary.Read(f, binary.LittleEndian, &ret)
	DEOF = err == io.EOF
	qerr.QERR(err)
	return ret
}

func ReadInt(f io.Reader) int {
	return int(ReadInt32(f))
}

func ReadInt64(f io.Reader) int64 {
	var ret int64 = 0
	err := binary.Read(f, binary.LittleEndian, &ret)
	qerr.QERR(err)
	DEOF = err == io.EOF
	return ret
}

func RawReadString(f io.Reader, l int32) string {
	if l<0 { fmt.Printf("WARNING! Negative len poped up: %d \n",l); return "" }
	ret := make([]byte, l)
	_, err := f.Read(ret)
	DEOF = err == io.EOF
	return qstr.BA2S(ret)
}

func WriteInt32(f io.Writer,i int32) {
	err := binary.Write(f, binary.LittleEndian, &i)
	qerr.QERR(err)
}


// Reads a 32 bit int for the string length and then read the string
// based on that data
func ReadString(f io.Reader) string {
	l := ReadInt32(f)
	return RawReadString(f, l)
}

func RawWriteString(f io.Writer, s string) {
	ws:=[]byte(s)
	_,err:=f.Write(ws)
	qerr.QERR(err)
}

func WriteString(f io.Writer,s string) {
	WriteInt32(f,int32(len(s)))
	RawWriteString(f,s)
}

func ReadByte(r io.Reader) byte {
	buf := make([]byte, 1)
	_, err := r.Read(buf)
	DEOF = err == io.EOF
	qerr.QERR(err)
	return buf[0]

}

func WriteByte(w io.Writer,b byte){
	buf := make([]byte,1)
	buf[0]=b
	_,err:=w.Write(buf)
	qerr.QERR(err)
}

func Seek(r os.File, offs int) {
	r.Seek(int64(offs), 0)
}

func Pos(file os.File) int {
	// Find the current position by getting the
	// return value from Seek after moving 0 bytes
	currentPosition, err := file.Seek(0, 1)
	if err != nil {
		panic(err)
	}
	// fmt.Println("Current position:", currentPosition)
	return int(currentPosition)
}

func Size(file os.File) int {
	fi, err := file.Stat()
	if err != nil {
		panic(err)
	}
	return int(fi.Size())
}

// the function of this one is trivial at best
// Go does (for some reasons far beyond anybody who thinks logicall) not support a NORMAL way of EOF detection or languages do, but has a pretty fucked up way of doing this.
// DEOF can help a little, I hope...
func EOF(fi os.File) bool {
	return !(Pos(fi) < Size(fi))
}


// You want to go all the sh.... Go puts us through just to get the filesize?
// NAAH! This routine will do that quickly :P
func FileSize(filename string) int {
     file, err := os.Open(filename)
     if err != nil {
         // handle the error here
         return -1
     }
     defer file.Close()
     // get the file size
     stat, err := file.Stat()
     if err != nil {
       return -2
     }
     return int(stat.Size())
}

func IsDir(filename string) bool {
	file, err := os.Open(filename)
	var ret bool
	if err != nil {
		// handle the error and return
		return false
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
	// handle the error and return
		return false
	}
	if fi.IsDir() {
    // it's a directory
		ret = true
	} else {
		// it's not a directory
		ret = false
	}
	return ret
}

func IsFile(filename string) bool {
	file, err := os.Open(filename)
	var ret bool
	if err != nil {
		// handle the error and return
		return false
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
	// handle the error and return
		return false
	}
	if fi.Mode().IsRegular() {
    // it's a directory
		ret = true
	} else {
		// it's not a directory
		ret = false
	}
	return ret
}

func GetFile(filename string) []byte {
	// Please note... this is not the fastest, but it is the most stable.
	// Files longer than 32767 bytes have shown to get truncated only loading
	// zero-characters after offset 32767 and I simply cannot risk that.
	// Rather a slow routine that works, than a fast one showing trouble.
	size:=FileSize(filename)
	bt,err:=os.Open(filename)
	defer bt.Close()
	if err!=nil{
		fmt.Printf("ERROR!\nGetFile(\"%s\"): %s\n\n",filename,err.Error())
		return make([]byte,size)
	}
	ret:=make([]byte,size)
	b:=make([]byte,1)
	for i:=0;i<size;i++{
		bt.Read(b)
		ret[i]=b[0]
	}
	return ret
}

func EGetFile(filename string) ([]byte,error) {
	// Please note... this is not the fastest, but it is the most stable.
	// Files longer than 32767 bytes have shown to get truncated only loading
	// zero-characters after offset 32767 and I simply cannot risk that.
	// Rather a slow routine that works, than a fast one showing trouble.
	size:=FileSize(filename)
	bt,err:=os.Open(filename)
	defer bt.Close()
	if err!=nil{
		//fmt.Printf("ERROR!\nGetFile(\"%s\"): %s\n\n",filename,err.Error())
		return []byte{},err
	}
	ret:=make([]byte,size)
	b:=make([]byte,1)
	for i:=0;i<size;i++{
		bt.Read(b)
		ret[i]=b[0]
	}
	return ret,err
}

func MergeFiles(source1,source2,target string) error {
	s1:=GetFile(source1)
	s2:=GetFile(source2)
	out, err := os.Create(target)
	defer out.Close()
	if err!=nil { return err }
	out.Write(s1)
	out.Write(s2)
	return nil
}

// Reads entire file as a string
func GetString(filename string) string {
	return string(GetFile(filename))
}

func EGetString(filename string) (string,error) {
	r,e:= EGetFile(filename)
	return string(r),e
}

// Reads a text file and returns it in lines.
// The system does try to detect the difference between a Windows
// text files were line breaks contain both <cr> and <lf> and a unix
// file that only has <lf>
func GetLines(filename string) []string {
	s:=GetString(filename)
	// The two character line break MUST be checked first, or else the
	// one character types will dominate everything causing faulty
	// results.
	if strings.Index(s,"\r\n")>=0 {
		return strings.Split(s,"\r\n")
	} else if strings.Index(s,"\n\r")>=0 {
		return strings.Split(s,"\r\n")
	} else if strings.Index(s,"\n")>=0 {
		return strings.Split(s,"\n")
	} else if strings.Index(s,"\r")>=0 {
		return strings.Split(s,"\r")
	} else {
		return []string{s}
	}
}

func PWD() string {
  dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
  return dir
}

func MD5File(filename string) string{
	// This is adapted code from the original site
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("MD5File(\"%s\"): Error opening file!",filename)
		log.Fatal(err)
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Printf("MD5File(\"%s\"): Error reading data to hash!",filename)
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func FileMode(filename string) int{
    info,_ := os.Stat(filename)
    mode := info.Mode()
    return int(mode)
}

func TimeStamp(filename string) int64{
	info,_:=os.Stat(filename)
	stamp :=info.ModTime()
	return stamp.Unix()
}

func init() {
mkl.Version("Tricky's Go Units - qff.go","18.07.24")
mkl.Lic    ("Tricky's Go Units - qff.go","ZLib License")
}



func WriteStringToFile(filepath, s string) error {
	fo, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fo.Close()

	_, err = io.Copy(fo, strings.NewReader(s))
	if err != nil {
		return err
	}

	return nil
}



// Gets directory and returns it as a listed string
// if t==0 everything
// if t==1 only files
// if t==2 only directories
// Hidden is only read in unix style, so a file being prefixed with a "." or not. If false, hidden files are filtered out!
func GetDir(dir string, t byte, hidden bool) ([]string,error) {
	if t>2 { return []string{},errors.New("Invalid search type") }
	d := strings.Replace(dir,"\\","/",-2)
	files, err := ioutil.ReadDir(d)
	if err != nil {
		return []string{},err
	}

	ret:=[]string{}
	for _, ifile := range files {
		file:=ifile.Name()
		//fmt.Println(file.Name())
		if hidden || qstr.Left(file,1)!="." {
			switch t{
				case 0: ret=append(ret,file)
				case 1: if !ifile.IsDir() { ret=append(ret,file) }
				case 2: if  ifile.IsDir() { ret=append(ret,file) }
			}
		}
	}
	return ret,nil
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
// Source: https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
// By; markc
func CopyFile(src, dst string) (err error) {
    sfi, err := os.Stat(src)
    if err != nil {
        return
    }
    if !sfi.Mode().IsRegular() {
        // cannot copy non-regular files (e.g., directories,
        // symlinks, devices, etc.)
        return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
    }
    dfi, err := os.Stat(dst)
    if err != nil {
        if !os.IsNotExist(err) {
            return
        }
    } else {
        if !(dfi.Mode().IsRegular()) {
            return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
        }
        if os.SameFile(sfi, dfi) {
            return
        }
    }
    if err = os.Link(src, dst); err == nil {
        return
    }
    err = copyFileContents(src, dst)
    return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
    in, err := os.Open(src)
    if err != nil {
        return
    }
    defer in.Close()
    out, err := os.Create(dst)
    if err != nil {
        return
    }
    defer func() {
        cerr := out.Close()
        if err == nil {
            err = cerr
        }
    }()
    if _, err = io.Copy(out, in); err != nil {
        return
    }
    err = out.Sync()
    return
}
