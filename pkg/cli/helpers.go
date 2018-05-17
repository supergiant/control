package cli

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/imdario/mergo"
	"github.com/urfave/cli"

	"github.com/supergiant/supergiant/pkg/client"
)

type GlobalConfig struct {
	Server   string `json:"server"`
	Token    string `json:"token"`
	CertFile string `json:"cert_file"`
}

func newGlobalConf(c *cli.Context) GlobalConfig {
	return GlobalConfig{
		c.String("server"),
		c.String("api-token"),
		c.String("cert-file"),
	}
}

func Client(c *cli.Context) *client.Client {
	globalConf := GlobalConfig{}

	if b, err := ioutil.ReadFile(globalConfFile); err == nil {
		// NOTE no error handling here, silently continues
		json.Unmarshal(b, &globalConf)
	}

	conf := newGlobalConf(c)

	if err := mergo.Merge(&conf, globalConf); err != nil {
		panic(err)
	}

	return client.New(conf.Server, "token", conf.Token, conf.CertFile)
}

//------------------------------------------------------------------------------

func printObj(obj interface{}) error {
	out, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func printList(c *cli.Context, list interface{}) error {
	// Optional formatting
	if format := c.String("format"); format != "" {
		tmpl, err := template.New("format").Parse(format)
		if err != nil {
			return err
		}
		items := reflect.ValueOf(list).Elem().FieldByName("Items")
		for i := 0; i < items.Len(); i++ {
			if err := tmpl.Execute(os.Stdout, items.Index(i).Interface()); err != nil {
				return err
			}
		}
		return nil
	}
	return printObj(list)
}

func listFilters(c *cli.Context) (map[string][]string, error) {
	filters := make(map[string][]string)
	for _, filter := range c.StringSlice("filter") {
		segments := strings.Split(filter, ":")
		if len(segments) != 2 {
			return nil, fmt.Errorf("Invalid filter flag '%s'", filter)
		}
		field := segments[0]
		values := strings.Split(segments[1], ",")
		filters[field] = values
	}
	return filters, nil
}
