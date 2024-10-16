package openrdap

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testData struct {
	Name           string
	SampleDataPath string
	ExpectedVCard  VCard
	ExpectedError  bool
}

func TestParseJCard(t *testing.T) {
	// load sample data from file
	sampleData := []testData{
		{
			Name:           "TestSuccessNoLabel",
			SampleDataPath: "test/jcard/example.json",
			ExpectedVCard: VCard{
				Version:  "4.0",
				FullName: "Simon Perreault",
				Address: Address{
					Label:           "",
					PostOfficeBox:   "",
					ExtendedAddress: "Suite D2-630",
					StreetAddress:   "2875 Laurier",
					Locality:        "Quebec",
					Region:          "QC",
					PostalCode:      "G1V 2M2",
					Country:         "Canada",
				},
				Kind:      "",
				Email:     "simon.perreault@viagenie.ca",
				Telephone: "tel:+1-418-262-6501",
				Org:       "Viagenie",
			},
			ExpectedError: false,
		},
		{
			Name:           "TestSuccessWithLabel",
			SampleDataPath: "test/jcard/example_label.json",
			ExpectedVCard: VCard{
				Version:  "4.0",
				FullName: "Simon Perreault",
				Address: Address{
					Label:           "123 Maple Ave\nSuite 901\nVancouver\nBC\nA1B 2C9\nCanada",
					PostOfficeBox:   "",
					ExtendedAddress: "Suite 901",
					StreetAddress:   "123 Maple Ave",
					Locality:        "Vancouver",
					Region:          "BC",
					PostalCode:      "A1B 2C9",
					Country:         "Canada",
				},
				Kind:      "",
				Email:     "simon.perreault@viagenie.ca",
				Telephone: "tel:+1-418-262-6501",
				Org:       "Viagenie",
			},
			ExpectedError: false,
		},
	}

	for _, test := range sampleData {
		fileData, err := os.ReadFile(test.SampleDataPath)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}

		// Unmarshal the JSON file into a slice of interfaces
		var jcardData []interface{}
		if err := json.Unmarshal(fileData, &jcardData); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// Test parseJCard function
		vcard, err := parseJCard(jcardData)
		if err != nil {
			t.Errorf("Failed to parse jCard: %v", err)
		}

		if !cmp.Equal(vcard, test.ExpectedVCard) {
			t.Fatalf("test %s failed: values are not the same %s", test.Name, cmp.Diff(vcard, test.ExpectedVCard))
		}
	}
}
