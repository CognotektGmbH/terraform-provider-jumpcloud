package jumpcloud

import jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"

// UserGroup is like jcapiv2.UserGroup with Attributes
type UserGroup struct {
	// ID uniquely identifies a User Group.
	ID string `json:"id,omitempty"`

	// Type is the type of the group.
	Type string `json:"type,omitempty"`

	// Display name of a User Group.
	Name       string                      `json:"name,omitempty"`
	Attributes jcapiv2.UserGroupAttributes `json:"attributes,omitempty"`
}
