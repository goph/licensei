package licensematch

// Merge merges two lists of license matches.
func Merge(licenses map[string]float32, newlicenses map[string]float32) map[string]float32 {
	merged := make(map[string]float32)

	// Merge new license matches with the existing list
	for license, confidence := range licenses {
		merged[license] = (confidence + newlicenses[license]) / 2
	}

	// Add new licenses to the list
	for license, confidence := range newlicenses {
		if _, ok := merged[license]; !ok {
			merged[license] = confidence
		}
	}

	return merged
}
