package conf

// Options : argument struct
type Options struct {
	Endpoint     string `short:"e" long:"endpoint" description:"Full endpoint for alert API" default:"http://127.0.0.1:9000/api/alert"`
	ShodanKey    string `short:"s" long:"shodanKey" description:"Shodan Api Key" required:"true"`
	CaseTemplate string `short:"c" long:"caseTemplate" description:"Case template for alert creation. Can be empty" default:""`
	TheHiveKey   string `short:"t" long:"theHiveKey" description:"TheHive api key" required:"true"`
	Verbose      bool   `short:"v" long:"verbose" description:"Show verbose debug information"`
}

// Config : Used to retrieve conf key/values. Globally available
var Config Options
