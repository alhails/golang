// Code generated by "go-enum -type ConnState -trimprefix=ConnState"; DO NOT EDIT.

package cmux

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"strconv"
)

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ConnStateNew-0]
	_ = x[ConnStateActive-1]
	_ = x[ConnStateIdle-2]
	_ = x[ConnStateHijacked-3]
	_ = x[ConnStateClosed-4]
}

const _ConnState_name = "NewActiveIdleHijackedClosed"

var _ConnState_index = [...]uint8{0, 3, 9, 13, 21, 27}

func _() {
	var _nil_ConnState_value = func() (val ConnState) { return }()

	// An "cannot convert ConnState literal (type ConnState) to type fmt.Stringer" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ fmt.Stringer = _nil_ConnState_value
}

func (i ConnState) String() string {
	if i < 0 || i >= ConnState(len(_ConnState_index)-1) {
		return "ConnState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ConnState_name[_ConnState_index[i]:_ConnState_index[i+1]]
}

var _ConnState_values = []ConnState{0, 1, 2, 3, 4}

var _ConnState_name_to_values = map[string]ConnState{
	_ConnState_name[0:3]:   0,
	_ConnState_name[3:9]:   1,
	_ConnState_name[9:13]:  2,
	_ConnState_name[13:21]: 3,
	_ConnState_name[21:27]: 4,
}

// ParseConnStateString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func ParseConnStateString(s string) (ConnState, error) {
	if val, ok := _ConnState_name_to_values[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%[1]s does not belong to ConnState values", s)
}

// ConnStateValues returns all values of the enum
func ConnStateValues() []ConnState {
	return _ConnState_values
}

// IsAConnState returns "true" if the value is listed in the enum definition. "false" otherwise
func (i ConnState) Registered() bool {
	for _, v := range _ConnState_values {
		if i == v {
			return true
		}
	}
	return false
}

func _() {
	var _nil_ConnState_value = func() (val ConnState) { return }()

	// An "cannot convert ConnState literal (type ConnState) to type encoding.BinaryMarshaler" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ encoding.BinaryMarshaler = &_nil_ConnState_value

	// An "cannot convert ConnState literal (type ConnState) to type encoding.BinaryUnmarshaler" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ encoding.BinaryUnmarshaler = &_nil_ConnState_value
}

// MarshalBinary implements the encoding.BinaryMarshaler interface for ConnState
func (i ConnState) MarshalBinary() (data []byte, err error) {
	return []byte(i.String()), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for ConnState
func (i *ConnState) UnmarshalBinary(data []byte) error {
	var err error
	*i, err = ParseConnStateString(string(data))
	return err
}

func _() {
	var _nil_ConnState_value = func() (val ConnState) { return }()

	// An "cannot convert ConnState literal (type ConnState) to type json.Marshaler" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ json.Marshaler = _nil_ConnState_value

	// An "cannot convert ConnState literal (type ConnState) to type encoding.Unmarshaler" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ json.Unmarshaler = &_nil_ConnState_value
}

// MarshalJSON implements the json.Marshaler interface for ConnState
func (i ConnState) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for ConnState
func (i *ConnState) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("ConnState should be a string, got %[1]s", data)
	}

	var err error
	*i, err = ParseConnStateString(s)
	return err
}

func _() {
	var _nil_ConnState_value = func() (val ConnState) { return }()

	// An "cannot convert ConnState literal (type ConnState) to type encoding.TextMarshaler" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ encoding.TextMarshaler = _nil_ConnState_value

	// An "cannot convert ConnState literal (type ConnState) to type encoding.TextUnmarshaler" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ encoding.TextUnmarshaler = &_nil_ConnState_value
}

// MarshalText implements the encoding.TextMarshaler interface for ConnState
func (i ConnState) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for ConnState
func (i *ConnState) UnmarshalText(text []byte) error {
	var err error
	*i, err = ParseConnStateString(string(text))
	return err
}

//func _() {
//	var _nil_ConnState_value = func() (val ConnState) { return }()
//
//	// An "cannot convert ConnState literal (type ConnState) to type yaml.Marshaler" compiler error signifies that the base type have changed.
//	// Re-run the go-enum command to generate them again.
//	var _ yaml.Marshaler = _nil_ConnState_value
//
//	// An "cannot convert ConnState literal (type ConnState) to type yaml.Unmarshaler" compiler error signifies that the base type have changed.
//	// Re-run the go-enum command to generate them again.
//	var _ yaml.Unmarshaler = &_nil_ConnState_value
//}

// MarshalYAML implements a YAML Marshaler for ConnState
func (i ConnState) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for ConnState
func (i *ConnState) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	var err error
	*i, err = ParseConnStateString(s)
	return err
}

func _() {
	var _nil_ConnState_value = func() (val ConnState) { return }()

	// An "cannot convert ConnState literal (type ConnState) to type driver.Valuer" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ driver.Valuer = _nil_ConnState_value

	// An "cannot convert ConnState literal (type ConnState) to type sql.Scanner" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ sql.Scanner = &_nil_ConnState_value
}

func (i ConnState) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *ConnState) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		bytes, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("value is not a byte slice")
		}

		str = string(bytes[:])
	}

	val, err := ParseConnStateString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}

// ConnStateSliceContains reports whether sunEnums is within enums.
func ConnStateSliceContains(enums []ConnState, sunEnums ...ConnState) bool {
	var seenEnums = map[ConnState]bool{}
	for _, e := range sunEnums {
		seenEnums[e] = false
	}

	for _, v := range enums {
		if _, has := seenEnums[v]; has {
			seenEnums[v] = true
		}
	}

	for _, seen := range seenEnums {
		if !seen {
			return false
		}
	}

	return true
}

// ConnStateSliceContainsAny reports whether any sunEnum is within enums.
func ConnStateSliceContainsAny(enums []ConnState, sunEnums ...ConnState) bool {
	var seenEnums = map[ConnState]struct{}{}
	for _, e := range sunEnums {
		seenEnums[e] = struct{}{}
	}

	for _, v := range enums {
		if _, has := seenEnums[v]; has {
			return true
		}
	}

	return false
}
