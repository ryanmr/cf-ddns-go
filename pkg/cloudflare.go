package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

func UpdateCloudflare(ip string) error {
	log.Info().Msg("Updating cloudflare dns settings")
	return nil
}

// CloudflareZoneRecord represents a DNS record in Cloudflare
type CloudflareZoneRecord struct {
	ID        string `json:"id"`
	ZoneID    string `json:"zone_id"`
	ZoneName  string `json:"zone_name"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	Proxiable bool   `json:"proxiable"`
	Proxied   bool   `json:"proxied"`
	TTL       int    `json:"ttl"`
	Locked    bool   `json:"locked"`
	Meta      struct {
		AutoAdded           bool   `json:"auto_added"`
		ManagedByApps       bool   `json:"managed_by_apps"`
		ManagedByArgoTunnel bool   `json:"managed_by_argo_tunnel"`
		Source              string `json:"source"`
	} `json:"meta"`
	Comment    string   `json:"comment"`
	Tags       []string `json:"tags"`
	CreatedOn  string   `json:"created_on"`
	ModifiedOn string   `json:"modified_on"`
}

// CloudflareListDnsRecordsResponse represents the response from listing DNS records in Cloudflare
type CloudflareListDnsRecordsResponse struct {
	Result     []CloudflareZoneRecord `json:"result"`
	Success    bool                   `json:"success"`
	Errors     []string               `json:"errors"`
	Messages   []string               `json:"messages"`
	ResultInfo struct {
		Page       int `json:"page"`
		PerPage    int `json:"per_page"`
		Count      int `json:"count"`
		TotalCount int `json:"total_count"`
		TotalPages int `json:"total_pages"`
	} `json:"result_info"`
}

// CloudflarePatchDnsRecordRequest represents the request body for patching a DNS record in Cloudflare
type CloudflarePatchDnsRecordRequest struct {
	Content string `json:"content"`
	Name    string `json:"name"`
	Type    string `json:"type"`
}

const cloudflareApi = "https://api.cloudflare.com/client/v4"

func getAccessToken() string {
	value := os.Getenv("CLOUDFLARE_ACCESS_TOKEN")
	if value == "" {
		panic("cloudflare access token was not set")
	}
	return value
}

func getListDnsRecordsEndpoint(zone string) string {
	return cloudflareApi + "/zones/" + zone + "/dns_records"
}

func getListDnsRecords(zone string) (*CloudflareListDnsRecordsResponse, error) {
	url := getListDnsRecordsEndpoint(zone)
	token := getAccessToken()
	fmt.Printf("🔑 token: %s...\n", token[:8])

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data CloudflareListDnsRecordsResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func selectZoneRecordByName(records []CloudflareZoneRecord, name string) (*CloudflareZoneRecord, error) {
	for _, record := range records {
		if record.Name == name {
			return &record, nil
		}
	}
	return nil, fmt.Errorf("zone record was not found")
}

func getPatchDnsRecordEndpoint(zoneId string, recordId string) string {
	return cloudflareApi + "/zones/" + zoneId + "/dns_records/" + recordId
}

func makePatchBody(existingRecord CloudflareZoneRecord, newIp string) CloudflarePatchDnsRecordRequest {
	return CloudflarePatchDnsRecordRequest{
		Content: newIp,
		Name:    existingRecord.Name,
		Type:    existingRecord.Type,
	}
}

func patchDnsRecord(zoneId string, recordId string, body CloudflarePatchDnsRecordRequest) (*CloudflareZoneRecord, error) {
	url := getPatchDnsRecordEndpoint(zoneId, recordId)
	token := getAccessToken()

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data CloudflareZoneRecord
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func example() {
	// Example usage:
	zone := "example_zone"
	name := "example_name"
	newIp := "123.123.123.123"

	dnsRecordsResp, err := getListDnsRecords(zone)
	if err != nil {
		fmt.Println("Error getting DNS records:", err)
		return
	}

	record, err := selectZoneRecordByName(dnsRecordsResp.Result, name)
	if err != nil {
		fmt.Println("Error selecting DNS record by name:", err)
		return
	}

	patchBody := makePatchBody(*record, newIp)
	updatedRecord, err := patchDnsRecord(zone, record.ID, patchBody)
	if err != nil {
		fmt.Println("Error patching DNS record:", err)
		return
	}

	fmt.Println("Updated DNS record:", updatedRecord)
}
