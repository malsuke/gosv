package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/malsuke/govs/pkg/vuln/domain"
)

func main() {
	vuln, err := domain.GetVulnerabilityByCVEID("CVE-2020-22452", "")
	if err != nil {
		log.Fatalf("failed to get vulnerability: %v", err)
	}
	json, err := json.Marshal(vuln)
	if err != nil {
		log.Fatalf("failed to marshal vulnerability: %v", err)
	}
	fmt.Println(string(json))
}
