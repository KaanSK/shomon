// This package is derived from https://github.com/ns3777k/go-shodan

package shodan

import (
	"encoding/json"
	"math/big"
	"net"
	"strconv"
)

// Facet is a property to get summary information on.
type Facet struct {
	Count int    `json:"count"`
	Value string `json:"value"`
}

// IntString is string with custom unmarshaling.
type IntString string

// UnmarshalJSON handles either a string or a number
// and casts it to string.
func (v *IntString) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		*v = IntString(s)
		return nil
	}

	var n int
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}

	*v = IntString(strconv.Itoa(n))

	return nil
}

// String method just returns string out of IntString.
func (v *IntString) String() string {
	return string(*v)
}

// HostServicesOptions is options for querying services.
type HostServicesOptions struct {
	History bool `url:"history,omitempty"`
	Minify  bool `url:"minify,omitempty"`
}

// HostLocation is the location of the host.
type HostLocation struct {
	City         string  `json:"city"`
	RegionCode   string  `json:"region_code"`
	AreaCode     int     `json:"area_code"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Country      string  `json:"country_name"`
	CountryCode  string  `json:"country_code"`
	CountryCode3 string  `json:"country_code3"`
	Postal       string  `json:"postal_code"`
	DMA          int     `json:"dma_code"`
}

// HostDHParams is the Diffie-Hellman parameters if available.
type HostDHParams struct {
	Prime       string     `json:"prime"`
	PublicKey   string     `json:"public_key"`
	Bits        int        `json:"bits"`
	Generator   *IntString `json:"generator"`
	Fingerprint string     `json:"fingerprint"`
}

// HostTLSExtEntry contains id and name.
type HostTLSExtEntry struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// HostCipher is a cipher description.
type HostCipher struct {
	Version string `json:"version"`
	Bits    int    `json:"bits"`
	Name    string `json:"name"`
}

// HostCertificatePublicKey holds type and bits length of the key.
type HostCertificatePublicKey struct {
	Type string `json:"type"`
	Bits int    `json:"bits"`
}

// HostCertificateAttributes is an ordinary certificate attributes description.
type HostCertificateAttributes struct {
	CountryName         string `json:"C,omitempty"`
	CommonName          string `json:"CN,omitempty"`
	Locality            string `json:"L,omitempty"`
	Organization        string `json:"O,omitempty"`
	StateOrProvinceName string `json:"ST,omitempty"`
	OrganizationalUnit  string `json:"OU,omitempty"`
}

// HostCertificateExtension represent single cert extension.
type HostCertificateExtension struct {
	Data       string `json:"data"`
	Name       string `json:"name"`
	IsCritical bool   `json:"critical,omitempty"`
}

// HostCertificate contains common certificate description.
type HostCertificate struct {
	SignatureAlgorithm string                      `json:"sig_alg"`
	IsExpired          bool                        `json:"expired"`
	Version            int                         `json:"version"`
	Serial             *big.Int                    `json:"serial"`
	Issued             string                      `json:"issued"`
	Expires            string                      `json:"expires"`
	Fingerprint        map[string]string           `json:"fingerprint"`
	Issuer             *HostCertificateAttributes  `json:"issuer"`
	Subject            *HostCertificateAttributes  `json:"subject"`
	PublicKey          *HostCertificatePublicKey   `json:"pubkey"`
	Extensions         []*HostCertificateExtension `json:"extensions"`
}

// HostSSL holds ssl host information.
type HostSSL struct {
	Versions    []string           `json:"versions"`
	Chain       []string           `json:"chain"`
	DHParams    *HostDHParams      `json:"dhparams"`
	TLSExt      []*HostTLSExtEntry `json:"tlsext"`
	Cipher      *HostCipher        `json:"cipher"`
	Certificate *HostCertificate   `json:"cert"`
}

// HostData is all services that have been found on the given host IP.
type HostData struct {
	Product      string                 `json:"product"`
	Hostnames    []string               `json:"hostnames"`
	Version      IntString              `json:"version"`
	Title        string                 `json:"title"`
	SSL          *HostSSL               `json:"ssl"`
	IP           net.IP                 `json:"ip_str"`
	OS           string                 `json:"os"`
	Organization string                 `json:"org"`
	ISP          string                 `json:"isp"`
	CPE          []string               `json:"cpe"`
	Data         string                 `json:"data"`
	ASN          string                 `json:"asn"`
	Port         int                    `json:"port"`
	HTML         string                 `json:"html"`
	Banner       string                 `json:"banner"`
	Link         string                 `json:"link"`
	Transport    string                 `json:"transport"`
	Domains      []string               `json:"domains"`
	Timestamp    string                 `json:"timestamp"`
	DeviceType   string                 `json:"devicetype"`
	Location     *HostLocation          `json:"location"`
	ShodanData   map[string]interface{} `json:"_shodan"`
	Opts         map[string]interface{} `json:"opts"`
}

// Host is the all information about the host.
type Host struct {
	OS              string      `json:"os"`
	Ports           []int       `json:"ports"`
	IP              net.IP      `json:"ip_str"`
	ISP             string      `json:"isp"`
	Hostnames       []string    `json:"hostnames"`
	Organization    string      `json:"org"`
	Vulnerabilities []string    `json:"vulns"`
	ASN             string      `json:"asn"`
	LastUpdate      string      `json:"last_update"`
	Data            []*HostData `json:"data"`
	HostLocation
}

// HostQueryOptions is Shodan search query options.
type HostQueryOptions struct {
	Query  string `url:"query"`
	Facets string `url:"facets,omitempty"`
	Minify bool   `url:"minify,omitempty"`
	Page   int    `url:"page,omitempty"`
}

// HostMatch is the search results with all matched hosts.
type HostMatch struct {
	Total   int                 `json:"total"`
	Facets  map[string][]*Facet `json:"facets"`
	Matches []*HostData         `json:"matches"`
}

// HostQueryTokens is filters are being used by the query string and what
// parameters were provided to the filters.
type HostQueryTokens struct {
	Filters    []string               `json:"filters"`
	String     string                 `json:"string"`
	Errors     []string               `json:"errors"`
	Attributes map[string]interface{} `json:"attributes"`
}
