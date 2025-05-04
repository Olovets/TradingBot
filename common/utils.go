package common

func InUintSlice(slice []uint, value uint) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
