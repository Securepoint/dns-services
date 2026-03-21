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
	groups, err := loadGroups("groups.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading groups: %v\n", err)
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

		if _, ok := groups[svc.Group]; !ok {
			fmt.Fprintf(os.Stderr, "unknown group %q in %s\n", svc.Group, e.Name())
			os.Exit(1)
		}

		compiled.Services[serviceID] = svc
		fmt.Printf("compiled: %s -> %s\n", serviceID, svc.Group)
	}

	fmt.Printf("found %d service files\n", cnt)

	out, err := json.MarshalIndent(compiled, "", "    ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling output: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("services.json", out, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing services.json: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("successfully compiled %d services in %d groups to services.json\n", cnt, len(compiled.Groups))
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
