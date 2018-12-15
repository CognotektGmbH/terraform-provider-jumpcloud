package jumpcloud

import jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"

// UserGroup is like github.com/TheJumpCloud/jcapi-go/v2.UserGroup but with Attributes and go best practices
type UserGroup struct {
	// ID uniquely identifies a User Group.
	ID string `json:"id,omitempty"`

	// Type is the type of the group.
	Type string `json:"type,omitempty"`

	// Display name of a User Group.
	Name       string     `json:"name,omitempty"`
	Attributes Attributes `json:"attributes,omitempty"`
}

// Attributes holds UserGroup properties relevant for updates
type Attributes struct {
	POSIXGroups []jcapiv2.UserGroupPostAttributesPosixGroups `json:"posixGroups,omitempty"`
	// SambaEnabled bool                                         `json:"sambaEnable,omitempty"`
}
