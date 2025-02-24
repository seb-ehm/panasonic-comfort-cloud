package comfortcloud

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type Authentication struct {
	username   string
	password   string
	token      *Token
	raw        bool
	appVersion string
}

func NewAuthentication(username, password string, token *Token) *Authentication {
	return &Authentication{
		username:   username,
		password:   password,
		token:      token,
		appVersion: XAppVersion,
	}
}

func (a *Authentication) GetNewToken() error {
	slog.Info("Starting token retrieval")
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Prevent automatic redirects
		},
		Jar: jar, //store cookies
	}

	state, codeVerifier, codeChallenge := generateOAuthParameters()
	slog.Debug("OAuth parameters generated", "state", state, "codeChallenge", codeChallenge)

	// Step 1: Authorize

	resp, err := makeAuthorizeRequest(codeChallenge, state, client)
	if err != nil {
		slog.Error("Authorization request failed", "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		slog.Error("Unexpected authorize response", "expected", 302, "got", resp.StatusCode)
		return fmt.Errorf("authorize: expected status 302, got %d", resp.StatusCode)
	}

	location, state, err := handleRedirect(resp, state)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(location, RedirectUri) {
		req, _ := http.NewRequest("GET", BasePathAuth+location, nil)
		req.Header.Set("User-Agent", "okhttp/4.10.0")
		resp, err = client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("authorize_redirect: expected status 200, got %d", resp.StatusCode)
		}

		// Extract `_csrf` from cookies
		var csrf string
		for _, c := range resp.Cookies() {
			if c.Name == "_csrf" {
				csrf = c.Value
				break
			}
		}

		// Step 3: Login
		resp, err = a.submitLoginForm(csrf, state, client)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("login: expected status 200, got %d", resp.StatusCode)
		}

		resp, err = performLoginCallback(resp, client)
		if err != nil {
			return err
		}
	}

	// Step 5: Get Token
	location = resp.Header.Get("Location")
	tokenResponse, err := getToken(location, codeVerifier, client)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	accessToken := tokenResponse.AccessToken

	now := time.Unix(tokenResponse.AccessTokenIssuedAt, 0).UTC()

	apiKey := tokenResponse.getAPIKey(now)

	postUrl := BasePathAcc + "/auth/v2/login"
	timestamp := now.Format("2006-01-02 15:04:05")
	reqBody := `{"language": 0}`
	req, _ := http.NewRequest("POST", postUrl, strings.NewReader(reqBody))
	//req, _ = http.NewRequest("POST", "http://localhost:8080", strings.NewReader(reqBody))
	req.Header.Set("User-Agent", "G-RAC")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("x-app-name", "Comfort Cloud")
	req.Header.Set("x-app-timestamp", timestamp)
	req.Header.Set("x-app-type", "1")
	req.Header.Set("x-app-version", XAppVersion)
	req.Header.Set("x-cfc-api-key", apiKey)
	req.Header.Set("x-user-authorization-v2", "Bearer "+accessToken)

	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("get_acc_client_id: expected status 200, got %d", resp.StatusCode)
	}

	// Extract ACC Client ID
	var accClientResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&accClientResponse)
	if err != nil {
		return err
	}
	accClientID := accClientResponse["clientId"].(string)

	token := tokenResponse
	token.AccClientID = accClientID
	a.token = &token

	return nil
}

func getToken(location string, codeVerifier string, client *http.Client) (Token, error) {
	parsedURL, err := url.Parse(location)
	if err != nil {
		return Token{}, fmt.Errorf("failed to parse redirect URL: %w", err)
	}
	queryParams := parsedURL.Query()
	code := queryParams.Get("code")

	tokenRequest := map[string]string{
		"scope":         "openid",
		"client_id":     AppClientId,
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  RedirectUri,
		"code_verifier": codeVerifier,
	}

	jsonData, _ := json.Marshal(tokenRequest)
	req, _ := http.NewRequest("POST", BasePathAuth+"/oauth/token", strings.NewReader(string(jsonData)))
	req.Header.Set("Auth0-Client", Auth0Client)
	req.Header.Set("User-Agent", "okhttp/4.10.0")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return Token{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Token{}, fmt.Errorf("get_token: expected status 200, got %d", resp.StatusCode)
	}

	// Parse token response
	var tokenResponse Token
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return Token{}, fmt.Errorf("failed to decode token response: %w", err)
	}
	fmt.Println("Token Response:", tokenResponse)

	// Step 6: Get ACC Client ID
	if !tokenResponse.isValid() {
		return Token{}, errors.New("invalid token response")
	}
	return tokenResponse, nil
}

