package config

import (
	"github.com/fatih/structs"
	"testing"
)

var configData = []struct {
	name       string
	data       interface{}
	expectedOk bool
	valid      bool
}{
	{"working_config", map[string]interface{}{
		"AppAddress": "127.0.0.1",
		"AppPort":    3000,
		"DBType":     DBPostgres,
		"LogFile":    "tttt",
	}, true, true},
	{"invalid_datatype_config", map[int]interface{}{
		2: 3,
	}, false, false},
}

func TestLoad(t *testing.T) {
	l := DefaultLoader{}
	for _, d := range configData {
		c, err := l.Load(d.data)
		if d.expectedOk && err != nil {
			t.Fatalf("%s: valid config expected. | error: %s", d.name, err)
		}
		s := structs.New(c)
		if d.valid {
			for _, f := range s.Fields() {
				dataMap := d.data.(map[string]interface{})

				_, present := dataMap[f.Name()]
				if !present {
					t.Fatalf("%s: key is present in data but not present in config | key: %s", d.name, f.Name())
				}
			}
		}
	}
}

func TestLoadWithoutParam(t *testing.T) {
	l := DefaultLoader{}
	_, err := l.Load()
	if err == nil {
		t.Fatal("error expected, nil returned")
	}
}
