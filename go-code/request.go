package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func checkDuplicates(data Prequest) (*CheckDuplicateResponse, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://localhost:5000/check-duplicates", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result CheckDuplicateResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, err
}

func updateDuplicate(updates OutPayload) error {
	payload, err := json.Marshal(updates)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:5000/update-duplicates", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}