func (a *Authentication) RefreshToken() error {
	//def _refresh_token(self):
	//# do before, so that timestamp is older rather than newer
	//now = datetime.datetime.now()
	//unix_time_token_received = time.mktime(now.timetuple())

	// Prepare the request payload
	payload := map[string]interface{}{
		"scope":         a.token.Scope,
		"client_id":     AppClientId,
		"refresh_token": a.token.RefreshToken,
		"grant_type":    "refresh_token",
	}

	// Prepare the request URL
	tokenUrl := fmt.Sprintf("%s/oauth/token", BasePathAuth)
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data: %w", err)
	}

	// Create the HTTP POST request
	req, err := http.NewRequest(http.MethodPost, tokenUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("failed to create POST request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		err := a.GetNewToken()
		if err != nil {
			return fmt.Errorf("failed to get new token: %w", err)
		}
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	// Parse the response
	var tokenResponse map[string]interface{}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return fmt.Errorf("failed to parse token response: %v", err)
	}
	iat, exp, err := extractIATAndEXPFromJWT(a.token.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to extract IAT: %w", err)
	}
	// Update the token
	a.token = &Token{
		AccessToken:          tokenResponse["access_token"].(string),
		RefreshToken:         tokenResponse["refresh_token"].(string),
		IDToken:              tokenResponse["id_token"].(string),
		AccessTokenIssuedAt:  iat,
		AccessTokenExpiresAt: exp,
		ExpiresInSec:         int(tokenResponse["expires_in"].(float64)),
		AccClientID:          a.token.AccClientID,
		Scope:                tokenResponse["scope"].(string),
	}

	return nil
}

func performLoginCallback(resp *http.Response, client *http.Client) (*http.Response, error) {
	// Step 4: Extract login callback parameters
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)
	hiddenInputs, err := ExtractHiddenInputValues(bodyStr)
	if err != nil {
		return nil, fmt.Errorf("failed to extract hidden input values: %v", err)
	}

	formData := url.Values{}
	for key, value := range hiddenInputs {
		formData.Set(key, value)
	}

	// Step 4.5: Perform the login callback request
	userAgent := "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/113.0.0.0 Mobile Safari/537.36"

	req, err := http.NewRequest("POST", BasePathAuth+"/login/callback", strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating POST request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)

	resp, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound { // 302 Found
		return nil, fmt.Errorf("Unexpected status code: %d\n", resp.StatusCode)
	}
	// Follow redirect
	location := resp.Header.Get("Location")

	req, _ = http.NewRequest("GET", BasePathAuth+location, nil)
	req.Header.Set("User-Agent", "okhttp/4.10.0")
	resp, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making POST request: %w", err)
	}
	if resp.StatusCode != http.StatusFound {
		return nil, fmt.Errorf("Unexpected status code: %d\n", resp.StatusCode)
	}
	return resp, nil
}

func (a *Authentication) submitLoginForm(csrf string, state string, client *http.Client) (*http.Response, error) {
	loginData := map[string]string{
		"client_id":     AppClientId,
		"redirect_uri":  RedirectUri,
		"tenant":        "pdpauthglb-a1",
		"response_type": "code",
		"scope":         OAuthScopes,
		"audience":      OAuthAudience,
		"_csrf":         csrf,
		"state":         state,
		"_intstate":     "deprecated",
		"username":      a.username,
		"password":      a.password,
		"lang":          "en",
		"connection":    "PanasonicID-Authentication",
	}

	jsonData, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", BasePathAuth+"/usernamepassword/login", bytes.NewReader(jsonData))
	req.Header.Set("Auth0-Client", Auth0Client)
	req.Header.Set("User-Agent", "okhttp/4.10.0")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "_csrf="+csrf)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error submitting login form: %w", err)
	}
	return resp, nil
}

func handleRedirect(resp *http.Response, state string) (string, string, error) {
	location := resp.Header.Get("Location")
	parsedURL, err := url.Parse(location)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse redirect URL: %w", err)
	}
	queryParams := parsedURL.Query()
	newState := queryParams.Get("state")
	if newState != "" {
		state = newState
	}
	return location, state, nil
}

func makeAuthorizeRequest(codeChallenge string, state string, client *http.Client) (*http.Response, error) {
	params := url.Values{
		"scope":                 {OAuthScopes},
		"audience":              {OAuthAudience},
		"protocol":              {"oauth2"},
		"response_type":         {"code"},
		"code_challenge":        {codeChallenge},
		"code_challenge_method": {"S256"},
		"auth0Client":           {Auth0Client},
		"client_id":             {AppClientId},
		"redirect_uri":          {RedirectUri},
		"state":                 {state},
	}

	req, err := http.NewRequest("GET", BasePathAuth+"/authorize?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("error building authorize request %w", err)
	}
	req.Header.Set("User-Agent", "okhttp/4.10.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making authorization request %w", err)
	}
	return resp, nil
}

