package types

import (
	"time"

	"github.com/google/uuid"
)

// Release represents a plugin release
type Release struct {
	ID        uuid.UUID `db:"id,omitempty"`
	PluginID  uuid.UUID `db:"plugin_id" json:"plugin_id"`
	Version   string    `db:"version" json:"version"`
	Platform  string    `db:"platform" json:"platform"`
	GoVersion string    `db:"go_version" json:"go_version"`
	URL       string    `db:"url" json:"url"`
	Checksum  string    `db:"checksum" json:"checksum"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// RegistryPlugin represents plugin metadata from the registry
type RegistryPlugin struct {
	ID            uuid.UUID  `db:"id,omitempty" json:"id"`
	Name          string     `db:"name" json:"name"`
	Version       string     `db:"version" json:"version"`
	Description   string     `db:"description" json:"description"`
	Author        string     `db:"author" json:"author"`
	Repository    string     `db:"repository" json:"repository"`
	Official      bool       `db:"official" json:"official"`
	DownloadCount int64      `db:"download_count" json:"download_count"`
	Releases      []*Release `json:"releases"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
}
