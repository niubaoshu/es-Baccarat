package player

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Profile represents a player's persistent data.
type Profile struct {
	Username    string `json:"username"`
	Balance     int    `json:"balance"`
	HandsPlayed int    `json:"hands_played"`
	TotalWager  int    `json:"total_wager"`
}

var ErrPlayerNotFound = errors.New("player profile not found")
var ErrPlayerAlreadyExists = errors.New("player already exists")

const profileDir = "data/profiles"

func getProfilePath(username string) string {
	return filepath.Join(profileDir, fmt.Sprintf("%s.json", username))
}

// LoadProfile attempts to read a player's profile from disk.
func LoadProfile(username string) (*Profile, error) {
	path := getProfilePath(username)
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrPlayerNotFound
		}
		return nil, err
	}
	defer file.Close()

	var p Profile
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

// CreateProfile makes a new profile and saves it. Fails if it already exists.
func CreateProfile(username string, initBalance int) (*Profile, error) {
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		return nil, err
	}

	path := getProfilePath(username)
	if _, err := os.Stat(path); err == nil {
		return nil, ErrPlayerAlreadyExists
	}

	p := &Profile{
		Username: username,
		Balance:  initBalance,
	}

	if err := p.Save(); err != nil {
		return nil, err
	}

	return p, nil
}

// Save writes the current state of the profile to disk.
func (p *Profile) Save() error {
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}

	path := getProfilePath(p.Username)
	return os.WriteFile(path, data, 0644)
}
