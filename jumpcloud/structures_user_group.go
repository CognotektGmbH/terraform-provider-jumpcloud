package jumpcloud

// see https://www.terraform.io/docs/extend/writing-custom-providers.html#implementing-a-more-complex-read

import (
	"fmt"
	"strconv"
	"strings"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
)

func flattenAttributes(attr *jcapiv2.UserGroupPostAttributes) map[string]interface{} {
	return map[string]interface{}{
		"posix_groups": flattenPosixGroups(attr.PosixGroups),
		// "enable_samba": fmt.Sprintf("%t", attr.SambaEnabled),
	}
}

func flattenPosixGroups(pg []jcapiv2.UserGroupPostAttributesPosixGroups) string {
	out := []string{}
	for _, v := range pg {
		out = append(out, fmt.Sprintf("%d:%s", v.Id, v.Name))
	}
	return strings.Join(out, ",")
}

func expandAttributes(attr interface{}) (out *jcapiv2.UserGroupPostAttributes, ok bool) {
	if attr == nil {
		return
	}
	mapAttr, ok := attr.(map[string]interface{})
	if !ok {
		return
	}

	// var enableSamba bool
	// sambaStr, ok := mapAttr["enable_samba"].(string)
	// if ok {
	// 	enableSamba, _ = strconv.ParseBool(sambaStr)
	// }

	// TODO: empty string? nil?
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

	return &jcapiv2.UserGroupPostAttributes{
		PosixGroups: posixGroups,
		// SambaEnabled: enableSamba,
	}, true
}
