package main

import "testing"

func TestResolveListenAddress(t *testing.T) {
	tests := []struct {
		name       string
		configured string
		engineMode string
		want       string
		wantErr    bool
	}{
		{name: "fixture default", want: ":8080"},
		{name: "lab default", engineMode: "nmap-lab", want: "127.0.0.1:8080"},
		{name: "lab explicit loopback", configured: "127.0.0.1:9090", engineMode: "nmap-lab", want: "127.0.0.1:9090"},
		{name: "lab explicit ipv6 loopback", configured: "[::1]:9090", engineMode: "nmap-lab", want: "[::1]:9090"},
		{name: "lab rejects wildcard", configured: ":8080", engineMode: "nmap-lab", wantErr: true},
		{name: "lab rejects external address", configured: "192.0.2.10:8080", engineMode: "nmap-lab", wantErr: true},
		{name: "fixture accepts wildcard", configured: ":9090", want: ":9090"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := resolveListenAddress(test.configured, test.engineMode)
			if (err != nil) != test.wantErr {
				t.Fatalf("resolveListenAddress() error = %v, wantErr %t", err, test.wantErr)
			}
			if err == nil && got != test.want {
				t.Fatalf("resolveListenAddress() = %q, want %q", got, test.want)
			}
		})
	}
}
