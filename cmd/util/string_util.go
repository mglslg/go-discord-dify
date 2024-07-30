package util

func ContainsString(target string, slice []string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
