package api

import "time"

type JellyfinApiConfig struct {
	// BaseUrl is the base URL to the Jellyfin API (for example: http://127.0.0.1:8096/emby/)
	BaseUrl string

	// Token is the API token for authenticating requests
	Token string

	// SkipSslVerify will control disabling SSL verification for self-signed certificates if using https
	SkipSslVerify bool
}

type Session struct {
	User   string    `json:"UserName"`
	Device string    `json:"DeviceName"`
	Date   time.Time `json:"LastActivityDate"`
}
