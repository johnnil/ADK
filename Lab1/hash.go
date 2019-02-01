package main

func FirstThree(word []byte) []byte {
	switch len(word) {
	case 1:
		return append([]byte{0x23, 0x23}, word...) // Add # if len(word) < 3
	case 2:
		return append([]byte{0x23}, word...)
	default:
		return word[:3]
	}
}

func Hash(word []byte) uint32 {
	return 900*byteMap(word[0]) + 30*byteMap(word[1]) + byteMap(word[2])
}

func byteMap(b byte) uint32 {
	switch b {
	case 0x23: // a #
		return 0
	case 0xe5: // å
		return 27
	case 0xe4: // ä
		return 28
	case 0xf6: // ö
		return 29
	default: // a-z
		return uint32(b - 96)
	}
}