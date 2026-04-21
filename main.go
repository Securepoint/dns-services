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

type CountryDefinition struct {
	Name Description `json:"name"`
}

type Country struct {
	ID   int         `json:"id"`
	Code string      `json:"code"`
	Name Description `json:"name"`
}

type CompiledServicesOutput struct {
	Groups   map[string]Description `json:"groups"`
	Services map[string]Service     `json:"services"`
}

type CompiledCountriesOutput struct {
	Countries map[string]Country `json:"countries"`
}

type StableIDs struct {
	Services  map[string]int `json:"services"`
	Countries map[string]int `json:"countries"`
}

func main() {
	const stableIDsPath = "ids.json"
	const legacyServiceIDsPath = "service_ids.json"
	const legacyCountryIDsPath = "country_ids.json"

	groups, err := loadGroups("groups.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading groups: %v\n", err)
		os.Exit(1)
	}

	stableIDs, err := loadStableIDs(stableIDsPath, legacyServiceIDsPath, legacyCountryIDsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading stable ids: %v\n", err)
		os.Exit(1)
	}

	nextServiceID, err := nextAvailableID(stableIDs.Services)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error preparing service ids: %v\n", err)
		os.Exit(1)
	}

	compiledServices := CompiledServicesOutput{
		Groups:   groups,
		Services: make(map[string]Service),
	}

	serviceCount, err := compileServices("services", groups, stableIDs.Services, nextServiceID, compiledServices.Services)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error compiling services: %v\n", err)
		os.Exit(1)
	}

	if err := writeJSONFile("services.json", compiledServices); err != nil {
		fmt.Fprintf(os.Stderr, "error writing services.json: %v\n", err)
		os.Exit(1)
	}

	nextCountryID, err := nextAvailableID(stableIDs.Countries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error preparing country ids: %v\n", err)
		os.Exit(1)
	}

	compiledCountries := CompiledCountriesOutput{
		Countries: make(map[string]Country),
	}

	countryCount, err := compileCountries("countries", stableIDs.Countries, nextCountryID, compiledCountries.Countries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error compiling countries: %v\n", err)
		os.Exit(1)
	}

	if err := writeJSONFile(stableIDsPath, stableIDs); err != nil {
		fmt.Fprintf(os.Stderr, "error writing %s: %v\n", stableIDsPath, err)
		os.Exit(1)
	}

	if err := writeJSONFile("countries.json", compiledCountries); err != nil {
		fmt.Fprintf(os.Stderr, "error writing countries.json: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("successfully compiled %d services in %d groups to services.json using %d stable ids\n", serviceCount, len(compiledServices.Groups), len(stableIDs.Services))
	fmt.Printf("successfully compiled %d countries to countries.json using %d stable ids\n", countryCount, len(stableIDs.Countries))
}

func compileServices(dir string, groups map[string]Description, serviceIDs map[string]int, nextServiceID int, dst map[string]Service) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	var count int
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		count++

		d, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return 0, fmt.Errorf("reading %s: %w", e.Name(), err)
		}

		var svc Service
		if err := json.Unmarshal(d, &svc); err != nil {
			return 0, fmt.Errorf("parsing %s: %w", e.Name(), err)
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
			return 0, fmt.Errorf("unknown group %q in %s", svc.Group, e.Name())
		}

		svc.ID = numericID
		dst[serviceID] = svc
		fmt.Printf("compiled service: %s -> %s (id=%d)\n", serviceID, svc.Group, svc.ID)
	}

	fmt.Printf("found %d service files\n", count)
	return count, nil
}

func compileCountries(dir string, countryIDs map[string]int, nextCountryID int, dst map[string]Country) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	var count int
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		count++

		d, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return 0, fmt.Errorf("reading %s: %w", e.Name(), err)
		}

		var def CountryDefinition
		if err := json.Unmarshal(d, &def); err != nil {
			return 0, fmt.Errorf("parsing %s: %w", e.Name(), err)
		}

		countryCode := strings.ToUpper(strings.TrimSuffix(e.Name(), ".json"))
		numericID, ok := countryIDs[countryCode]
		if !ok {
			numericID = nextCountryID
			countryIDs[countryCode] = numericID
			nextCountryID++
			fmt.Printf("assigned new country id %d -> %s\n", numericID, countryCode)
		}

		dst[countryCode] = Country{
			ID:   numericID,
			Code: countryCode,
			Name: def.Name,
		}
		fmt.Printf("compiled country: %s (id=%d)\n", countryCode, numericID)
	}

	fmt.Printf("found %d country files\n", count)
	return count, nil
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

func loadStableIDs(path, legacyServicesPath, legacyCountriesPath string) (StableIDs, error) {
	d, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return loadLegacyStableIDs(legacyServicesPath, legacyCountriesPath)
		}
		return StableIDs{}, err
	}

	var ids StableIDs
	if err := json.Unmarshal(d, &ids); err != nil {
		return StableIDs{}, err
	}

	ensureStableIDMaps(&ids)

	return ids, nil
}

func loadLegacyStableIDs(servicesPath, countriesPath string) (StableIDs, error) {
	serviceIDs, err := loadLegacyStableIDMap(servicesPath)
	if err != nil {
		return StableIDs{}, err
	}

	countryIDs, err := loadLegacyStableIDMap(countriesPath)
	if err != nil {
		return StableIDs{}, err
	}

	return StableIDs{
		Services:  serviceIDs,
		Countries: countryIDs,
	}, nil
}

func loadLegacyStableIDMap(path string) (map[string]int, error) {
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

func ensureStableIDMaps(ids *StableIDs) {
	if ids.Services == nil {
		ids.Services = make(map[string]int)
	}
	if ids.Countries == nil {
		ids.Countries = make(map[string]int)
	}
}

func nextAvailableID(ids map[string]int) (int, error) {
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
