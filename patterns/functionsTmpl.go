package patterns

func Index(slice []string, item int) string {
	if item >= 0 && item < len(slice) {
		return slice[item]

	}
	return ""
}

func (rt *RootTemplate) GetVariable(key string) string {
	for _, variable := range rt.Variables {
		if variable.Key == key {
			if variable.Default != "" {
				return variable.Default
			}
			// Return empty string or handle missing value
			return ""
		}
	}
	return ""
}
