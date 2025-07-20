//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"syscall/js"

	"github.com/xtls/xray-core/common/errors"
	"github.com/xtls/xray-core/common/platform/filesystem"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/infra/conf"
	json_reader "github.com/xtls/xray-core/infra/conf/json"
)

//go:embed assets/geoip.dat
var geoipRaw []byte

//go:embed assets/geosite.dat
var geositeRaw []byte

func main() {
	filesystem.NewFileReader = func(path string) (io.ReadCloser, error) {
		if strings.HasSuffix(path, "geoip.dat") {
			return io.NopCloser(bytes.NewReader(geoipRaw)), nil
		}

		if strings.HasSuffix(path, "geosite.dat") {
			return io.NopCloser(bytes.NewReader(geositeRaw)), nil
		}

		return nil, errors.New(path + " cannot be opened in the browser")
	}

	js.Global().Set("XrayGetVersion", js.FuncOf(func(this js.Value, args []js.Value) any {
		return core.Version()
	}))

	js.Global().Set("XrayParseConfig", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			fmt.Println("invalid number of args")
			return nil
		}

		arg := args[0]
		if arg.Type() != js.TypeString {
			fmt.Println("the argument should be a string")
			return nil
		}

		var configObj interface{}
		if err := json.Unmarshal([]byte(arg.String()), &configObj); err != nil {
			return err.Error()
		}

		replaceKeyFileAndCertificateFile(configObj)
		removeExtDatRules(configObj)

		modifiedConfigBytes, err := json.Marshal(configObj)
		if err != nil {
			return err.Error()
		}

		jsonConfig := &conf.Config{}
		jsonReader := &json_reader.Reader{
			Reader: bytes.NewReader(modifiedConfigBytes),
		}
		decoder := json.NewDecoder(jsonReader)

		if err := decoder.Decode(jsonConfig); err != nil {
			return err.Error()
		}

		if _, err := jsonConfig.Build(); err != nil {
			return err.Error()
		}

		return nil
	}))

	js.Global().Get("onWasmInitialized").Invoke()

	// Prevent the program from exiting.
	// Note: the exported func should be released if you don't need it any more,
	// and let the program exit after then. To simplify this demo, this is
	// omitted. See https://pkg.go.dev/syscall/js#Func.Release for more information.
	select {}
}

