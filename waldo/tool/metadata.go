package tool

import (
	"time"
)

type ArtifactMetadata struct {
	AppID       string    `json:"appID,omitempty"`
	BuildPath   string    `json:"buildPath"`
	UploadTime  time.Time `json:"uploadTime,omitempty"`
	UploadToken string    `json:"uploadToken,omitempty"`
}
