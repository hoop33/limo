package serr

import (
	"errors"
	"fmt"
	"testing"
)

func TestSErr(t *testing.T) {
	const strErr1 = "Ok. This is a test err"
	// todo: constantize most test literals

	fmt.Println("Testing SErr")

	// We should safely ignore a nil err
	ret := Wrap(nil, "We should be able to handle a nil error without crashing")
	if _, ok := ret.(SErr); !ok {
		t.Error("We should return a generated error when a nil error is wrapped")
	}

	ser := NewSErr(strErr1, "thing1", "thing1val", "thing2", "thing2val")
	if ser.Error() != strErr1 {
		t.Errorf("Expected custom error to contain '%s', got '%s'", strErr1, ser.Error())
		t.FailNow()
	}
	if _, ok := ser.(SErr); !ok {
		t.Error("ser should be a SErr")
		t.FailNow()
	}

	// Add some fields to an existing sErr
	err := Wrap(ser, "thing2", "thing2NewVal")
	se, ok := err.(SErr)
	if !ok {
		t.Error("Wrap should return an error containing a concrete SErr type")
	} else {
		// Test SErr#OriginalErr
		strErr := se.OriginalErr().Error()
		if strErr != strErr1 {
			t.Errorf(`Expected wrapped error string to be "%s", got "%s"`, strErr1, strErr)
		}
		// Test SErr#FieldsSlice
		strFlds := se.FieldsSlice()
		fmt.Printf("[Debug] strFlds: %#v\n", strFlds) // debug
		if len(strFlds) != 8 {
			t.Error("Expected length of SErr.Fields() to be 8, got", len(strFlds))
		}
		// Test SErr#FieldsMap
		mapFlds := se.FieldsMap()
		if len(mapFlds) != 3 {
			t.Error("Expected length of SErr.MapFlds() to be 3, got", len(mapFlds))
			t.FailNow()
		}

		if val, ok := mapFlds["thing1"]; ok {
			if val != "thing1val" {
				t.Errorf("Expected thing1 to be 'thing1Val', got '%s'", val)
			}
		} else {
			t.Error("mapFlds should contain the key: thing1")
		}

		if val, ok := mapFlds["thing2"]; ok {
			if val != "thing2NewVal - thing2val" {
				t.Errorf("Expected thing2 to be 'thing2NewVal - thing2val', got '%s'", val)
			}
		} else {
			t.Error("mapFlds should contain the key: thing2")
		}
	}

	// We should be able to wrap with a single field which becomes `"msg": field`
	const thisIsMyMessage = "This is my message"
	er := Wrap(errors.New(strErr1), thisIsMyMessage)
	se, ok = er.(SErr)
	if !ok {
		t.Error("er should be a SErr")
		t.FailNow()
	}
	sl := se.FieldsSlice()
	if len(sl) != 4 {
		t.Error("Structured error from an error wrapped with a single field should contain 4 fields, got", len(sl))
	}
	if len(sl) > 0 && sl[0] != "msg" {
		t.Error("Structured error from an error wrapped with a single field should have 'msg' as the first field, got", sl[0])
	}
	if len(sl) > 1 && sl[1] != thisIsMyMessage {
		t.Errorf(`The structured error should have "%s" as the second field, got "%s"`, thisIsMyMessage, sl[1])
		fmt.Println()
	}

	// We should be able to wrap with an odd number of fields
	er2 := Wrap(errors.New(strErr1), "Wrapping odd number of fields", "key1", "value1", "key2", "value2")
	se, ok = er2.(SErr)
	if !ok {
		t.Error("er should be a SErr")
		t.FailNow()
	}
	arrFlds := se.FieldsSlice()
	if ln := len(arrFlds); ln != 8 {
		t.Error("Odd number of fields should use first field as a message. Expected 8 items in slice, got", ln)
		t.FailNow()
	}
	if len(arrFlds) > 0 && arrFlds[0] != "msg" {
		t.Error("Structured error from an error wrapped with a single field should have 'msg' as the first field, got", arrFlds[0])
	}
}
