package main

import (
	"net"
	"reflect"
	"testing"
)

// TestConvertToCIDR tests the convertToCIDR function.
func TestConvertToCIDR(t *testing.T) {
	testCases := []struct {
		name     string
		record   Record
		expected []string
	}{
		{
			name: "768 from 192.168.0.0",
			record: Record{
				StartIP:  net.ParseIP("192.168.0.0"),
				NumHosts: 768,
			},
			expected: []string{"192.168.0.0/23", "192.168.2.0/24"},
		},
		{
			name: "512 from 192.168.3.0",
			record: Record{
				StartIP:  net.ParseIP("192.168.3.0"),
				NumHosts: 512,
			},
			expected: []string{"192.168.3.0/24", "192.168.4.0/24"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := convertToCIDR(tc.record)
			if err != nil {
				t.Errorf("convertToCIDR(%v) returned an error: %v", tc.record, err)
			}
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("convertToCIDR(%v) = %v, want %v", tc.record, got, tc.expected)
			}
		})
	}
}
