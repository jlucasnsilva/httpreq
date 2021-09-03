package httpreq

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

type (
	decodeMacros struct {
		Macros map[string]string `yaml:"where,omitempty"`
	}
)

var (
	appFS = afero.NewOsFs()

	example = Group{
		Config: Config{
			HTTPS: false,
			Host:  "localhost",
			Port:  "8080",
			Header: map[string]string{
				"Authorization": "Bearer my.awesome.token",
			},
		},
		EndpointCalls: []EndpointCall{
			{
				ID:      "login",
				Call:    "POST /api/v1/login",
				Payload: `{ "email": "{{ .email }}", "password": "{{ .password }}" }`,
				Header: map[string]string{
					"X-Parameter": "magic",
				},
			},
		},
		Macros: map[string]string{
			"email":    "mr.frodo@lotr.com",
			"password": "$PASSWORD",
		},
	}

	httpMethods = map[string]struct{}{
		"POST":   {},
		"GET":    {},
		"PUT":    {},
		"DELETE": {},
		"HEAD":   {},
	}
)

// CreateFile ...
func CreateFile(filename string) error {
	return createFile(appFS, filename)
}

func createFile(afs afero.Fs, filename string) error {
	b, err := yaml.Marshal(&example)
	if err != nil {
		return fmt.Errorf("create file: marshaling error: %v", err)
	}
	err = afero.WriteFile(afs, filename, b, 0666)
	if err != nil {
		return fmt.Errorf("create file: writing error: %v", err)
	}
	return nil
}

// Parse ...
func Parse(filename string) (*Group, error) {
	return parse(appFS, filename)
}

func parse(afs afero.Fs, filename string) (*Group, error) {
	data, err := afero.ReadFile(afs, filename)
	if err != nil {
		return nil, fmt.Errorf("parse: reading error: %v", err)
	}

	macros, err := parseMacros(data)
	if err != nil {
		return nil, fmt.Errorf("parse: macros error: %v", err)
	}

	tmpl, err := template.New("replace").Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("parse: macros substitution error: %v", err)
	}

	bs := bytes.Buffer{}
	if err := tmpl.Execute(&bs, macros); err != nil {
		return nil, fmt.Errorf("parse: macros substitution error: %v", err)
	}

	g := new(Group)
	if err := yaml.Unmarshal(bs.Bytes(), g); err != nil {
		return nil, fmt.Errorf("parse: group unmarshalling error: %v", err)
	}
	g.Macros = macros
	return g, nil
}

func parseMacros(in []byte) (map[string]string, error) {
	ms := Group{}
	if err := yaml.Unmarshal(in, &ms); err != nil {
		return nil, fmt.Errorf("macros: unmarshalling error: %v", err)
	}

	replaceMacrosEnvVariables(ms.Macros)
	return ms.Macros, nil
}

func replaceMacrosEnvVariables(ms map[string]string) {
	for macro, value := range ms {
		v := strings.TrimSpace(value)
		if strings.HasPrefix(v, "$") {
			s := strings.TrimLeft(v, "$")
			ev, ok := os.LookupEnv(s)
			if !ok {
				panic("environment variable not defined: " + s)
			}
			v = ev
		}
		ms[macro] = v
	}
}
