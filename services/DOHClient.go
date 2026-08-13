package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

// [struct(JSON)]
type DoHResponse struct {
	Answer []struct {
		Data string `json:"data"`
		Type int    `json:"type"`
	} `json:"answer"`
}

func resolveDoH(ctx context.Context, host string) (string, error) {
	dohURL := fmt.Sprintf("https://1.1.1.1/dns-query?name=%s&type=A", host)
	req, err := http.NewRequestWithContext(ctx, "GET", dohURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/dns-json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result DoHResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	for _, ans := range result.Answer {
		if ans.Type == 1 { // A record
			return ans.Data, nil
		}
	}
	return "", fmt.Errorf("no IPv4 record found for %s via DoH", host)
}

func NewDOHHTTPClient() *http.Client {
	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}

			if net.ParseIP(host) != nil {
				return dialer.DialContext(ctx, network, addr)
			}

			ip, err := resolveDoH(ctx, host)
			if err != nil {
				return dialer.DialContext(ctx, network, addr)
			}

			return dialer.DialContext(ctx, network, net.JoinHostPort(ip, port))
		},
	}

	return &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
	}
}
