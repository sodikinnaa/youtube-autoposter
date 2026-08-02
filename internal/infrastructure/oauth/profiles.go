package oauth

import (
	"fmt"
	"path/filepath"
	"strings"
)

type AccountProfile struct {
	Name        string `json:"name"`
	ProfileName string `json:"profile_name"`
	TokenFile   string `json:"token_file"`
}

// ListProfiles scans current directory for token files (token.json, token_*.json)
func ListProfiles() ([]AccountProfile, error) {
	files, err := filepath.Glob("token*.json")
	if err != nil {
		return nil, err
	}

	var profiles []AccountProfile
	for _, file := range files {
		rawName := strings.TrimSuffix(file, ".json")
		rawName = strings.TrimPrefix(rawName, "token_")
		if rawName == "token" {
			rawName = "default"
		}
		displayName := rawName
		if displayName == "default" {
			displayName = "Default Account"
		} else {
			displayName = strings.Title(strings.ReplaceAll(displayName, "_", " "))
		}
		profiles = append(profiles, AccountProfile{
			Name:        displayName,
			ProfileName: rawName,
			TokenFile:   file,
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
