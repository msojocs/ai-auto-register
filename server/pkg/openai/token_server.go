package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SentinelInfo struct {
	Persona     string `json:"persona"`
	Token       string `json:"token"`
	ExpireAfter int    `json:"expire_after"`
	ExpirAt     int64  `json:"expire_at"`
	Turnstile   struct {
		Required bool   `json:"required"`
		Dx       string `json:"dx"`
	} `json:"turnstile"`
	Proofofwork struct {
		Required   bool   `json:"required"`
		Seed       string `json:"seed"`
		Difficulty string `json:"difficulty"`
	} `json:"proofofwork"`
}

const defaultSentinelReqURL = "https://sentinel.openai.com/backend-api/sentinel/req"

type SentinelToken struct {
	baseUrl        string
	flowName       string
	did            string
	proof          string
	sentinelInfo   *SentinelInfo
	sentinelReqURL string
}

// NewSentinelToken creates a SentinelToken. Pass reqURL="" to use the default sentinel endpoint.
func NewSentinelToken(url, flowName, did, reqURL string) *SentinelToken {
	if reqURL == "" {
		reqURL = defaultSentinelReqURL
	}
	return &SentinelToken{
		baseUrl:        url,
		flowName:       flowName,
		did:            did,
		sentinelReqURL: reqURL,
	}
}

func (s *SentinelToken) Req(client *http.Client) (*SentinelInfo, error) {
	// 1. 生成proof
	var proof string
	{
		url := s.baseUrl + "/proof"
		// GET
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get proof: %s", resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		proof = string(body)
		s.proof = proof
	}

	// 2. 请求OpenAI
	{
		reqURL := s.sentinelReqURL
		payload := map[string]string{
			"p":    proof,
			"id":   s.did,
			"flow": "authorize_continue",
		}
		headers := map[string]string{
			"User-Agent": "",
			"Referer":    "https://sentinel.openai.com/backend-api/sentinel/frame.html?sv=20260219f9f6",
		}
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest("POST", reqURL, io.NopCloser(io.Reader(bytes.NewReader(jsonData))))
		if err != nil {
			return nil, err
		}
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get sentinel token: %s", resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var sentinelInfo SentinelInfo
		if err := json.Unmarshal(body, &sentinelInfo); err != nil || sentinelInfo.Token == "" {
			return nil, fmt.Errorf("failed to parse sentinel token response: %w; response: %s", err, string(body))
		}
		s.sentinelInfo = &sentinelInfo
	}

	return s.sentinelInfo, nil
}

func (s *SentinelToken) GetSentinelHeader() (map[string]string, error) {
	var proof, turnstileToken, soToken string
	{
		url := s.baseUrl + "/turnstile"
		if s.proof == "" || s.sentinelInfo == nil {
			return nil, fmt.Errorf("proof or sentinelInfo is empty")
		}
		str, err := json.Marshal(map[string]interface{}{
			"proof":        s.proof,
			"sentinelInfo": s.sentinelInfo,
		})
		if err != nil {
			return nil, err
		}
		resp, err := http.Post(url, "application/json", bytes.NewReader(str))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get proof: %s", resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var data map[string]string
		if err := json.Unmarshal(body, &data); err != nil {
			return nil, fmt.Errorf("failed to parse proof response: %w; response: %s", err, string(body))
		}
		proof = data["enforcementToken"]
		turnstileToken = data["turnstileToken"]
		soToken = data["soToken"]
	}

	headers := map[string]string{
		"p":  proof,
		"t":  turnstileToken,
		"c":  s.sentinelInfo.Token,
		"so": soToken,
	}
	return headers, nil

}
