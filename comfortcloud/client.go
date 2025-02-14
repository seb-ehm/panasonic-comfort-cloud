package comfortcloud

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

type Client struct {
	auth          *Authentication
	groups        []Group
	devices       []Device
	tokenFileName string
}

func NewClient(username string, password string, tokenFileName string) *Client {
	auth := NewAuthentication(username, password, nil)

	return &Client{
		auth:          auth,
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

func (c *Client) FetchGroupsAndDevices() error {
	// Ensure the client is logged in
	if err := c.ensureLoggedIn(); err != nil {
		return err
	}

	// Fetch and parse groups
	groupURL := c.getGroupURL()
	response, err := c.auth.ExecuteGet(groupURL, "get_groups", http.StatusOK)
	if err != nil {
		return fmt.Errorf("failed to fetch groups: %w", err)
	}

	var result Response
	if err := json.Unmarshal(response, &result); err != nil {
		return fmt.Errorf("failed to parse groups response: %w", err)
	}

	// Reset devices
	c.groups = result.GroupList
	c.devices = nil

	// Populate devices
	for _, group := range c.groups {
		for _, device := range group.DeviceList {
			// Append to devices slice
			c.devices = append(c.devices, device)
		}
		fmt.Println(c.devices)
	}

	return nil
}

func (c *Client) GetDevice(deviceID string) (*Device, error) {
	// Ensure the client is logged in
	if err := c.ensureLoggedIn(); err != nil {
		return nil, err
	}

	// Find the device by DeviceHashGuid or DeviceGuid
	var device *Device
	for _, d := range c.devices {
		if d.DeviceHashGuid == deviceID || hashMD5(d.DeviceGuid) == deviceID {
			device = &d
			break
		}
	}

	if device == nil {
		return nil, fmt.Errorf("device not found: %s", deviceID)
	}

	// Fetch device status using DeviceGuid
	deviceURL := c.getDeviceStatusURL(device.DeviceGuid)
	response, err := c.auth.ExecuteGet(deviceURL, "get_device", http.StatusOK)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch device status: %w", err)
	}
	fmt.Println("Device Status")
	fmt.Println(string(response))
	// Parse response

	if err := json.Unmarshal(response, &device); err != nil {
		return nil, fmt.Errorf("failed to parse device status: %w", err)
	}

	return device, nil
}

func hashMD5(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
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
