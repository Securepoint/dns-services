package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Description struct {
	En string `json:"en"`
	De string `json:"de"`
}

type Service struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Group       string      `json:"group"`
	Description Description `json:"description"`
	Domains     []string    `json:"domains,omitempty"`
	Patterns    []string    `json:"patterns,omitempty"`
}

type CompiledOutput struct {
	Groups   map[string]Description `json:"groups"`
	Services map[string]Service     `json:"services"`
}

func main() {
	const serviceIDsPath = "service_ids.json"

	groups, err := loadGroups("groups.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading groups: %v\n", err)
		os.Exit(1)
	}

	serviceIDs, err := loadServiceIDs(serviceIDsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading service ids: %v\n", err)
		os.Exit(1)
	}

	nextServiceID, err := nextAvailableServiceID(serviceIDs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error preparing service ids: %v\n", err)
		os.Exit(1)
	}

	compiled := CompiledOutput{
		Groups:   groups,
		Services: make(map[string]Service),
	}

	entries, err := os.ReadDir("services")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading services directory: %v\n", err)
		os.Exit(1)
	}

	var cnt int
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		cnt++

		d, err := os.ReadFile(filepath.Join("services", e.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", e.Name(), err)
			os.Exit(1)
		}

		var svc Service
		if err := json.Unmarshal(d, &svc); err != nil {
			fmt.Fprintf(os.Stderr, "error parsing %s: %v\n", e.Name(), err)
			os.Exit(1)
		}

		serviceID := strings.TrimSuffix(e.Name(), ".json")
		numericID, ok := serviceIDs[serviceID]
		if !ok {
			numericID = nextServiceID
			serviceIDs[serviceID] = numericID
			nextServiceID++
			fmt.Printf("assigned new service id %d -> %s\n", numericID, serviceID)
		}

		if _, ok := groups[svc.Group]; !ok {
			fmt.Fprintf(os.Stderr, "unknown group %q in %s\n", svc.Group, e.Name())
			os.Exit(1)
		}

		svc.ID = numericID
		compiled.Services[serviceID] = svc
		fmt.Printf("compiled: %s -> %s (id=%d)\n", serviceID, svc.Group, svc.ID)
	}

	fmt.Printf("found %d service files\n", cnt)

	if err := writeJSONFile(serviceIDsPath, serviceIDs); err != nil {
		fmt.Fprintf(os.Stderr, "error writing %s: %v\n", serviceIDsPath, err)
		os.Exit(1)
	}

	out, err := json.MarshalIndent(compiled, "", "    ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling output: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("services.json", out, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing services.json: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("successfully compiled %d services in %d groups to services.json using %d stable ids\n", cnt, len(compiled.Groups), len(serviceIDs))
}

func loadGroups(path string) (map[string]Description, error) {
	d, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var g map[string]Description
	if err := json.Unmarshal(d, &g); err != nil {
		return nil, err
	}

	return g, nil
}

func loadServiceIDs(path string) (map[string]int, error) {
	d, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]int), nil
		}
		return nil, err
	}

	var ids map[string]int
	if err := json.Unmarshal(d, &ids); err != nil {
		return nil, err
	}

	if ids == nil {
		ids = make(map[string]int)
	}

	return ids, nil
}

func nextAvailableServiceID(ids map[string]int) (int, error) {
	maxID := 0
	seen := make(map[int]string, len(ids))

	for serviceKey, id := range ids {
		if id <= 0 {
			return 0, fmt.Errorf("service %q has invalid id %d", serviceKey, id)
		}
		if existingService, ok := seen[id]; ok {
			return 0, fmt.Errorf("duplicate id %d for services %q and %q", id, existingService, serviceKey)
		}
		seen[id] = serviceKey
		if id > maxID {
			maxID = id
		}
	}

	return maxID + 1, nil
}

func writeJSONFile(path string, value any) error {
	out, err := json.MarshalIndent(value, "", "    ")
	if err != nil {
		return err
	}

	out = append(out, '\n')
	return os.WriteFile(path, out, 0644)
}
