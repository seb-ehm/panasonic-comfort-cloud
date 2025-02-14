package comfortcloud

import (
	"encoding/json"
	"fmt"
	"os"
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

	tokenFile, err := os.ReadFile(c.tokenFileName)
	if err != nil {
		return fmt.Errorf("failed to read token file: %v", err)
	}

	var token Token
	if err := json.Unmarshal(tokenFile, &token); err != nil {
		return fmt.Errorf("failed to parse token file: %w", err)
	}
	if !token.isValid() {
		return fmt.Errorf("invalid token file")
	}

	c.auth.token = &token

	err = c.auth.Login()
	if err != nil {
		return err
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
