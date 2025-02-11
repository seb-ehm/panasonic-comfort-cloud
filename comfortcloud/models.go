package comfortcloud

type Token struct {
	AccessToken                string `json:"access_token"`
	RefreshToken               string `json:"refresh_token"`
	IDToken                    string `json:"id_token"`
	UnixTimestampTokenReceived int64  `json:"unix_timestamp_token_received"`
	ExpiresInSec               int    `json:"expires_in"`
	AccClientID                string `json:"acc_client_id"`
	Scope                      string `json:"scope"`
}
