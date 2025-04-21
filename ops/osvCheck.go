package ops

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Structs for the OSV API request/response
type osvQuery struct {
	Package osvPackage `json:"package"`
	Version string     `json:"version"`
}

type osvPackage struct {
	Name      string `json:"name"`
	Ecosystem string `json:"ecosystem"`
}

type osvBatchRequest struct {
	Queries []osvQuery `json:"queries"`
}

type osvBatchResponse struct {
	Results []struct {
		Vulns []struct {
			ID      string   `json:"id"`
			Details string   `json:"details"`
			Aliases []string `json:"aliases"`
			Refs    []struct {
				Type string `json:"type"`
				URL  string `json:"url"`
			} `json:"references"`
		} `json:"vulns"`
	} `json:"results"`
}

// Provide a map of dependencies with versions
func CheckVulnerabilities(deps map[string]string) {
	var queries []osvQuery
	for name, version := range deps {
		queries = append(queries, osvQuery{
			Package: osvPackage{
				Name:      name,
				Ecosystem: "Go",
			},
			Version: version,
		})
	}

	reqBody, err := json.Marshal(osvBatchRequest{Queries: queries})
	if err != nil {
		fmt.Println("Failed to marshal request:", err)
		return
	}

	resp, err := http.Post("https://api.osv.dev/v1/querybatch", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		fmt.Println("Failed to query OSV:", err)
		return
	}
	defer resp.Body.Close()

	var osvResp osvBatchResponse
	if err := json.NewDecoder(resp.Body).Decode(&osvResp); err != nil {
		fmt.Println("Failed to decode OSV response:", err)
		return
	}

	vulnCount := 0
	for i, result := range osvResp.Results {
		if len(result.Vulns) == 0 {
			continue
		}
		// Match back to the dependency name
		fmt.Printf("\n📦 %s %s\n", queries[i].Package.Name, queries[i].Version)
		for _, vuln := range result.Vulns {
			vulnCount++
			fmt.Printf("  🚨 %s: %s\n", vuln.ID, vuln.Details)
			for _, ref := range vuln.Refs {
				if ref.Type == "WEB" {
					fmt.Printf("    🔗 %s\n", ref.URL)
				}
			}
		}
	}

	if vulnCount == 0 {
		fmt.Println("No known vulnerabilities found.")
	}
}
