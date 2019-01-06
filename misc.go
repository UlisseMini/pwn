// misc.go contains random useful functions

package pwn

// ToBytes takes type interface{} and converts it to []byte, if it can't convert
// to []byte it will panic
func ToBytes(t interface{}) (output []byte) {
	switch x := t.(type) {
	case string:
		output = []byte(x)
	case []byte:
		output = x
	case byte:
		output = append(data, x)
	case rune:
		output = []byte(string(x))
	default:
		panic("failed to convert t to type []byte")
	}

	return output
}