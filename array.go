package util

func InStrings(s string, ins []string) bool {
	for _, v := range ins {
		if s == v {
			return true
		}
	}
	return false
}

func AllinStrings(ss []string, ins []string) bool {
	for _, s := range ss {
		if false == InStrings(s, ins) {
			return false
		}
	}
	return true
}