func replaceKeyFileAndCertificateFile(obj interface{}) {
	switch v := obj.(type) {
	case map[string]interface{}:
		for key, value := range v {
			if key == "keyFile" {
				fmt.Printf("Replacing keyFile: %v with mocked key\n", value)
				v["key"] = []string{`-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCw++pBy6bPL23L
lcJVeVlFp7nv1is6pkmFxUEr2ipvEejn4ueiwMPAi4+/PEQ28qfDHRVTqsWSRf7H
MCDWNuG81Rmjfz2d27lHKDHUthGyUglcz8yu2pFGv/sP+jZTwbUCSasSs3hwOald
Kl6Z9HXxi1AkK8r3/Fzb9MDPagzJnNKHTRIjYdAGnzgJHb7XZdulesYLzmeZ6+yw
TwK7sCnP8FDZTtzdRZFf+sugxewfBjrwkByC1w134KsrJzKsa1h3KxR+QyG3dcV2
zylIg8l/I+kF6MoItlg5kr2+zbIgMAAqgo5U1O/Yv0ZrgkYkq1r7Pxq/xyCBzSVR
hsulnI7HAgMBAAECggEBAJeD8zbEzFfPOOXznd1i9+REBVCoP3YEaikViAeszGsu
IAB1Ju0UrgEm+bc0Nwn7j1fgjCOGrYHeXpHUSChb78GjhkGMawm709BxPsfl3xyU
UuOcGpUPFRRGwv8YrG3kBbyhqM1nzBH3q2DRJxASUu6o38E3pJnM+QptbKulFLF8
HdOEEg3gnbNd1T1wpgczkiQhWNCB6K5QvRjrsw8G9lddLO8YmieMUKyk45Enrv93
WDAvZQFDSlpj0r7qGunG6oEEdc+ZKvqq5I7dxEqGGg2HXI65ib3UwiUOsdxolEtv
O9xFKg2qqbUveTXpbNv8NuM2czjSyHR6M4UpC3zf4DECgYEA2MoUsU5ZKywiLJmj
JpkBkYENY9nEBEVW9d+ebfIxpo4GYAsTlryox4X4448Twf1wcz5MGZXwmhJl2Syc
Dp1WIodJvZ/xxNGkAUhr1oJLm4W2EXry6Xp0uhNVafB4ekaYJkm5Vb8N2IEWzDWo
dU4/MEcTL1Vs2ktp8j9v8BlS9P0CgYEA0P6+DUBNTDLRLgZVnHOIJ/o7oGe9yE3Q
OB6Sr5XqklcwEcQWBIEvtZj4rfm3o6RoaxnoJdrpjiLFXKwlDcGr2sN9Lqks0kfm
2ZdGpT2plKE0oh8Kdl045yemcyy51WmIfKBBHInR5hL25p2wHhHxnYib6ng940tk
P7sTA0Wj4BMCgYAOoy5cfmbE5Hj2O/VpIMGbWnRV/pkelLP3a/7de5HgpgxGJdlP
vzNCLYiNjNaPrZYPIfCvdZFGReG8lSeAUR1EvY+8DvWbDXCeaY5mcGu8d33AlmWa
YBtLiQymV4c68lNJmoa6BGlh6e1pImJacUsQ7mucMY9k+dyQb7oWIw+X3QKBgQDJ
TaWXs+kATS5Iy5cok+uAvjkyntohFjpJ48DcWWVQwaQuaJXgjuJ6YzactJwahiCB
kLmXxM0TuBAr3C/wmSxEEeoAyLjAbrs/uMM2JDe0TrYYthdovRAzLnDYHSt+ESGD
EQTGTUWc+4VPynE59YSpfUzYwiqyRQrxi+qbSze2ewKBgDGt21sta++EbNv9F/Wl
X4dAB3/G4SMQo/teUhlK1iXgZXjLeOKkuOq8HErLU2pJer2s7fKXNP/9BpRX7pEE
08bG0JRxlgYpFgPVnEhF9NTNEbS2wR+8Z7Sgm36WF8eHrgeGmRz2IWBMsZpfeJ/a
7BHM8pzGDJJG/PiONMYk5fNj
-----END PRIVATE KEY-----
`}
				delete(v, "keyFile")
			} else if key == "certificateFile" {
				fmt.Printf("Replacing certificateFile: %v with mocked certificate\n", value)
				v["certificate"] = []string{`-----BEGIN CERTIFICATE-----
MIIDcDCCAligAwIBAgIEGg8C6zANBgkqhkiG9w0BAQsFADBRMQswCQYDVQQGEwJK
UDEPMA0GA1UEAwwGTW9ja2VkMQ8wDQYDVQQIDAZNb2NrZWQxDzANBgNVBAcMBk1v
Y2tlZDEPMA0GA1UECgwGTW9ja2VkMB4XDTI0MTEyOTIzMTE1MVoXDTM0MTEyNzIz
MTE1MVowUTELMAkGA1UEBhMCSlAxDzANBgNVBAMMBk1vY2tlZDEPMA0GA1UECAwG
TW9ja2VkMQ8wDQYDVQQHDAZNb2NrZWQxDzANBgNVBAoMBk1vY2tlZDCCASIwDQYJ
KoZIhvcNAQEBBQADggEPADCCAQoCggEBALD76kHLps8vbcuVwlV5WUWnue/WKzqm
SYXFQSvaKm8R6Ofi56LAw8CLj788RDbyp8MdFVOqxZJF/scwINY24bzVGaN/PZ3b
uUcoMdS2EbJSCVzPzK7akUa/+w/6NlPBtQJJqxKzeHA5qV0qXpn0dfGLUCQryvf8
XNv0wM9qDMmc0odNEiNh0AafOAkdvtdl26V6xgvOZ5nr7LBPAruwKc/wUNlO3N1F
kV/6y6DF7B8GOvCQHILXDXfgqysnMqxrWHcrFH5DIbd1xXbPKUiDyX8j6QXoygi2
WDmSvb7NsiAwACqCjlTU79i/RmuCRiSrWvs/Gr/HIIHNJVGGy6WcjscCAwEAAaNQ
ME4wHQYDVR0OBBYEFGumFmmfzuM18BJIBVJxVapQzNtcMB8GA1UdIwQYMBaAFGum
FmmfzuM18BJIBVJxVapQzNtcMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQELBQAD
ggEBAFR4rpOQmu1kaT8ScCSozVQogna3b2YV5bvGGrU3TUkJOE3cBEkEAonPA9uE
YY8ddQqSdJdKr7oJNBWYo2ywltG/p/dvRu4FGLICD9OgTcN3220fcetr5JU7Ngi7
kNAFpJZxy6aDhRpJJvIkN81rrxXgqFE0u7VoFJfrHUoyoNXgzibh7em1zau4g/zk
vseUcPm90XJePTDQTzlWPMyuRSa2BE2q8nWcpX7suJstBHAGTmKs/iErmjJ4A4k+
1ZNBxGrrFbGi7rcO9pPD4iBnhL6jjOYDuRgjCn9evJI4Np9E4QP+x4xbpoZILLbl
EncRJJEl8/Q4jzWq+/JOGQVf6ds=
-----END CERTIFICATE-----
`}
				delete(v, "certificateFile")
			} else {
				replaceKeyFileAndCertificateFile(value)
			}
		}
	case []interface{}:
		for _, item := range v {
			replaceKeyFileAndCertificateFile(item)
		}
	}
}

