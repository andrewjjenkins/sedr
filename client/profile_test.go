package client

import (
	"encoding/json"
	"os"
	"testing"
)

func TestParseProfile(t *testing.T) {
	var profile Profile
	profileFile, err := os.Open("test-profile.json")
	if err != nil {
		t.Fatalf("Failed to open test-profile.json: %v", err)
	}

	err = json.NewDecoder(profileFile).Decode(&profile)
	if err != nil {
		t.Fatalf("Failed to parse test-profile.json: %v", err)
	}

	/* This should be an assert instead
	fmt.Printf("Commander: %v\n", profile.Commander)
	fmt.Printf("Ships: %v\n", profile.Ships)
	fmt.Printf("Commodities[0]: %v\n", profile.LastStarport.Commodities[0])
	*/
}
