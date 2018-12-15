package jumpcloud

import jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"

type UserGroup struct {
	// ObjectId uniquely identifying a User Group.
	Id string `json:"id,omitempty"`

	// The type of the group.
	Type_ string `json:"type,omitempty"`

	// Display name of a User Group.
	Name       string     `json:"name,omitempty"`
	Attributes Attributes `json:"attributes,omitempty"`
}

type Attributes struct {
	LDAPGroups  []LDAPGroup                                  `json:"ldapGroups,omitempty"`
	POSIXGroups []jcapiv2.UserGroupPostAttributesPosixGroups `json:"posixGroups,omitempty"`
	// SambaEnabled bool                                         `json:"sambaEnable,omitempty"`
}

type LDAPGroup struct {
	Name string `json:"name,omitempty"`
}
