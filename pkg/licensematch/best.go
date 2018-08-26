package licensematch

// Best returns the first best match from a license list.
//
// Note: thanks to the unordered nature of maps "first" means first in alphabetical order.
// This is to ensure consistent results for this function.
func Best(licenses map[string]float32) (license string, confidence float32) {
	for l, c := range licenses {
		if c > confidence {
			license = l
			confidence = c
		} else if c == confidence && license > l {
			license = l
		}
	}

	return
}

// Bests returns the all best matches and their confidence level from a license list.
func Bests(licenses map[string]float32) (matches map[string]float32, maxConfidence float32) {
	for l, c := range licenses {
		if c > maxConfidence {
			matches = make(map[string]float32)
			maxConfidence = c
		}

		if c == maxConfidence {
			matches[l] = c
		}
	}

	return
}
