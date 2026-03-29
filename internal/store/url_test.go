package store

import (
	"os"
	"path/filepath"
	"testing"
)

func prepareRepository(t *testing.T) *UrlRepository {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "short-links-test.db")
	if err := os.Setenv("DB_NAME_FILE", dbPath); err != nil {
		t.Fatalf("set env: %v", err)
	}

	InitStore()

	repo := GetUrlRepository()
	if err := repo.ClearAll(); err != nil {
		t.Fatalf("clear repo: %v", err)
	}

	t.Cleanup(func() {
		if err := repo.ClearAll(); err != nil {
			t.Fatalf("cleanup repo: %v", err)
		}
	})

	return repo
}

func TestDeleteByShortRemovesBothKeys(t *testing.T) {
	repo := prepareRepository(t)

	created, err := repo.Create("https://example.com/some/page")
	if err != nil {
		t.Fatalf("create url: %v", err)
	}

	if err := repo.DeleteByShort(created.Short); err != nil {
		t.Fatalf("delete by short: %v", err)
	}

	if _, err := repo.FindByShort(created.Short); err == nil {
		t.Fatalf("expected short code to be removed")
	}

	if _, err := repo.FindLink(created.Original); err == nil {
		t.Fatalf("expected original url record to be removed")
	}
}

func TestClearAllRemovesAllCreatedLinks(t *testing.T) {
	repo := prepareRepository(t)

	if _, err := repo.Create("https://example.com/one"); err != nil {
		t.Fatalf("create first url: %v", err)
	}

	if _, err := repo.Create("https://example.com/two"); err != nil {
		t.Fatalf("create second url: %v", err)
	}

	if err := repo.ClearAll(); err != nil {
		t.Fatalf("clear all: %v", err)
	}

	urls, err := repo.ListAll()
	if err != nil {
		t.Fatalf("list all: %v", err)
	}

	if len(urls) != 0 {
		t.Fatalf("expected empty list after clear, got %d items", len(urls))
	}
}
