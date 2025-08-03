package patterns

func Index(slice []string, item int) string {
	if item >= 0 && item < len(slice) {
		return slice[item]

	}
	return ""
}
