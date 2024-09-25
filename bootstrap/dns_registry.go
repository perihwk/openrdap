package bootstrap

/*
type DNSRegistry struct {
	Version     string              `json:"version"`
	Publication string              `json:"publication"`
	Description string              `json:"description"`
	Services    map[string][]string `json:"services"`
}

func (r *DNSRegistry) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}

	temp := &struct {
		Services [][]interface{} `json:"services"`
	}{}

	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}

	r.Services = make(map[string][]string)

	for _, service := range temp.Services {
		for _, domain := range service[0].([]interface{}) {
			for _, url := range service[1].([]interface{}) {
				r.Services[domain.(string)] = append(r.Services[domain.(string)], url.(string))
			}
		}
	}
	return nil
}

func (r *DNSRegistry) MarshalJSON() ([]byte, error) {
	if len(r.Services) == 0 {
		return json.Marshal(nil) // Handles empty Services map as "null"
	}

	temp := &struct {
		Services [][]interface{} `json:"services"`
	}{}

	for domain, urls := range r.Services {
		urlsAsInterface := make([]interface{}, len(urls))
		for i, url := range urls {
			urlsAsInterface[i] = url
		}

		service := []interface{}{
			[]interface{}{domain}, // domain as a list
			urlsAsInterface,       // urls as a list
		}
		temp.Services = append(temp.Services, service)
	}

	return json.Marshal(temp)
}
*/
