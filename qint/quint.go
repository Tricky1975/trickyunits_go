/*
  quint.go

  version: 18.06.12
  Copyright (C) 2018 Jeroen P. Broks
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
package qint


import (
    "encoding/binary"
    "bytes"
)


func Int64toBytes(num int64) ([]byte,error){
  buf1 := new(bytes.Buffer)
  err1 := binary.Write(buf1, binary.LittleEndian, num)
  if err1 != nil {
     return nil,err1
  }
  bs:=buf1.Bytes()
  //fmt.Println( bs )
  return bs,nil
}

func Int32toBytes(num int32) ([]byte,error){
  buf1 := new(bytes.Buffer)
  err1 := binary.Write(buf1, binary.LittleEndian, num)
  if err1 != nil {
     return nil,err1
  }
  bs:=buf1.Bytes()
  //fmt.Println( bs )
  return bs,nil
}

func BytesToInt64(b []byte) (int64,error){
    v := int64(0)
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &v)
  if err != nil {
     return 0,err
  }
  return 0,nil
}

func BytesToInt32(b []byte) (int32,error){
    v := int32(0)
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &v)
  if err != nil {
     return 0,err
  }
  return 0,nil
}


func FloatToBytes(num float64) ([]byte,error){
  buf1 := new(bytes.Buffer)
  err1 := binary.Write(buf1, binary.LittleEndian, num)
  if err1 != nil {
     return nil,err1
  }
  bs:=buf1.Bytes()
  //fmt.Println( bs )
  return bs,nil
}

func BytesToFloat(b []byte) (float64,error){
    v := int64(0)
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, &v)
  if err != nil {
     return 0,err
  }
  return 0,nil
}