func removeExtDatRules(obj interface{}) {
	switch v := obj.(type) {
	case map[string]interface{}:
		fmt.Printf("Processing map with keys: %v\n", getKeys(v))
		for key, value := range v {
			switch val := value.(type) {
			case []interface{}:
				fmt.Printf("Processing array in key '%s' with %d elements\n", key, len(val))
				filteredRules := make([]interface{}, 0)
				for _, rule := range val {
					switch r := rule.(type) {
					case string:
						if strings.Contains(r, "ext:") && strings.Contains(r, ".dat") {
							fmt.Printf("Removing string rule: %s\n", r)
							continue
						}
						filteredRules = append(filteredRules, r)
					case map[string]interface{}:
						if domains, ok := r["domain"].([]interface{}); ok {
							fmt.Printf("Found nested domain rules with %d elements\n", len(domains))
							filteredDomains := make([]interface{}, 0)
							for _, domain := range domains {
								if domainStr, ok := domain.(string); ok {
									if strings.Contains(domainStr, "ext:") && strings.Contains(domainStr, ".dat") {
										fmt.Printf("Removing nested rule: %s\n", domainStr)
										continue
									}
									filteredDomains = append(filteredDomains, domain)
								}
							}
							r["domain"] = filteredDomains
						}
						if ips, ok := r["ip"].([]interface{}); ok {
							fmt.Printf("Found nested ip rules with %d elements\n", len(ips))
							filteredIps := make([]interface{}, 0)
							for _, ip := range ips {
								if ipStr, ok := ip.(string); ok {
									if strings.Contains(ipStr, "ext:") && strings.Contains(ipStr, ".dat") {
										fmt.Printf("Removing nested ip rule: %s\n", ipStr)
										continue
									}
									filteredIps = append(filteredIps, ip)
								}
							}
							r["ip"] = filteredIps
						}
						filteredRules = append(filteredRules, r)
					default:
						filteredRules = append(filteredRules, rule)
					}
				}
				v[key] = filteredRules
			case map[string]interface{}:
				removeExtDatRules(val)
			}
		}
	case []interface{}:
		for _, item := range v {
			removeExtDatRules(item)
		}
	}
}

func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
