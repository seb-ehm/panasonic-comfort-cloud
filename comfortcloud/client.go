package comfortcloud

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
)

type Client struct {
	auth          *Authentication
	groups        []map[string]interface{}
	devices       []Device
	deviceIndexer map[string]string // Maps device ID to device GUID
	tokenFileName string
}

func NewClient(username string, password string, tokenFileName string) *Client {
	auth := NewAuthentication(username, password, nil)

	return &Client{
		auth:          auth,
		deviceIndexer: make(map[string]string),
		tokenFileName: tokenFileName,
	}
}

func (c *Client) Login() error {
	if c.auth.token.isValid() {
		return nil
	}
	var token Token
	tokenFile, err := os.ReadFile(c.tokenFileName)
	if err == nil {
		if err := json.Unmarshal(tokenFile, &token); err == nil {
			if token.isValid() {
				c.auth.token = &token
				return nil
			}
		}
	}
	err2 := c.auth.Login()
	if err2 != nil {
		return fmt.Errorf("token file invalid: %w, failed to login to Comfort Cloud. %w", err, err2)
	}
	updatedToken := c.auth.token
	tokenJSON, err := json.MarshalIndent(updatedToken, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	if err := os.WriteFile(c.tokenFileName, tokenJSON, 0644); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}
	return nil
}

func (c *Client) ensureLoggedIn() error {
	err := c.auth.Login()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Logout() error {
	return c.auth.Logout()
}

func (c *Client) GetDevices() error {
	//TODO implement me
	panic("implement me")
}

func (c *Client) GetDevice() error {
	//TODO implement me
	panic("implement me")
}

func (c *Client) GetDeviceState() error {
	//TODO implement me
	panic("implement me")
}

func (c *Client) SetDeviceState() error {
	//TODO implement me
	panic("implement me")
}

// getGroupURL returns the URL for retrieving groups.
func (c *Client) getGroupURL() string {
	//return "http://localhost:8080"
	return fmt.Sprintf("%s/device/group", BasePathAcc)
}

// getDeviceStatusURL returns the URL for retrieving device status.
func (c *Client) getDeviceStatusURL(guid string) string {
	escapedGUID := regexp.MustCompile(`(?i)%2f`).ReplaceAllString(url.QueryEscape(guid), "f")
	return fmt.Sprintf("%s/deviceStatus/%s", BasePathAcc, escapedGUID)
}

// getDeviceStatusControlURL returns the URL for controlling device status.
func (c *Client) getDeviceStatusControlURL() string {
	return fmt.Sprintf("%s/deviceStatus/control", BasePathAcc)
}

// getDeviceHistoryURL returns the URL for retrieving device history.
func (c *Client) getDeviceHistoryURL() string {
	return fmt.Sprintf("%s/deviceHistoryData", BasePathAcc)
}
