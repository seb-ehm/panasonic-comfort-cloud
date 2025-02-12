package comfortcloud

import "fmt"

const (
	AppClientId      = "Xmy6xIYIitMxngjB2rHvlm6HSDNnaMJx"
	Auth0Client      = "eyJuYW1lIjoiQXV0aDAuQW5kcm9pZCIsImVudiI6eyJhbmRyb2lkIjoiMzAifSwidmVyc2lvbiI6IjIuOS4zIn0="
	RedirectUri      = "panasonic-iot-cfc://authglb.digital.panasonic.com/android/com.panasonic.ACCsmart/callback"
	BasePathAuth     = "https://authglb.digital.panasonic.com"
	BasePathAcc      = "https://accsmart.panasonic.com"
	XAppVersion      = "1.22.0"
	OAuthScopes      = "openid offline_access comfortcloud.control a2w.control"
	OAuthAudienceURL = "https://digital.panasonic.com/%s/api/v1/"
)

var OAuthAudience = fmt.Sprintf(OAuthAudienceURL, AppClientId)
