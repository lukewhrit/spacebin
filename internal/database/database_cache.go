package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"time"

	"github.com/lukewhrit/spacebin/internal/config"
	"github.com/patrickmn/go-cache"
)

type EphemeralDb struct {
	cache *cache.Cache
}

func NewEphemeralDb() (*EphemeralDb, error) {
	expire := time.Duration(config.Config.ExpirationAge) * time.Hour
	return &EphemeralDb{cache: cache.New(
		expire, expire*2,
	)}, nil
}

func (e EphemeralDb) Migrate(ctx context.Context) error {
	return nil
}

func (e EphemeralDb) Close() error {
	return nil
}

func (e EphemeralDb) GetDocument(ctx context.Context, id string) (Document, error) {
	jsonIntf, found := e.cache.Get(fmt.Sprintf("document_%s", id))
	if !found {
		return Document{}, sql.ErrNoRows
	}

	jsonBytes, ok := jsonIntf.([]byte)
	if !ok {
		return Document{}, fmt.Errorf("document corrupted")
	}

	var document Document
	err := json.Unmarshal(jsonBytes, &document)
	if err != nil {
		return Document{}, fmt.Errorf("document corrupted: %w", err)
	}

	return document, nil

}

func (e EphemeralDb) CreateDocument(ctx context.Context, id, content string) error {
	t := time.Now()
	document := Document{
		ID:        id,
		Content:   content,
		CreatedAt: t,
		UpdatedAt: t,
	}
	jsonBytes, err := json.Marshal(document)
	if err != nil {
		return err
	}
	e.cache.Set(fmt.Sprintf("document_%s", id), jsonBytes, 0)
	return nil
}

func (e EphemeralDb) ListDocuments(ctx context.Context) ([]Document, error) {
	items := e.cache.Items()
	documents := make([]Document, 0)
	for k, v := range items {
		if v.Expired() {
			continue
		}
		jsonIntf := v.Object
		jsonBytes, ok := jsonIntf.([]byte)
		if !ok {
			log.Err(fmt.Errorf("document '%s' corrupted", k))
		}

		var document Document
		err := json.Unmarshal(jsonBytes, &document)
		if err != nil {
			log.Err(fmt.Errorf("document '%s' corrupted: %w", k, err))
		}

		documents = append(documents, document)
	}

	return documents, nil
}
