package comfortcloud

type Authentication struct {
	username   string
	password   string
	token      *Token
	raw        bool
	appVersion string
}

func NewAuthentication(username, password string, token *Token, raw bool) *Authentication {
	return &Authentication{
		username:   username,
		password:   password,
		token:      token,
		raw:        raw,
		appVersion: XAppVersion,
	}
}
