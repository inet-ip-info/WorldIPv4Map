package main

import (
	"net"
	"reflect"
	"testing"
)

type Record struct {
	StartIP  net.IP
	NumHosts int
}

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
		{
			name: "256 from 192.168.5.0",
			record: Record{
				StartIP:  net.ParseIP("192.168.5.0"),
				NumHosts: 256,
			},
			expected: []string{"192.168.5.0/24"},
		},
		{
			name: "4 from 192.168.6.0",
			record: Record{
				StartIP:  net.ParseIP("192.168.6.0"),
				NumHosts: 4,
			},
			expected: []string{"192.168.6.0/30"},
		},
		{
			name: "2048 from 192.168.10.0",
			record: Record{
				StartIP:  net.ParseIP("192.168.10.0"),
				NumHosts: 2048,
			},
			//expected: []string{"192.168.10.0/21"},
			expected: []string{"192.168.10.0/23", "192.168.12.0/22", "192.168.16.0/23"},
		},
		{
			name: "1 from 192.168.100.0",
			record: Record{
				StartIP:  net.ParseIP("192.168.100.0"),
				NumHosts: 1,
			},
			expected: []string{"192.168.100.0/32"},
		},
		{
			name: "256 from 192.168.100.0",
			record: Record{
				StartIP:  net.ParseIP("192.168.100.0"),
				NumHosts: 256,
			},
			expected: []string{"192.168.100.0/24"},
		},
		{
			name: "257 from 192.168.100.0",
			record: Record{
				StartIP:  net.ParseIP("192.168.100.0"),
				NumHosts: 257,
			},
			expected: []string{"192.168.100.0/24", "192.168.101.0/32"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := convertToCIDR(tc.record.StartIP, tc.record.NumHosts)
			if err != nil {
				t.Errorf("convertToCIDR(%v) returned an error: %v", tc.record, err)
			}
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("convertToCIDR(%v) = %v, want %v", tc.record, got, tc.expected)
			}
		})
	}
}

func TestConvertCIDRToRanges(t *testing.T) {
	testCases := []struct {
		name     string
		cidrs    []string
		expected []IPRange
		wantErr  bool
	}{
		{
			name:  "Adjacent CIDR Blocks",
			cidrs: []string{"192.168.0.0/24", "192.168.1.0/24"},
			expected: []IPRange{
				{StartIP: net.ParseIP("192.168.0.0"), NumHosts: 512},
			},
			wantErr: false,
		},
		{
			name:  "Single CIDR Block",
			cidrs: []string{"192.168.0.0/24"},
			expected: []IPRange{
				{StartIP: net.ParseIP("192.168.0.0"), NumHosts: 256},
			},
			wantErr: false,
		},
		{
			name:  "Non-Adjacent CIDR Blocks",
			cidrs: []string{"192.168.0.0/24", "192.168.2.0/24"},
			expected: []IPRange{
				{StartIP: net.ParseIP("192.168.0.0"), NumHosts: 256},
				{StartIP: net.ParseIP("192.168.2.0"), NumHosts: 256},
			},
			wantErr: false,
		},
		{
			name:  "Combining Different Sizes of CIDR Blocks",
			cidrs: []string{"192.168.0.0/25", "192.168.0.128/26", "192.168.0.192/26"},
			expected: []IPRange{
				{StartIP: net.ParseIP("192.168.0.0"), NumHosts: 256}, // 192.168.0.0/24と同等
			},
			wantErr: false,
		},
		{
			name:  "Combining Different Sizes of CIDR Blocks into Multiple Ranges",
			cidrs: []string{"192.168.0.0/25", "192.168.0.128/26", "192.168.1.0/24", "192.168.3.0/23"},
			expected: []IPRange{
				{StartIP: net.ParseIP("192.168.0.0"), NumHosts: 192}, // 192.168.0.0/25 と 192.168.0.128/26 を結合
				{StartIP: net.ParseIP("192.168.1.0"), NumHosts: 256}, // 192.168.1.0/24
				{StartIP: net.ParseIP("192.168.3.0"), NumHosts: 512}, // 192.168.3.0/23
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := convertCIDRToRanges(tc.cidrs)
			if (err != nil) != tc.wantErr {
				t.Errorf("convertCIDRToRanges() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("convertCIDRToRanges() = %v, want %v", got, tc.expected)
			}
		})
	}
}
