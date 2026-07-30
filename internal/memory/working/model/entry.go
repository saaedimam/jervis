package model

import (
	"github.com/ioriimasu/jervis/internal/memory/contracts"
	"time"
)

type entry struct {
	id        string
	content   any
	metadata  map[string]string
	timestamp time.Time
}

func NewEntry(id string, content any, metadata map[string]string) contracts.Entry {
	// Defensive copy of metadata
	meta := make(map[string]string)
	for k, v := range metadata {
		meta[k] = v
	}

	return &entry{
		id:        id,
		content:   content,
		metadata:  meta,
		timestamp: time.Now().UTC(),
	}
}

func (e *entry) ID() string {
	return e.id
}

func (e *entry) Content() any {
	return e.content
}

func (e *entry) Metadata() map[string]string {
	// Defensive copy
	meta := make(map[string]string)
	for k, v := range e.metadata {
		meta[k] = v
	}
	return meta
}

func (e *entry) Timestamp() time.Time {
	return e.timestamp
}
