
// This function was set up by PeterSO on StackOverflow
// https://stackoverflow.com/questions/14230145/what-is-the-best-way-to-convert-byte-array-to-string
func CToGoString(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}

// Quicker way :P
func BA2S(c []byte) string {
	return CToGoString(c[:])
}
   