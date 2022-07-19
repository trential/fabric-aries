package utils

import (
	"crypto/ed25519"
	"encoding/json"
	"reflect"
	"sort"
	"strings"

	"github.com/btcsuite/btcutil/base58"
)

func VerifyRequest(verKey, message, signature string) bool {
	pub := ed25519.PublicKey(base58.Decode(verKey))
	sig := base58.Decode(signature)
	return ed25519.Verify(pub, []byte(serialize(message)), sig)
}

func serialize(msg string) string {
	d := json.NewDecoder(strings.NewReader(msg))
	d.UseNumber()
	var object interface{}
	if err := d.Decode(&object); err != nil {
		return ""
	}
	return serializeInternal(object, true)
}

func serializeInternal(v interface{}, isTopLevel bool) string {
	if v == nil {
		return ""
	}
	switch reflect.TypeOf(v).Kind() {
	case reflect.Bool:
		if v.(bool) {
			return "True"
		} else {
			return "False"

		}
	case reflect.String:
		if _, ok := v.(json.Number); ok {
			return v.(json.Number).String()
		}
		return v.(string)
	case reflect.Slice:
		out := ""
		for i, vv := range v.([]interface{}) {
			out += serializeInternal(vv, isTopLevel)
			if i != len(v.([]interface{}))-1 {
				out += ","
			}
		}
		return out
	case reflect.Map:
		out := ""
		inMiddle := false

		m := v.(map[string]interface{})
		keys := []string{}
		for k, _ := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if k == "signature" && isTopLevel {
				continue
			}
			if inMiddle {
				out += "|"
			}
			out = out + k + ":" + serializeInternal(m[k], false)
			inMiddle = true
		}
		return out
	}

	return ""
}
