package badgercache

import "testing"

func TestIsJSON(t *testing.T) {
	if !IsJSON(`{"name":"hello"}`) {
		t.Errorf("Should be json")
	}

	if !IsJSON(`[{"name":"hello"}]`) {
		t.Errorf("Should be json")
	}

	if IsJSON(`[{"name":"hello"}`) {
		t.Errorf("Should not be json")
	}
}

func TestIsYAML(t *testing.T) {
	var d = `foo: 1
hello:
- one
- two`

	if !IsYAML(d) {
		t.Errorf("Should be yaml")
	}

	d = `--- !test.com
version: "some_version"
list:
  - val1:
    sub: 1
    end: true
  - val2
  - val3:
    sub: 3
    end: false
`

	if !IsYAML(d) {
		t.Errorf("Should be yaml")
	}
}
