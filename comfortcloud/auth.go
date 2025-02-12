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

func NewAuthentication(username, password string, token *Token, raw bool) *Authentication {
	return &Authentication{
		username:   username,
		password:   password,
		token:      token,
		raw:        raw,
		appVersion: XAppVersion,
	}
}

func (a *Authentication) GetNewToken() error {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Prevent automatic redirects
		},
		Jar: jar, //store cookies
	}

	state, codeVerifier, codeChallenge := generateOAuthParameters()

	// Step 1: Authorize

	resp, err := makeAuthorizeRequest(codeChallenge, state, client)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
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
		resp, err := a.submitLoginForm(csrf, state, client)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("login: expected status 200, got %d", resp.StatusCode)
		}

		err = performLoginCallback(resp, client)
		if err != nil {
			return err
		}
	}
	// Step 5: Get Token
	location = resp.Header.Get("Location")
	//get code
	parsedURL, err := url.Parse(location)
	if err != nil {
		return fmt.Errorf("failed to parse redirect URL: %w", err)
	}
	queryParams := parsedURL.Query()
	code := queryParams.Get("code")
	//unixTimeTokenReceived := time.Now().Unix()

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

	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("get_token: expected status 200, got %d", resp.StatusCode)
	}

	// Parse token response
	var tokenResponse Token
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}
	fmt.Println("Token Response:", tokenResponse)
	tokenResponse.setIATAndEXP()

	// Step 6: Get ACC Client ID
	if !tokenResponse.isValid() {
		return errors.New("invalid token response")
	}

	accessToken := tokenResponse.AccessToken

	now := time.Unix(tokenResponse.AccessTokenIssuedAt, 0).UTC()

	apiKey := tokenResponse.getAPIKey(now)

	postUrl := BasePathAcc + "/auth/v2/login"
	timestamp := now.Format("2006-01-02 15:04:05")
	reqBody := `{"language": 0}`
	req, _ = http.NewRequest("POST", postUrl, strings.NewReader(reqBody))
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

	fmt.Println("Access Token:", tokenResponse.AccessToken)
	fmt.Println("ACC Client ID:", accClientID)
	token := tokenResponse
	token.AccClientID = accClientID
	a.token = &token

	return nil
}

func performLoginCallback(resp *http.Response, client *http.Client) error {
	// Step 4: Extract login callback parameters
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)
	hiddenInputs, err := ExtractHiddenInputValues(bodyStr)
	if err != nil {
		return fmt.Errorf("failed to extract hidden input values: %v", err)
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
		return fmt.Errorf("error creating POST request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)

	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("error making POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound { // 302 Found
		return fmt.Errorf("Unexpected status code: %d\n", resp.StatusCode)
	}
	// Follow redirect
	location := resp.Header.Get("Location")

	req, _ = http.NewRequest("GET", BasePathAuth+location, nil)
	req.Header.Set("User-Agent", "okhttp/4.10.0")
	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("error making POST request: %w", err)
	}
	if resp.StatusCode != http.StatusFound {
		return fmt.Errorf("Unexpected status code: %d\n", resp.StatusCode)
	}
	return nil
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
