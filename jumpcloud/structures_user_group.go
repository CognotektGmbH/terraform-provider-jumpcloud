package jumpcloud

import (
	"fmt"
	"strconv"
	"strings"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
)

func flattenAttributes(attr *Attributes) map[string]interface{} {
	return map[string]interface{}{
		"ldap_groups":  flattenLDAPGroups(attr.LDAPGroups),
		"posix_groups": flattenPOSIXGroups(attr.POSIXGroups),
	}
}

func flattenPOSIXGroups(pg []jcapiv2.UserGroupPostAttributesPosixGroups) string {
	out := []string{}
	for _, v := range pg {
		out = append(out, fmt.Sprintf("%d:%s", v.Id, v.Name))
	}
	return strings.Join(out, ",")
}

func flattenLDAPGroups(lg []LDAPGroup) string {
	out := []string{}
	for _, v := range lg {
		out = append(out, v.Name)
	}
	return strings.Join(out, ",")
}

//	Note: PosixGroups cannot be edited after group creation, only first member of slice is considered
func expandAttributes(attr interface{}) (out *jcapiv2.UserGroupPostAttributes, ok bool) {
	if attr == nil {
		return
	}
	mapAttr, ok := attr.(map[string]interface{})
	if !ok {
		return
	}
	posixStr, ok := mapAttr["posix_groups"].(string)
	if !ok {
		return
	}

	groups := strings.Split(posixStr, ",")
	posixGroups := []jcapiv2.UserGroupPostAttributesPosixGroups{}
	for _, v := range groups {
		g := strings.Split(v, ":")
		if len(g) != 2 {
			return
		}
		id, err := strconv.ParseInt(g[0], 10, 32)
		if err != nil {
			continue
		}
		posixGroups = append(posixGroups,
			jcapiv2.UserGroupPostAttributesPosixGroups{
				Id: int32(id), Name: g[1],
			})
	}

	if len(posixGroups) == 0 {
		return
	}

	return &jcapiv2.UserGroupPostAttributes{PosixGroups: posixGroups}, true
}
