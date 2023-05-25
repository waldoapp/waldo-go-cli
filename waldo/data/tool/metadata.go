package tool

import (
	"time"
)

type ArtifactMetadata struct {
	BuildPath   string    `json:"build_path"`
	UploadTime  time.Time `json:"uploadTime,omitempty"`
	UploadToken string    `json:"uploadToken,omitempty"`
}
