package bootstrap

/*
type ASNRegistry struct {
	// List of ASNs & their RDAP base URLs.
	//
	// Stored in a sorted order for fast search.
	Version     string     `json:"version"`
	Publication string     `json:"publication"`
	Description string     `json:"description"`
	Services    []asnRange `json:"services"`
}

type asnRange struct {
	MinASN uint32     // First AS number.
	MaxASN uint32     // Last AS number.
	URLs   []*url.URL // RDAP base URLs.
}

func (r *ASNRegistry) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" || string(data) == `""` {
		return nil
	}

	temp := struct {
		Services [][][]string `json:"services"`
	}{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal ASNRegistry: %w", err)
	}

	// Iterate over services to populate asnRange
	for _, service := range temp.Services {
		if len(service) < 2 {
			return fmt.Errorf("service data is malformed")
		}

		asnRanges, urls := service[0], service[1]
		for _, asnRangeStr := range asnRanges {
			minASN, maxASN, err := parseASNRange(asnRangeStr)
			if err != nil {
				return fmt.Errorf("invalid ASN range %s: %w", asnRangeStr, err)
			}

			parsedURLs, err := parseURLs(urls)
			if err != nil {
				return fmt.Errorf("failed to parse URLs: %w", err)
			}

			r.Services = append(r.Services, asnRange{
				MinASN: minASN,
				MaxASN: maxASN,
				URLs:   parsedURLs,
			})
		}
	}

	return nil
}

func (r *ASNRegistry) MarshalJSON() ([]byte, error) {
	if r == nil {
		return []byte("null"), nil
	}

	temp := struct {
		Services [][][]string `json:"services"`
	}{}

	// Iterate over Services to convert asnRange back into [][]string format
	for _, service := range r.Services {
		var asnRangeStrs []string
		if service.MinASN == service.MaxASN {
			asnRangeStrs = append(asnRangeStrs, fmt.Sprintf("%d", service.MinASN))
		} else {
			asnRangeStrs = append(asnRangeStrs, fmt.Sprintf("%d-%d", service.MinASN, service.MaxASN))
		}
		urlStrs := urlsToStrings(service.URLs)

		temp.Services = append(temp.Services, [][]string{asnRangeStrs, urlStrs})
	}

	// Marshal the transformed data back into JSON
	return json.Marshal(temp)
}



// parseASNRange parses a string range like "36864-37887" into min and max ASN.
func parseASNRange(asnRangeStr string) (uint32, uint32, error) {
	asns := strings.Split(asnRangeStr, "-")
	minASN, err := strconv.ParseUint(asns[0], 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid min ASN: %w", err)
	}

	maxASN := minASN
	if len(asns) == 2 {
		maxASN, err = strconv.ParseUint(asns[1], 10, 32)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid max ASN: %w", err)
		}
	}

	if minASN > maxASN {
		minASN, maxASN = maxASN, minASN
	}

	return uint32(minASN), uint32(maxASN), nil
}

// String returns "ASxxxx" for a single AS, or "ASxxxx-ASyyyy" for a range.
func (a asnRange) String() string {
	if a.MinASN == a.MaxASN {
		return fmt.Sprintf("AS%d", a.MinASN)
	}

	return fmt.Sprintf("AS%d-AS%d", a.MinASN, a.MaxASN)
}
*/
