// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package bootstrap

/*
type NetRegistry struct {
	Version     string                    `json:"version"`
	Publication string                    `json:"publication"`
	Description string                    `json:"description"`
	Services    map[*net.IPNet][]*url.URL `json:"services"`
}

func (r *NetRegistry) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}

	temp := &struct {
		Services [][][]string `json:"services"`
	}{}

	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}

	r.Services = make(map[*net.IPNet][]*url.URL)

	for _, service := range temp.Services {
		for _, cidrStr := range service[0] {
			_, ipNet, err := net.ParseCIDR(cidrStr)
			if err != nil {
				return fmt.Errorf("invalid CIDR block %s: %w", cidrStr, err)
			}

			parsedURL, err := parseURLs(service[1])
			if err != nil {
				return fmt.Errorf("failed to parse URLs: %w", err)
			}

			r.Services[ipNet] = parsedURL
		}
	}
	return nil
}

func (r *NetRegistry) MarshalJSON() ([]byte, error) {
	// If the Services map is nil or empty, return null or an empty string as the JSON output
	if len(r.Services) == 0 {
		return json.Marshal("")
	}

	// Create a temporary structure to hold the formatted services
	temp := struct {
		Services [][][]string `json:"services"`
	}{
		Services: make([][][]string, 0, len(r.Services)),
	}

	// Iterate over the Services map
	for ipNet, urls := range r.Services {
		// Convert *net.IPNet to its CIDR string representation
		cidrStr := ipNet.String()

		// Convert []*url.URL to a list of URL strings
		urlStrs := urlsToStrings(urls)

		// Append the CIDR and URLs to the temp.Services slice in the expected format
		temp.Services = append(temp.Services, [][]string{
			{cidrStr}, // First slice contains CIDR block(s)
			urlStrs,   // Second slice contains URLs
		})
	}

	// Marshal the entire struct into JSON
	return json.Marshal(temp)
}
*/
