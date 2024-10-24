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
const deletedDir = "loan_data/deleted/"

type FileRepo struct {
	users   map[string]*domain.User
	deleted map[string]*domain.User
}

func NewFileRepo() (*FileRepo, error) {
	// Load all users' data from files
	users, err := loadUsers()
	if err != nil {
		log.Error().Err(err).Msg("Error loading users")
		return nil, err
	}

	// Load all deleted users' data from files
	deleted, err := loadDeletedUsers()
	if err != nil {
		log.Error().Err(err).Msg("Error loading deleted users")
		return nil, err
	}

	return &FileRepo{
		users:   users,
		deleted: deleted,
	}, nil
}

// GetUser gets an user's data from the map.
func (r *FileRepo) GetUser(userName string) *domain.User {
	return r.users[userName]
}

// AddUser add an user's data to the map.
func (r *FileRepo) AddUser(user *domain.User) error {
	// Add user to the map
	r.users[user.UserName] = user
	return nil
}

// MoveUserToDeleted moves a user's data to the deleted map.
func (r *FileRepo) MoveUserToDeleted(userID string) error {
	// Move user data to deleted map
	r.deleted[userID] = r.users[userID]

	// Remove user from the map
	delete(r.users, userID)
	return nil
}

// PersistUserData saves all user data to files.
func (r *FileRepo) PersistUserData() error {
	// Save all users to files
	for _, user := range r.users {
		err := saveUser(*user)
		if err != nil {
			return err
		}
	}
	r.users = nil

	for _, deleted := range r.deleted {
		err := moveUserToDeleted(deleted.UserName)
		if err != nil {
			return err
		}
	}

	r.deleted = nil
	return nil
}

// loadUsers loads all users' data from files.
func loadUsers() (map[string]*domain.User, error) {
	users := make(map[string]*domain.User)

	// Read all files in the data directory
	files, err := os.ReadDir(dataDir)
	if err != nil {
		log.Error().Err(err).Str("dir", dataDir).Msg("Error reading data directory")
		return users, fmt.Errorf("error reading data directory: %w", err)
	}

	// Load user data from each file
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		userName := file.Name()[:len(file.Name())-5]
		user, err := loadUser(userName, dataDir)
		if err != nil {
			return users, err
		}

		users[user.UserName] = &user
	}

	log.Info().Int("count", len(users)).Msg("All users loaded successfully")
	return users, nil
}

// loadDeletedUsers loads all deleted users' data from files.
func loadDeletedUsers() (map[string]*domain.User, error) {
	users := make(map[string]*domain.User)

	// Read all files in the deleted directory
	files, err := os.ReadDir(deletedDir)
	if err != nil {
		log.Error().Err(err).Str("dir", deletedDir).Msg("Error reading deleted directory")
		return users, fmt.Errorf("error reading deleted directory: %w", err)
	}

	// Load user data from each file
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		userName := file.Name()[:len(file.Name())-5]
		user, err := loadUser(userName, deletedDir)
		if err != nil {
			return users, err
		}

		users[user.UserName] = &user
	}

	log.Info().Int("count", len(users)).Msg("All deleted users loaded successfully")
	return users, nil
}

// loadUser loads a user's data from a file.
func loadUser(userName, dataPath string) (domain.User, error) {
	var user domain.User

	filePath := fmt.Sprintf("%s%s.json", dataPath, userName)
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

// saveUser saves a user's data to a file.
func saveUser(user domain.User) error {
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

// moveUserToDeleted moves the user's JSON file to the deleted directory.
func moveUserToDeleted(userID string) error {
	// Define the paths
	userFilePath := filepath.Join(dataDir, userID+".json")

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
