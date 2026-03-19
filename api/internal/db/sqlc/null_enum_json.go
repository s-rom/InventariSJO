package dbsqlc

import "encoding/json"

// MarshalJSON serializes NullRamTypeEnum as a plain string or JSON null.
func (ns NullRamTypeEnum) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(string(ns.RamTypeEnum))
}

// UnmarshalJSON deserializes NullRamTypeEnum from a plain string or JSON null.
func (ns *NullRamTypeEnum) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		ns.Valid = false
		ns.RamTypeEnum = ""
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	ns.RamTypeEnum = RamTypeEnum(s)
	ns.Valid = true
	return nil
}

// MarshalJSON serializes NullStorageTypeEnum as a plain string or JSON null.
func (ns NullStorageTypeEnum) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(string(ns.StorageTypeEnum))
}

// UnmarshalJSON deserializes NullStorageTypeEnum from a plain string or JSON null.
func (ns *NullStorageTypeEnum) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		ns.Valid = false
		ns.StorageTypeEnum = ""
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	ns.StorageTypeEnum = StorageTypeEnum(s)
	ns.Valid = true
	return nil
}

// MarshalJSON serializes NullAuditEventEnum as a plain string or JSON null.
func (ns NullAuditEventEnum) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(string(ns.AuditEventEnum))
}

// UnmarshalJSON deserializes NullAuditEventEnum from a plain string or JSON null.
func (ns *NullAuditEventEnum) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		ns.Valid = false
		ns.AuditEventEnum = ""
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	ns.AuditEventEnum = AuditEventEnum(s)
	ns.Valid = true
	return nil
}

// MarshalJSON serializes NullShiftEnum as a plain string or JSON null.
func (ns NullShiftEnum) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(string(ns.ShiftEnum))
}

// UnmarshalJSON deserializes NullShiftEnum from a plain string or JSON null.
func (ns *NullShiftEnum) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		ns.Valid = false
		ns.ShiftEnum = ""
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	ns.ShiftEnum = ShiftEnum(s)
	ns.Valid = true
	return nil
}
