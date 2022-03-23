package utils

func KeyStartWith(source, prefix []byte) bool {
	if len(prefix) > len(source) {
		return false
	}

	for i := 0; i < len(prefix); i++ {
		if prefix[i] != source[i] {
			return false
		}
	}
	return true
}
