package jumpcloud

import (
	jcapiv1 "github.com/TheJumpCloud/jcapi-go/v1"
)

func flattenPhoneNumbers(pn []jcapiv1.SystemuserreturnPhoneNumbers)  []interface{} {
	if pn == nil {
		return make([]interface{}, 0)
	}

	phoneNumbers := make([]interface{}, 0)
	for _, v := range pn {
		phoneNumbers = append(phoneNumbers, map[string]interface{}{
			"type": v.Type_,
			"number": v.Number,
		})
	}
	return phoneNumbers
}

func expandPhoneNumbers(input []interface{}) []map[string]string {
	if input == nil || len(input) == 0 {
		return nil
	}

	var phoneNumbers []map[string]string

	for _, v := range input {
		if phoneNumber, ok := v.(map[string]interface{}); ok {
			phoneNumbers = append(phoneNumbers, map[string]string{
				"number": phoneNumber["number"].(string),
				"type":  phoneNumber["type"].(string),
			})
		}

	}
	return phoneNumbers
}
