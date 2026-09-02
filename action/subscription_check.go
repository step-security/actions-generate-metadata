// Copyright IBM Corp. 2021, 2025
// Copyright 2026 StepSecurity
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func runSubscriptionCheck() {
	eventPath := os.Getenv("GITHUB_EVENT_PATH")
	var repoPrivate *bool

	if eventPath != "" {
		if eventData, err := os.ReadFile(eventPath); err == nil {
			var event struct {
				Repository struct {
					Private *bool `json:"private"`
				} `json:"repository"`
			}
			if err := json.Unmarshal(eventData, &event); err == nil {
				repoPrivate = event.Repository.Private
			}
		}
	}

	upstream := "hashicorp/actions-generate-metadata"
	action := os.Getenv("GITHUB_ACTION_REPOSITORY")
	docsURL := "https://docs.stepsecurity.io/actions/stepsecurity-maintained-actions"

	fmt.Println()
	fmt.Println("\x1b[1;36mStepSecurity Maintained Action\x1b[0m")
	fmt.Printf("Secure drop-in replacement for %s\n", upstream)
	if repoPrivate != nil && !*repoPrivate {
		fmt.Println("\x1b[32m✓ Free for public repositories\x1b[0m")
	}
	fmt.Printf("\x1b[36mLearn more:\x1b[0m %s\n", docsURL)
	fmt.Println()

	if repoPrivate != nil && !*repoPrivate {
		return
	}

	serverURL := os.Getenv("GITHUB_SERVER_URL")
	if serverURL == "" {
		serverURL = "https://github.com"
	}

	body := map[string]string{"action": action}
	if serverURL != "https://github.com" {
		body["ghes_server"] = serverURL
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Timeout or API not reachable. Continuing to next step.")
		return
	}

	apiURL := fmt.Sprintf("https://agent.api.stepsecurity.io/v1/github/%s/actions/maintained-actions-subscription", os.Getenv("GITHUB_REPOSITORY"))

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Post(apiURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Timeout or API not reachable. Continuing to next step.")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		fmt.Printf("::error::\x1b[1;31mThis action requires a StepSecurity subscription for private repositories.\x1b[0m\n")
		fmt.Printf("::error::\x1b[31mLearn how to enable a subscription: %s\x1b[0m\n", docsURL)
		os.Exit(1)
	}
}
