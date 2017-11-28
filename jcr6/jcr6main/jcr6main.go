/*
        jcr6main.go
	(c) 2017 Jeroen Petrus Broks.
	
	This Source Code Form is subject to the terms of the 
	Mozilla Public License, v. 2.0. If a copy of the MPL was not 
	distributed with this file, You can obtain one at 
	http://mozilla.org/MPL/2.0/.
        Version: 17.11.28
*/

func Entries(J TJCR6Dir) string {
	ret := ""
	for _, v := range J.entries {
		if ret != "" {
			ret += "\n"
		}
		ret += v.entry
	}
	return ret
}

func init() {
mkl.Version("Tricky's Go Units - jcr6main.go","17.11.28")
mkl.Lic    ("Tricky's Go Units - jcr6main.go","Mozilla Public License 2.0")
	JCR6Drivers["JCR6"] = &TJCR6Driver{"JCR6", func(file string) bool {
		if !qff.Exists(file) {
			chat("File " + file + " does not exist so it cannot be JCR6!")
			return false
		}
		bt, e := os.Open(file)
		if e != nil {
			JCR6Error = e.Error()
			chat("File " + file + " gave the error: " + JCR6Error)
			return false
		}
		head := make([]byte, 5)
		ch1, e := bt.Read(head)
		bt.Close()
		if ch1 != 5 {
			if debugchat {
				fmt.Printf("File %s did not have 5 bytes but had %i instead", file, ch1)

			}
			return false
		}
		chead := make([]byte, 5)
		chead[0] = 74
		chead[1] = 67
		chead[2] = 82
		chead[3] = 54
		chead[4] = 26
		for i := 0; i < 5; i++ {
			if chead[i] != head[i] {
				if debugchat {
					fmt.Printf("Byte %i is %i but had to be %i", i, head[i], chead[i])
				}
				return false
			}
		}
		chat("All is fine for the JCR6 type")
		return true
	}, func(file string) TJCR6Dir {
		JCR6Error = ""
		ret := TJCR6Dir{} //make(map[string]TJCR6Entry), make(map[string]string), make(map[string]int32), make(map[string]bool), make(map[string]string), 0, 0, 0, "Store"}
		ret.cfgbool = map[string]bool{}
		ret.cfgint = map[string]int32{}
		ret.cfgstr = map[string]string{}
		bt, e := os.Open(file)
		qerr.QERR(e)
		if qff.RawReadString(bt, 5) != "JCR6\x1a" {
			panic("YIKES!!! A NONE JCR6 FILE!!!! HELP! HELP! I'M DYING!!!")
		}
		ret.fatoffset = qff.ReadInt32(bt)
		if ret.fatoffset <= 0 {
			JCR6Error = "Invalid FAT offset. Maybe you are trying to read a JCR6 file that has never been properly finalized"
			bt.Close()
			return ret
		}
		chat(fmt.Sprintf("FAT Offest %i", ret.fatoffset))
		var TTag byte
		var Tag string
		TTag = qff.ReadByte(bt)
		for TTag != 255 {
			Tag = qff.ReadString(bt)
			chat(fmt.Sprintf("cfgtag %i/%s", TTag, Tag))
			switch TTag {
			case 1:
				ret.cfgstr[Tag] = qff.ReadString(bt)
			case 2:
				ret.cfgbool[Tag] = qff.ReadByte(bt) == 1
			case 3:
				ret.cfgint[Tag] = qff.ReadInt32(bt)
			case 255:
			default:
				JCR6Error = "Invalid config tag"
				bt.Close()
				return ret
			}
			TTag = qff.ReadByte(bt)
		}
		if ret.cfgbool["_CaseSensitive"] {
			JCR6Error = "Case Sensitive dir support was already deprecated and removed from JCR6 before it went to the Go language. It's only obvious that support for this was never implemented in Go in the first place."
			bt.Close()
			return ret
		}
		chat("Reading FAT")
		chats("Seeking at: %d", ret.fatoffset)
		qff.Seek(*bt, int(ret.fatoffset))
		chats("Positioned at: %i of %d", qff.Pos(*bt), qff.Size(*bt))
		theend := false
		chats("The End: %s", theend)
		chats("EOF:     %s", qff.EOF(*bt))
		ret.fatsize = qff.ReadInt(bt)
		ret.fatcsize = qff.ReadInt(bt)
		ret.fatstorage = qff.ReadString(bt)
		ret.entries = map[string]TJCR6Entry{}
		for (!qff.EOF(*bt)) && (!theend) {
			mtag := qff.ReadByte(bt)
			chat(fmt.Sprintf("FAT MAIN TAG %d", mtag))
			switch mtag {
			case 0xff:
				theend = true
			case 0x01:
				tag := strings.ToUpper(qff.ReadString(bt))
				chats("FAT TAG %s", tag)
				switch tag {
				case "FILE":
					newentry := TJCR6Entry{}
					newentry.mainfile = file
					newentry.datastring = map[string]string{}
					newentry.dataint = map[string]int{}
					newentry.databool = map[string]bool{}
					ftag := qff.ReadByte(bt)
					chats("FILE TAG %d", ftag)
					for ftag != 255 {
						switch ftag {
						case 1:
							k := qff.ReadString(bt)
							chats("string key %s", k)
							v := qff.ReadString(bt)
							chats("string value %s", v)
							newentry.datastring[k] = v
						case 2:
							kb := qff.ReadString(bt)
							vb := qff.ReadByte(bt) > 0
							newentry.databool[kb] = vb
						case 3:
							ki := qff.ReadString(bt)
							vi := qff.ReadInt32(bt)
							newentry.dataint[ki] = int(vi)
						case 255:

						default:
							JCR6Error = "Illegal tag"
							bt.Close()
							return ret
						}
						newentry.entry = newentry.datastring["__Entry"]
						newentry.size = newentry.dataint["__Size"]
						newentry.compressedsize = newentry.dataint["__CSize"]
						newentry.offset = newentry.dataint["__Offset"]
						newentry.storage = newentry.datastring["__Storage"]
						newentry.author = newentry.datastring["__Author"]
						newentry.notes = newentry.datastring["__notes"]
						centry := strings.ToUpper(newentry.entry)
						ret.entries[centry] = newentry
						ftag = qff.ReadByte(bt)

					}
				}
			default:
				JCR6Error = fmt.Sprintf("Unknown main tag %d", mtag)
				bt.Close()
				return ret
			}
		}

		bt.Close()
		return ret
	}}

	JCR6StorageDrivers["Store"] = &TJCR6StorageDriver{func(b []byte) []byte {
		return b
	}, func(b []byte) []byte {
		return b
	}}

}
