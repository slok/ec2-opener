package rule

import "testing"

func TestValidProtocol(t *testing.T) {
	tests := []struct {
		proto string
		valid bool
		want  Protocol
	}{
		{proto: "udp", valid: true, want: UDP},
		{proto: "UDP", valid: true, want: UDP},
		{proto: "tcp", valid: true, want: TCP},
		{proto: "TCP", valid: true, want: TCP},
		{proto: "icmp", valid: true, want: ICMP},
		{proto: "ICMP", valid: true, want: ICMP},
		{proto: "all", valid: true, want: ALL},
		{proto: "ALL", valid: true, want: ALL},
		{proto: "", valid: false},
		{proto: "1234", valid: false},
		{proto: "UDp", valid: false},
	}

	for _, test := range tests {

		parsedProto, err := ParseProtocol(test.proto)

		// Check for errors
		if err != nil && test.valid {
			t.Errorf("For '%s' protocol, got error: %s", test.proto, err)
		} else { // if no error check result is what we want
			if parsedProto != test.want {
				t.Errorf("For '%s' protocol, got '%s'; want '%s'", test.proto, parsedProto, test.want)
			}
		}
	}
}
