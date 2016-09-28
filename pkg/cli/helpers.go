package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/imdario/mergo"
	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/urfave/cli"
)

func commandList(title string, fn func(*cli.Context) error) cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "list " + title + "s",
		Flags: append(baseFlags, []cli.Flag{
			cli.StringSliceFlag{
				Name:  "filter",
				Usage: "--filter=name:this,or,that --filter=other_field:value",
			},
			cli.StringFlag{
				Name:  "format",
				Usage: "--format=\"{{ .ThisField }}\"",
			},
		}...),
		Action: fn,
	}
}

func commandCreate(title string, fn func(*cli.Context) error) cli.Command {
	return cli.Command{
		Name:  "create",
		Usage: "create a new " + title,
		Flags: append(baseFlags, []cli.Flag{
			cli.StringFlag{
				Name:  "file, f",
				Usage: "JSON input file",
			},
		}...),
		Action: fn,
	}
}

func commandGet(title string, fn func(*cli.Context) error) cli.Command {
	return cli.Command{
		Name:  "get",
		Usage: "get a " + title,
		Flags: append(baseFlags, []cli.Flag{
			cli.StringFlag{
				Name:  "id",
				Usage: "the resource ID",
			},
		}...),
		Action: fn,
	}
}

func commandUpdate(title string, fn func(*cli.Context) error) cli.Command {
	return cli.Command{
		Name:  "update",
		Usage: "update a " + title,
		Flags: append(baseFlags, []cli.Flag{
			cli.StringFlag{
				Name:  "id",
				Usage: "the resource ID",
			},
			cli.StringFlag{
				Name:  "file, f",
				Usage: "JSON input file",
			},
		}...),
		Action: fn,
	}
}

func commandAction(action string, title string, fn func(*cli.Context) error) cli.Command {
	return cli.Command{
		Name:  action,
		Usage: action + " a " + title,
		Flags: append(baseFlags, []cli.Flag{
			cli.StringFlag{
				Name:  "id",
				Usage: "the resource ID",
			},
		}...),
		Action: fn,
	}
}

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

type Client struct {
	*client.Client
	Config GlobalConfig
}

func newClient(c *cli.Context) *Client {
	globalConf := GlobalConfig{}

	if b, err := ioutil.ReadFile(globalConfFile); err == nil {
		// NOTE no error handling here, silently continues
		json.Unmarshal(b, &globalConf)
	}

	conf := newGlobalConf(c)

	if err := mergo.Merge(&conf, globalConf); err != nil {
		panic(err)
	}

	return &Client{client.New(conf.Server, "token", conf.Token, conf.CertFile), conf}
}

func decodeInputFileInto(c *cli.Context, item model.Model) (err error) {
	var file *os.File

	switch filepath := c.String("f"); filepath {
	case "":
		return errors.New("-f required")
	case "-":
		file = os.Stdin
	default:
		file, err = os.Open(filepath)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	err = json.NewDecoder(file).Decode(item)
	return err
}

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
