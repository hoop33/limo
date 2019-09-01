package serr

import (
	"errors"
	"fmt"
)

// Structured Error wrapper
// Supports wrapping of errors with a list of key, values to nicely support structured logging
// Works nicely with github.com/rohanthewiz/logger.LogErr
type SErr struct {
	err error // the usual error
	// support structured logging of the format key1, val1, key2, val2
	// Repeated keys are allowed and will be concatenated in log output
	fields []string
}

// Create a new SErr (structured err) from an error string and
// embellish with string pairs of attributes and values
// Returns an error (SErr satisfies the error interface)
// This requires an even number of fields otherwise the first field is added under the key "msg".
// Examples are given in serr_test.go
func NewSErr(er string, fields ...string) error {
	se := SErr{err: errors.New(er)}
	se.fields = buildFields(fields)
	return se
}

// Yes we are also an error type -- sweet!
// Satisfy the `error` interface
func (se SErr) Error() string {
	if se.err == nil {
		return ""
	}
	return se.err.Error()
}

// Return all SErr attributes as a map of string keys and values
func (se SErr) FieldsMap() map[string]string {
	flds := map[string]string{}
	key := ""
	lenFields := len(se.fields)
	for i, str := range se.fields {
		if i%2 == 0 { // we should have a key
			if i == lenFields-1 { // A key should not be the last item
				warn := fmt.Sprintf(`[SErr] key: "%s" has no matching value. Location: %s, Fields: %#v`,
					str, FuncLoc(3), se.fields)
				fmt.Println(warn)
				flds["serrWarn"] = warn
				break // this the last item
			}
			key = str
		} else { // we have a value
			if orig, ok := flds[key]; ok { // we've seen this key before - accumulate
				flds[key] = str + " - " + orig
			} else {
				flds[key] = str
			}
		}
	}
	return flds
}

// Wrap an existing error with string fields of attribute key and value pairs.
// Returns an SErr (structured err)
// This requires an even number of fields unless a single field is given
// in which case it is added under the key "msg".
// Examples are given in serr_test.go
func Wrap(err error, fields ...string) error {
	var flds []string
	err = handleNilError(err)

	// Add Existing fields
	if se, ok := err.(SErr); ok && len(se.fields) > 0 {
		flds = append(flds, se.fields...)
	}
	// Add new fields
	flds = append(flds, buildFields(fields)...)
	// Add location
	flds = append(flds, []string{"location", FuncLoc(CallerIndirection.Parent)}...)

	return SErr{err, flds} // return
}

// SpecialVersion of wrap that takes callers indirection as a parameter
// callerIndir should at least be grandparent
func LogWrap(err error, callerIndir int, fields ...string) error {
	var flds []string
	err = handleNilError(err)
	// Add Existing fields
	if se, ok := err.(SErr); ok && len(se.fields) > 0 {
		flds = append(flds, se.fields...)
	}
	// Add new fields
	flds = append(flds, buildFields(fields)...)
	// Add location
	if callerIndir < CallerIndirection.GrandParent { callerIndir = CallerIndirection.GrandParent }
	flds = append(flds, []string{"location", FuncLoc(callerIndir)}...)
	return SErr{err, flds} // return
}

// Return the wrapped error
func (se SErr) OriginalErr() error {
	return se.err
}

// Return the internal list of keys and values
func (se SErr) FieldsSlice() []string {
	return se.fields
}

// Process given array of strings, return sequence of attribute value pairs
// A Single field gets added as {"msg", "field"}
// For an odd number of multiple fields, the first field is considered a message value {"msg", "field"}
// An even number of fields are added without any change in sequence
func buildFields(fields []string) (flds []string) {
	ln := len(fields)
	if ln == 1 { // Single field becomes a "msg: field" pair
		flds = append(flds, []string{"msg", fields[0]}...)
	} else {
		if ln % 2 != 0 { // for odd fields, first is a message
			msg := fields[0]
			fields = fields[1:]  // drop the first
			flds = append(flds, []string{"msg", msg}...)
		}
		flds = append(flds, fields...)
	}
	return
}

func handleNilError(err error) error {
	if err == nil {
		warn := fmt.Sprintf(`[SErr] nil error provided at %s. That is weird since this is an err function`,
			FuncLoc(CallerIndirection.GrandParent))
		err = errors.New(warn)
		fmt.Println(warn)
	}
	return err
}
