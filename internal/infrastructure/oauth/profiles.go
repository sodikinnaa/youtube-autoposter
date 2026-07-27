package oauth

import (
	"fmt"
	"path/filepath"
	"strings"
)

type AccountProfile struct {
	Name      string
	TokenFile string
}

// ListProfiles scans current directory for token files (token.json, token_*.json)
func ListProfiles() ([]AccountProfile, error) {
	files, err := filepath.Glob("token*.json")
	if err != nil {
		return nil, err
	}

	var profiles []AccountProfile
	for _, file := range files {
		name := strings.TrimSuffix(file, ".json")
		name = strings.TrimPrefix(name, "token_")
		if name == "token" {
			name = "Default Account"
		} else {
			name = strings.Title(strings.ReplaceAll(name, "_", " "))
		}
		profiles = append(profiles, AccountProfile{
			Name:      name,
			TokenFile: file,
		})
	}
	return profiles, nil
}

// GetTokenFileForProfile returns path for profile name
func GetTokenFileForProfile(profileName string) string {
	if profileName == "" || profileName == "default" {
		return "token.json"
	}
	clean := strings.ToLower(strings.ReplaceAll(profileName, " ", "_"))
	return fmt.Sprintf("token_%s.json", clean)
}
