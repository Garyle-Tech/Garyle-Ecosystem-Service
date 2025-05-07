package ota

import "time"

type OTA struct {
	ID           int       `json:"id" db:"id"`
	AppID        string    `json:"app_id" db:"app_id"`
	VersionName  string    `json:"version_name" db:"version_name"`
	VersionCode  int       `json:"version_code" db:"version_code"`
	URL          string    `json:"url" db:"url"`
	ReleaseNotes string    `json:"release_notes" db:"release_notes"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
