package amadeus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// AmadeusProvider is a provider for the Amadeus API
type AmadeusProvider struct {
	clientID     string
	clientSecret string
}

func NewAmadeusProvider(clientID, clientSecret string) *AmadeusProvider {
	return &AmadeusProvider{clientID: clientID, clientSecret: clientSecret}
}

func (p *AmadeusProvider) getAccessToken() (string, error) {
	apiURL := "https://test.api.amadeus.com/v1/security/oauth2/token"
	data := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", p.clientID, p.clientSecret)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(data))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", fmt.Errorf("error unmarshalling response: %v", err)
	}

	return tokenResponse.AccessToken, nil
}

// GetFlightFees fetches flight fees from the Amadeus API
func (p *AmadeusProvider) GetFlightFees(from, to string) (map[string]interface{}, error) {
	token, err := p.getAccessToken()
	if err != nil {
		return nil, err
	}

	apiURL := fmt.Sprintf("https://test.api.amadeus.com/v2/shopping/flight-offers?originLocationCode=%s&destinationLocationCode=%s&departureDate=2025-02-25&adults=1", from, to)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return result, nil
}