func generateOAuthParameters() (string, string, string) {
	state := generateRandomString(20)
	codeVerifier := generateRandomString(43)

	// Generate Code Challenge
	hasher := sha256.New()
	hasher.Write([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))
	return state, codeVerifier, codeChallenge
}
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

func ExtractHiddenInputValues(htmlBody string) (map[string]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	parameters := make(map[string]string)

	// Find all <input> elements with type="hidden"
	doc.Find("input[type='hidden']").Each(func(_ int, s *goquery.Selection) {
		name, existsName := s.Attr("name")
		value, existsValue := s.Attr("value")

		if existsName && existsValue {
			parameters[name] = value
		}
	})

	return parameters, nil
}

func (a *Authentication) ExecuteGet(url, functionDescription string, expectedStatusCode int) ([]byte, error) {

	if !a.token.isValid() {
		err := a.GetNewToken()
		if err != nil {
			return nil, fmt.Errorf("invalid or expired token. error getting new token: %w", err)
		}
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %v", err)
	}

	// Add headers for the API call
	headers := a.getHeaderForAPICalls()
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != expectedStatusCode {
		return nil, fmt.Errorf(
			"%s: expected status code %d, got %d: %s",
			functionDescription,
			expectedStatusCode,
			resp.StatusCode,
			resp.Status,
		)
	}

	// Read and return the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return body, nil
}

func (a *Authentication) ExecutePost(url string, jsonData map[string]interface{}, functionDescription string, expectedStatusCode int) ([]byte, error) {
	// Ensure the token is valid
	if !a.token.isValid() {
		err := a.GetNewToken()
		if err != nil {
			return nil, fmt.Errorf("invalid or expired token. error getting new token: %w", err)
		}
	}

	// Convert JSON data to bytes
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON data: %v", err)
	}

	// Create the HTTP POST request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}

	// Add headers for the API call
	headers := a.getHeaderForAPICalls()
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != expectedStatusCode {
		return nil, fmt.Errorf(
			"%s: expected status code %d, got %d: %s",
			functionDescription,
			expectedStatusCode,
			resp.StatusCode,
			resp.Status,
		)
	}

	// Read and return the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return body, nil
}

func (a *Authentication) getHeaderForAPICalls() map[string]string {
	now := time.Now()

	headers := map[string]string{
		"Content-Type":            "application/json;charset=utf-8",
		"x-app-name":              "Comfort Cloud",
		"user-agent":              "G-RAC",
		"x-app-timestamp":         now.Format("2006-01-02 15:04:05"),
		"x-app-type":              "1",
		"x-app-version":           a.appVersion,
		"x-cfc-api-key":           a.token.getAPIKey(now),
		"x-client-id":             a.token.AccClientID,
		"x-user-authorization-v2": "Bearer " + a.token.AccessToken,
		"Accept-Encoding":         "gzip, deflate",
		"Accept":                  "*/*",
		"Connection":              "keep-alive",
	}

	return headers
}

func (a *Authentication) Login() error {
	if !a.token.isValid() {
		expired, err := a.token.isAccessTokenExpired()
		if err != nil {
			err := a.GetNewToken()
			if err != nil {
				return fmt.Errorf("invalid or expired token. error getting new token: %w", err)
			}
		}
		if expired {
			err := a.RefreshToken()
			if err != nil {
				err := a.GetNewToken()
				if err != nil {
					return fmt.Errorf("invalid or expired token. error getting new token: %w", err)
				}
			}
		}
	}
	return nil
}

// Logout logs out of the API.
func (a *Authentication) Logout() error {
	// Prepare the URL for the logout request
	logoutUrl := fmt.Sprintf("%s/auth/v2/logout", BasePathAcc)

	// Send the POST request
	response, err := a.ExecutePost(logoutUrl, nil, "logout", http.StatusOK)
	if err != nil {
		return fmt.Errorf("logout request failed: %v", err)
	}

	// Parse the response
	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return fmt.Errorf("failed to parse logout response: %v", err)
	}

	// Check if the logout was successful
	if result["result"].(float64) != 0 {
		// Logout failed, but we don't raise an error (as per the Python implementation)
		fmt.Println("Logout issue detected, but ignoring it")
	}

	return nil
}
