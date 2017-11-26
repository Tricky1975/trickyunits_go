/*   -- Start License block
        jcr6main.go
	(c) 2017 Jeroen Petrus Broks.

	This Source Code Form is subject to the terms of the
	Mozilla Public License, v. 2.0. If a copy of the MPL was not
	distributed with this file, You can obtain one at
	http://mozilla.org/MPL/2.0/.
        Version: 17.11.26
     -- End License block   */

package jcr6main

import (
	"trickyunits/mkl"
)

type TJCR6Entry struct {
	entry          string
	mainfile       string
	offset         int
	size           int
	compressedsize int
	storage        string
	author         string
	notes          string
	attrib         int
	data           map[string]string
}

type TJCR6Dir struct {
	entries  map[string]TJCR6Entry
	comments map[string]string
}

type TJCR6Driver struct {
	drvname   string
	recognize func(file string) bool
	dir       func(file string) TJCR6Dir
}

var JCR6Drivers = make(map[string]*TJCR6Driver)

type TJCR6StorageDriver struct {
	pack   func(b []byte) []byte
	unpack func(b []byte) []byte
}

var JCR6StorageDrivers = make(map[string]*TJCR6StorageDriver)

func init() {
	mkl.Version("Tricky's Go Units - jcr6main.go", "17.11.26")
	mkl.Lic("Tricky's Go Units - jcr6main.go", "Mozilla Public License 2.0")
	JCR6Drivers["JCR6"] = &TJCR6Driver{"JCR6", func(file string) bool {
		return true
	}, func(file string) TJCR6Dir {
		ret := TJCR6Dir{make(map[string]TJCR6Entry), make(map[string]string)}
		return ret
	}}

	JCR6StorageDrivers["Store"] = &TJCR6StorageDriver{func(b []byte) []byte {
		return b
	}, func(b []byte) []byte {
		return b
	}}

}
