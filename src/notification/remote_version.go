package notification

import "strings"

func NormalizeVersion(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	normalized = strings.ReplaceAll(normalized, "_", "-")
	normalized = strings.ReplaceAll(normalized, " ", "")
	return normalized
}

func IsVersionAllowed(currentVersion string, allowed []string) bool {
	target := NormalizeVersion(currentVersion)
	if target == "" {
		return false
	}
	for _, candidate := range allowed {
		if NormalizeVersion(candidate) == target {
			return true
		}
	}
	return false
}
