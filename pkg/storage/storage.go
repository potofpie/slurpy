package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bobby/slurpy/pkg/models"
)

const (
	SlurpyDir     = ".config/slurpy"
	LogsSubdir    = "logs"
	RequestSuffix = ".json"
)

// Storage manages the persistence of logged requests
type Storage struct {
	baseDir string
}

// New creates a new Storage instance
func New() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	baseDir := filepath.Join(homeDir, SlurpyDir)
	if err := os.MkdirAll(filepath.Join(baseDir, LogsSubdir), 0755); err != nil {
		return nil, fmt.Errorf("failed to create slurpy directories: %w", err)
	}

	return &Storage{baseDir: baseDir}, nil
}

// SaveRequest saves a logged request to disk
func (s *Storage) SaveRequest(req *models.LoggedRequest) error {
	logsDir := filepath.Join(s.baseDir, LogsSubdir)
	filename := fmt.Sprintf("%s_%s%s", req.Namespace, req.ID, RequestSuffix)
	filepath := filepath.Join(logsDir, filename)

	data, err := req.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	return os.WriteFile(filepath, data, 0644)
}

// LoadRequests loads all requests for a given namespace
func (s *Storage) LoadRequests(namespace string) ([]*models.LoggedRequest, error) {
	logsDir := filepath.Join(s.baseDir, LogsSubdir)

	entries, err := os.ReadDir(logsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*models.LoggedRequest{}, nil
		}
		return nil, fmt.Errorf("failed to read logs directory: %w", err)
	}

	var requests []*models.LoggedRequest
	prefix := namespace + "_"

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), prefix) || !strings.HasSuffix(entry.Name(), RequestSuffix) {
			continue
		}

		data, err := os.ReadFile(filepath.Join(logsDir, entry.Name()))
		if err != nil {
			continue // Skip corrupted files
		}

		req, err := models.FromJSON(data)
		if err != nil {
			continue // Skip corrupted files
		}

		requests = append(requests, req)
	}

	// Sort by timestamp, newest first
	sort.Slice(requests, func(i, j int) bool {
		return requests[i].Timestamp.After(requests[j].Timestamp)
	})

	return requests, nil
}

// LoadAllRequests loads all requests from all namespaces
func (s *Storage) LoadAllRequests() ([]*models.LoggedRequest, error) {
	logsDir := filepath.Join(s.baseDir, LogsSubdir)

	entries, err := os.ReadDir(logsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*models.LoggedRequest{}, nil
		}
		return nil, fmt.Errorf("failed to read logs directory: %w", err)
	}

	var requests []*models.LoggedRequest

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), RequestSuffix) {
			continue
		}

		data, err := os.ReadFile(filepath.Join(logsDir, entry.Name()))
		if err != nil {
			continue // Skip corrupted files
		}

		req, err := models.FromJSON(data)
		if err != nil {
			continue // Skip corrupted files
		}

		requests = append(requests, req)
	}

	// Sort by timestamp, newest first
	sort.Slice(requests, func(i, j int) bool {
		return requests[i].Timestamp.After(requests[j].Timestamp)
	})

	return requests, nil
}

// GetNamespaces returns all unique namespaces
func (s *Storage) GetNamespaces() ([]string, error) {
	logsDir := filepath.Join(s.baseDir, LogsSubdir)

	entries, err := os.ReadDir(logsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read logs directory: %w", err)
	}

	namespaces := make(map[string]bool)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), RequestSuffix) {
			continue
		}

		// Extract namespace from filename (namespace_id.json)
		name := entry.Name()
		if idx := strings.LastIndex(name, "_"); idx > 0 {
			namespace := name[:idx]
			namespaces[namespace] = true
		}
	}

	var result []string
	for ns := range namespaces {
		result = append(result, ns)
	}

	sort.Strings(result)
	return result, nil
}

// ClearNamespace removes all requests for a given namespace
func (s *Storage) ClearNamespace(namespace string) error {
	logsDir := filepath.Join(s.baseDir, LogsSubdir)

	entries, err := os.ReadDir(logsDir)
	if err != nil {
		return fmt.Errorf("failed to read logs directory: %w", err)
	}

	prefix := namespace + "_"
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), prefix) && strings.HasSuffix(entry.Name(), RequestSuffix) {
			if err := os.Remove(filepath.Join(logsDir, entry.Name())); err != nil {
				return fmt.Errorf("failed to remove file %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}
