package types

import (
	"time"

	"github.com/google/uuid"
)

// Release represents a plugin release
type Release struct {
	ID        uuid.UUID `db:"id,omitempty"`
	PluginID  uuid.UUID `db:"plugin_id" json:"plugin_id"`
	Platform  string    `db:"platform" json:"platform"`
	GoVersion string    `db:"go_version" json:"go_version"`
	URL       string    `db:"url" json:"url"`
	Checksum  string    `db:"checksum" json:"checksum"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// RegistryPlugin represents plugin metadata from the registry
type RegistryPlugin struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Version     string     `json:"version"`
	Description string     `json:"description"`
	Author      string     `json:"author"`
	Repository  string     `json:"repository"`
	Official    bool       `json:"official"`
	Releases    []*Release `json:"releases"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
