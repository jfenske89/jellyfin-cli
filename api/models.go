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

// Session defines a Jellyfin user session as a generic map with getters for specific fields
type Session map[string]interface{}

func (s Session) GetUserName() string {
	return s.getString("UserName")
}

func (s Session) GetDeviceName() string {
	return s.getString("DeviceName")
}

func (s Session) GetLastActivityDate() time.Time {
	if rawVal, ok := s["LastActivityDate"]; ok && rawVal != nil {
		if val, ok := rawVal.(time.Time); ok {
			return val
		} else if val, ok := rawVal.(string); ok {
			// 2024-02-18T16:31:11.9906841Z
			if result, err := time.Parse("2006-01-02T15:04:05.9999999Z", val); err == nil {
				return result
			}
		}
	}

	return time.Time{}
}

func (s Session) getString(key string) string {
	if rawVal, ok := s[key]; ok && rawVal != nil {
		if val, ok := rawVal.(string); ok {
			return val
		}
	}

	return ""
}
