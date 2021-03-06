package pop

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/markbates/going/defaults"
	"gopkg.in/yaml.v2"
)

var lookupPaths = []string{"", "/config", "../", "../config", "../..", "../../config"}
var ConfigName = "database.yml"

func init() {
	LoadConfig()
}

func LoadConfig() {
	path, err := findConfigPath()
	if err == nil {
		loadConfig(path)
	}
}

func LookupPaths() []string {
	return lookupPaths
}

func AddLookupPaths(paths ...string) {
	lookupPaths = append(paths, lookupPaths...)
	LoadConfig()
}

type ConnectionDetails struct {
	Dialect  string
	Database string
	Host     string
	Port     string
	User     string
	Password string
	URL      string
}

func (cd *ConnectionDetails) Parse(port string) error {
	if cd.URL != "" {
		u, err := url.Parse(cd.URL)
		if err != nil {
			return err
		}
		cd.Database = strings.TrimPrefix(u.Path, "/")
		hp := strings.Split(u.Host, ":")
		cd.Host = hp[0]
		cd.Port = defaults.String(hp[1], port)
		cd.User = u.User.Username()
		cd.Password, _ = u.User.Password()
	}
	return nil
}

func findConfigPath() (string, error) {
	// pwd, err := getAppPath()
	// if err != nil {
	// 	return "", err
	// }
	//
	// lookup paths
	for _, p := range LookupPaths() {
		path, _ := filepath.Abs(filepath.Join(p, ConfigName))
		if _, err := os.Stat(path); err == nil {
			return path, err
		}
	}
	return "", errors.New("[POP]: Tried to load configuration file, but couldn't find it.")
}

func getAppPath() (string, error) {
	pwd := os.Getenv("APP_PATH")
	if pwd == "" {
		b, err := exec.Command("pwd").Output()
		if err != nil {
			return "", err
		}
		pwd = string(b)
	}
	return strings.TrimSuffix(pwd, "\n"), nil
}

func loadConfig(path string) error {
	if Debug {
		fmt.Printf("[POP]: Loading config file from %s\n", path)
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	tmpl := template.New("test")
	tmpl.Funcs(map[string]interface{}{
		"envOr": func(s1, s2 string) string {
			return defaults.String(os.Getenv(s1), s2)
		},
		"env": func(s1 string) string {
			return os.Getenv(s1)
		},
	})
	t, err := tmpl.Parse(string(b))
	if err != nil {
		return err
	}

	var bb bytes.Buffer
	err = t.Execute(&bb, nil)
	if err != nil {
		return err
	}

	conns := map[string]*ConnectionDetails{}
	err = yaml.Unmarshal(bb.Bytes(), &conns)
	if err != nil {
		return err
	}
	for n, c := range conns {
		con := NewConnection(c)
		Connections[n] = con
	}
	return nil
}
