package store

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"
	"urlShortener/pkg/db"

	"urlShortener/internal/model"

	"go.etcd.io/bbolt"
)

type UrlRepository struct {
	repo *db.Repository
	ttl  time.Duration
}

func GetUrlRepository() *UrlRepository {
	return &UrlRepository{
		repo: db.GetRepository(),
		ttl:  3 * 24 * time.Hour, // 3 дня
	}
}

func (ur *UrlRepository) ListAll() ([]model.Url, error) {
	result := make([]model.Url, 0)

	err := ur.repo.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(UrlBucket)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var u model.Url
			if err := json.Unmarshal(v, &u); err != nil {
				continue
			}
			// Чтобы скипать шорты
			if string(k) == u.Original {
				result = append(result, u)
			}

		}
		return nil
	})

	return result, err
}

func (ur *UrlRepository) Create(orig string) (model.Url, error) {
	shortUrl, err := ur.Shorten(orig)
	if err != nil {
		return model.Url{}, err
	}

	u := model.Url{
		Short:     shortUrl,
		Original:  orig,
		CreatedAt: time.Now(),
		Clicks:    0,
	}

	if _, err := ur.FindLink(u.Original); err == nil {
		return model.Url{}, fmt.Errorf("url %s already exists", u.Original)
	}

	data, err := json.Marshal(u)
	if err != nil {
		return model.Url{}, fmt.Errorf("failed to marshal user: %w", err)
	}

	err = ur.repo.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(UrlBucket)
		return b.Put([]byte(u.Original), data)
	})
	if err != nil {
		return model.Url{}, fmt.Errorf("failed to save user: %w", err)
	}

	return u, nil
}

func (ur *UrlRepository) FindByShort(short string) (model.Url, error) {
	var originalURL string

	err := ur.repo.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(UrlBucket)
		data := b.Get([]byte(short))
		if data == nil {
			return fmt.Errorf("short URL not found")
		}
		originalURL = string(data)
		return nil
	})

	if err != nil {
		return model.Url{}, err
	}

	return ur.FindLink(originalURL)
}

func (ur *UrlRepository) FindLink(url string) (model.Url, error) {
	var u model.Url
	err := ur.repo.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(UrlBucket)
		data := b.Get([]byte(url))
		if data == nil {
			return fmt.Errorf("url not found")
		}
		return json.Unmarshal(data, &u)
	})
	if err != nil {
		return model.Url{}, err
	}
	return u, nil
}

func (ur *UrlRepository) DeleteLink(url string) error {
	u, err := ur.FindLink(url)
	if err != nil {
		return err
	}

	return ur.repo.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(UrlBucket)

		// Удаляем запись по оригинальному URL
		if err := b.Delete([]byte(url)); err != nil {
			return err
		}

		// Удаляем запись по короткому URL
		return b.Delete([]byte(u.Short))
	})
}

func (ur *UrlRepository) IncrementClicks(short string) error {
	u, err := ur.FindByShort(short)
	if err != nil {
		return err
	}

	u.Clicks++

	data, err := json.Marshal(u)
	if err != nil {
		return err
	}

	return ur.repo.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(UrlBucket)
		return b.Put([]byte(u.Original), data)
	})
}

func (ur *UrlRepository) ClearAll() error {
	return ur.repo.ClearBucket(UrlBucket)
}

func (ur *UrlRepository) Shorten(originalURL string) (string, error) {
	// Валидация URL
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return "", fmt.Errorf("неверный URL: %v", err)
	}

	shortCode := generateShortCode(originalURL)

	if shortCode == "" {
		return "", fmt.Errorf("не удалось сгенерировать короткий код")
	}

	return shortCode, nil
}

func (ur *UrlRepository) StartAutoCleanup(interval time.Duration) *time.Ticker {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			count, err := ur.cleanExpiredOptimized()
			if err != nil {
				log.Printf("Ошибка очистки: %v", err)
			} else if count > 0 {
				log.Printf("Автоочистка: удалено %d просроченных ссылок", count)
			}
		}
	}()

	return ticker
}

func (ur *UrlRepository) cleanExpiredOptimized() (int, error) {
	var deletedCount int
	cutoffTime := time.Now().Add(-24 * time.Hour)

	// Множество для отслеживания уже обработанных URL
	processed := make(map[string]bool)

	err := ur.repo.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(UrlBucket)

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			key := string(k)

			// Пропускаем уже обработанные
			if processed[key] {
				continue
			}

			var u model.Url
			if err := json.Unmarshal(v, &u); err != nil {
				continue
			}

			// Проверяем срок
			if u.CreatedAt.Before(cutoffTime) {
				processed[u.Short] = true
				processed[u.Original] = true

				// Удаляем обе записи
				if err := b.Delete([]byte(u.Short)); err != nil {
					return err
				}
				if err := b.Delete([]byte(u.Original)); err != nil {
					return err
				}

				deletedCount++
			}
		}
		return nil
	})

	return deletedCount, err
}

func generateShortCode(originalURL string) string {
	hash := sha256.Sum256([]byte(originalURL + time.Now().String()))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return encoded[:8]
}
