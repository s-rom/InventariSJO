package dbsqlc

import "encoding/json"

// MarshalJSON serializes NullConnectionTypeEnum as a plain string or JSON null.
func (ns NullConnectionTypeEnum) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(string(ns.ConnectionTypeEnum))
}

// UnmarshalJSON deserializes NullConnectionTypeEnum from a plain string or JSON null.
func (ns *NullConnectionTypeEnum) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	ns.ConnectionTypeEnum = ConnectionTypeEnum(s)
	ns.Valid = true
	return nil
}
