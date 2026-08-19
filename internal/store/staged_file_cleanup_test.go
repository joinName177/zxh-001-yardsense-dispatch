package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/model"
)

func TestFailedDurableWritesRemoveStagedFiles(t *testing.T) {
	root := t.TempDir()
	statePath := filepath.Join(root, "state.json")
	database, err := Open(statePath)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	if err := os.Mkdir(statePath, 0o755); err != nil {
		t.Fatalf("create blocked state destination: %v", err)
	}
	worker := model.Worker{ID: "worker-cleanup", Name: "Cleanup Worker", Zones: []string{"north-yard"}, Skills: []string{"safety"}, Active: true, Capacity: 1}
	if err := database.UpsertWorker(worker); err == nil {
		t.Fatal("save worker unexpectedly succeeded")
	}
	if _, err := os.Stat(statePath + ".tmp"); !os.IsNotExist(err) {
		t.Fatalf("failed state write left staged file: %v", err)
	}

	backupDir := filepath.Join(root, "backups")
	stamp := time.Date(2026, 8, 19, 13, 0, 0, 0, time.UTC)
	backupName := "yardsense-" + stamp.Format("20060102T150405.000000000Z") + ".json"
	if err := os.MkdirAll(filepath.Join(backupDir, backupName), 0o755); err != nil {
		t.Fatalf("create blocked backup destination: %v", err)
	}
	if _, err := database.backupAt(backupDir, stamp); err == nil {
		t.Fatal("backup unexpectedly succeeded")
	}
	if _, err := os.Stat(filepath.Join(backupDir, backupName+".partial")); !os.IsNotExist(err) {
		t.Fatalf("failed backup left staged file: %v", err)
	}
}
