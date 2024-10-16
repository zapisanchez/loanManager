package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zapisanchez/loanMgr/internal/core/domain"

	"github.com/rs/zerolog/log"
)

const dataDir = "loan_data/"

// SaveUser saves a user's data to a file.
func SaveUser(user domain.User) error {
	// Create directory if not exists
	err := os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		log.Error().Err(err).Msg("Error creating data directory")
		return fmt.Errorf("error creating data directory: %w", err)
	}

	// Convert user data to JSON
	userData, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling user data")
		return fmt.Errorf("error marshalling user data: %w", err)
	}

	// Save JSON data to file
	filePath := fmt.Sprintf("%s%s.json", dataDir, user.UserName)
	if err := os.WriteFile(filePath, userData, 0644); err != nil {
		log.Error().Err(err).Str("file", filePath).Msg("Error saving user data to file")
		return err
	}

	log.Info().Str("file", filePath).Msg("User data saved successfully")
	return nil
}

// LoadUser loads a user's data from a file.
func LoadUser(userName string) (domain.User, error) {
	var user domain.User

	filePath := fmt.Sprintf("%s%s.json", dataDir, userName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Error().Err(err).Str("file", filePath).Msg("Error reading user file")
		return user, fmt.Errorf("error reading user file: %w", err)
	}

	// Parse JSON data
	err = json.Unmarshal(data, &user)
	if err != nil {
		log.Error().Err(err).Msg("Error unmarshalling user data")
		return user, fmt.Errorf("error unmarshalling user data: %w", err)
	}

	log.Info().Str("user", userName).Msg("User data loaded successfully")
	return user, nil
}

// MoveUserToDeleted moves the user's JSON file to the deleted directory.
func MoveUserToDeleted(userID string) error {
	// Define the paths
	userFilePath := filepath.Join(dataDir, userID+".json")
	deletedDir := "deleted"

	// Create the deleted directory if it doesn't exist
	if err := os.MkdirAll(deletedDir, os.ModePerm); err != nil {
		return err
	}

	// Move the file to the deleted directory
	newFilePath := filepath.Join(deletedDir, userID+".json")
	err := os.Rename(userFilePath, newFilePath)
	if err != nil {
		return err
	}

	return nil
}
