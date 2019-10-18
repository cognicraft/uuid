package uuid

import "testing"

var (
	validTestCases = []struct {
		name    string
		uuid    UUID
		version uint8
		variant Variant
		str     string
	}{
		{"nil", Nil, 0, 0, "00000000-0000-0000-0000-000000000000"},
		{"ns:DNS", NamespaceDNS, 1, VariantRFC4122, "6ba7b810-9dad-11d1-80b4-00c04fd430c8"},
		{"ns:URL", NamespaceURL, 1, VariantRFC4122, "6ba7b811-9dad-11d1-80b4-00c04fd430c8"},
		{"ns:OID", NamespaceOID, 1, VariantRFC4122, "6ba7b812-9dad-11d1-80b4-00c04fd430c8"},
		{"ns:X500", NamespaceX500, 1, VariantRFC4122, "6ba7b814-9dad-11d1-80b4-00c04fd430c8"},
		{"v3:DNS:www.example.org", NewV3(NamespaceDNS, "www.example.org"), 3, VariantRFC4122, "0012416f-9eec-3ed4-a8b0-3bceecde1cd9"},
		{"v3:URL:www.example.org", NewV3(NamespaceURL, "www.example.org"), 3, VariantRFC4122, "288ea4f1-895d-3984-98a3-85aa3aa6db56"},
		{"v5:DNS:www.example.org", NewV5(NamespaceDNS, "www.example.org"), 5, VariantRFC4122, "74738ff5-5367-5958-9aee-98fffdcd1876"},
	}
)

func TestString(t *testing.T) {
	for _, tc := range validTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.uuid.String()
			if tc.str != got {
				t.Errorf("u.String() expected: %s, got: %s", tc.str, got)
			}
		})
	}
}

func TestParse(t *testing.T) {
	for _, tc := range validTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Parse(tc.str)
			if err != nil {
				t.Errorf("Parse() error: %s", err)
			}
			if !Equal(tc.uuid, got) {
				t.Errorf("Parse() expected: %s, got: %s", tc.uuid, got)
			}
		})
	}
	errorCases := []struct {
		name string
		str  string
	}{
		{"invalid-string-format", "6ba7b8109dad11d180b400c04fd430c8"},
		{"string-too-short", "6ba7"},
		{"string-too-long", "0012416f-9eec-3ed4-a8b0-3bceecde1cd9ab"},
	}
	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.str)
			if err == nil {
				t.Errorf("expected error")
			}
		})
	}
}

func TestVersion(t *testing.T) {
	for _, tc := range validTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.uuid.Version()
			if tc.version != got {
				t.Errorf("u.Version(): expected: %d, got %d", tc.version, got)
			}
		})
	}
}

func TestVariant(t *testing.T) {
	for _, tc := range validTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.uuid.Variant()
			if tc.variant != got {
				t.Errorf("u.Variant(): expected: %d, got %d", tc.variant, got)
			}
		})
	}
}
