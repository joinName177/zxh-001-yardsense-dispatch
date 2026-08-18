package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Backup struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Size      int64     `json:"size"`
}

func (s *JSONStore) Backup(directory string) (Backup, error) {
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return Backup{}, fmt.Errorf("create backup directory: %w", err)
	}
	stamp := time.Now().UTC()
	name := "yardsense-" + stamp.Format("20060102T150405.000000000Z") + ".json"
	path := filepath.Join(directory, name)
	data, err := json.MarshalIndent(s.Snapshot(), "", "  ")
	if err != nil {
		return Backup{}, fmt.Errorf("encode backup: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return Backup{}, fmt.Errorf("write backup: %w", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		return Backup{}, fmt.Errorf("stat backup: %w", err)
	}
	return Backup{Name: name, CreatedAt: stamp, Size: info.Size()}, nil
}

func ListBackups(directory string) ([]Backup, error) {
	entries, err := os.ReadDir(directory)
	if os.IsNotExist(err) {
		return []Backup{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read backup directory: %w", err)
	}
	backups := make([]Backup, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), "yardsense-") || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("read backup info: %w", err)
		}
		backups = append(backups, Backup{Name: entry.Name(), CreatedAt: info.ModTime().UTC(), Size: info.Size()})
	}
	sort.Slice(backups, func(i, j int) bool { return backups[i].CreatedAt.After(backups[j].CreatedAt) })
	return backups, nil
}

func (s *JSONStore) Restore(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read backup: %w", err)
	}
	var document Document
	if err := json.Unmarshal(data, &document); err != nil {
		return fmt.Errorf("decode backup: %w", err)
	}
	return s.Import(document)
}
