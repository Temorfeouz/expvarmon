package main

import "testing"

func TestUtils(t *testing.T) {
	str := "memstats.Alloc,memstats.Sys"

	vars, err := ParseVars(str)
	if err != nil {
		t.Fatalf("Err not nil: %v", err)
	}

	if len(vars) != 2 {
		t.Fatalf("vars should contain 2 elements, but has %d", len(vars))
	}

	str = "memstats.Alloc,memstats.Sys,goroutines,Counter.A"

	vars, err = ParseVars(str)
	if err != nil {
		t.Fatalf("Err not nil: %v", err)
	}

	if len(vars) != 4 {
		t.Fatalf("vars should contain 4 elements, but has %d", len(vars))
	}
}

func TestPorts(t *testing.T) {
	arg := "1234,1235"
	ports, err := ParsePorts(arg)
	if err != nil {
		t.Fatal(err)
	}
	if len(ports) != 2 || ports[0].Host != "localhost:1234" {
		t.Fatalf("ParsePorts returns wrong data: %v", ports)
	}

	arg = "1234-1237,2000"
	ports, err = ParsePorts(arg)
	if err != nil {
		t.Fatal(err)
	}
	if len(ports) != 5 || ports[0].Host != "localhost:1234" || ports[4].Host != "localhost:2000" {
		t.Fatalf("ParsePorts returns wrong data: %v", ports)
	}

	arg = "40000-40002,localhost:2000-2002,remote:1234-1235,https://example.com:1234-1236"
	ports, err = ParsePorts(arg)
	if err != nil {
		t.Fatal(err)
	}
	if len(ports) != 11 ||
		ports[0].Host != "localhost:40000" ||
		ports[3].Host != "localhost:2000" ||
		ports[7].Host != "remote:1235" ||
		ports[7].Path != "/debug/vars" ||
		ports[10].Host != "example.com:1236" ||
		ports[10].Scheme != "https" {
		t.Fatalf("ParsePorts returns wrong data: %v", ports)
	}

	// Test Auth
	arg = "http://user:pass@localhost:2000-2002"
	ports, err = ParsePorts(arg)
	if err != nil {
		t.Fatal(err)
	}
	pass, isSet := ports[0].User.Password()
	if len(ports) != 3 ||
		ports[0].User.Username() != "user" ||
		pass != "pass" || !isSet ||
		ports[0].Scheme != "http" {
		t.Fatalf("ParsePorts returns wrong data: %v", ports)
	}

	arg = "localhost:2000-2002,remote:1234-1235,some:weird:1234-123input"
	_, err = ParsePorts(arg)
	if err == nil {
		t.Fatalf("err shouldn't be nil")
	}

	arg = "string,sdasd"
	_, err = ParsePorts(arg)
	if err == nil {
		t.Fatalf("err shouldn't be nil")
	}
}
