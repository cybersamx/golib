package cli_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/cybersamx/golib/cli"
)

// targetVars should contain the fields of the data types that are supported by
// the function InitFlags.
type targetVars struct {
	strDefault    string // String with the default value.
	strReader     string // String to be overridden by a config file.
	strEnv        string // String to be overridden by an environment variable.
	strFlag       string // String to be overridden by a flag.
	nested        string // String bound to a nested flag.
	number        int
	boolean       bool
	singleKey     map[string]string
	multiKeys     map[string]string
	singleKeyJSON map[string]string
	multiKeysJSON map[string]string
	strList       []string
	duration      time.Duration
}

func newCommand(t *testing.T, target *targetVars, cfgFileReader io.Reader) *cobra.Command {
	t.Helper()

	bindings := []FlagBinding{
		{
			Usage:   "test string with default value",
			Name:    "str-default",
			Target:  &target.strDefault,
			Default: "str_default",
		},
		{
			Usage:   "test string to be overridden by reader",
			Name:    "str-reader",
			Target:  &target.strReader,
			Default: "str_reader_default",
		},
		{
			Usage:   "test string to be overridden by env variable",
			Name:    "str-env",
			Target:  &target.strEnv,
			Default: "str_env_default",
		},
		{
			Usage:   "test string to be overridden by flag",
			Name:    "str-flag",
			Target:  &target.strFlag,
			Default: "str_flag_default",
		},
		{
			Usage:   "test nested flag",
			Name:    "deep.nested.str",
			Target:  &target.nested,
			Default: "nested_flag_default",
		},
		{
			Usage:   "test int flag",
			Name:    "number",
			Target:  &target.number,
			Default: 123,
		},
		{
			Usage:   "test bool flag",
			Name:    "boolean",
			Target:  &target.boolean,
			Default: false,
		},
		{
			Usage:   "test map flag with a single key-value",
			Name:    "single-key",
			Target:  &target.singleKey,
			Default: map[string]string{"key": "default"},
		},
		{
			Usage:   "test map flag with multi key-values",
			Name:    "multi-keys",
			Target:  &target.multiKeys,
			Default: map[string]string{"map.key1": "map.default1", "map.key2": "map.default2"},
		},
		{
			Usage:   "test map flag with a single key-value (in json format)",
			Name:    "single-key-json",
			Target:  &target.singleKeyJSON,
			Default: map[string]string{"key": "defaultJSON"},
		},
		{
			Usage:   "test map flag with multi key-values (in json format)",
			Name:    "multi-keys-json",
			Target:  &target.multiKeysJSON,
			Default: map[string]string{"map.key1": "map.defaultJSON1", "map.key2": "map.defaultJSON2"},
		},
		{
			Usage:   "test string slice flag",
			Name:    "str-list",
			Target:  &target.strList,
			Default: []string{"list.default1", "list.default2", "list.default3"},
		},
		{
			Usage:   "test duration flag",
			Name:    "duration",
			Target:  &target.duration,
			Default: time.Hour + 2*time.Minute + 3*time.Second,
		},
	}

	v := NewViper("GL")
	cmd := cobra.Command{
		Use: "app_test",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	if cfgFileReader != nil {
		v.SetConfigType("yaml")
		err := v.ReadConfig(cfgFileReader)
		require.NoError(t, err)
	}

	err := InitFlags(v, cmd.Flags(), bindings)
	require.NoError(t, err)

	return &cmd
}

func TestInitFlags_Defaults(t *testing.T) {
	target := new(targetVars)

	cmd := newCommand(t, target, nil)
	cmd.SetArgs(nil)
	err := cmd.Execute()
	require.NoError(t, err)

	want := targetVars{
		strDefault:    "str_default",
		strReader:     "str_reader_default",
		strEnv:        "str_env_default",
		strFlag:       "str_flag_default",
		nested:        "nested_flag_default",
		number:        123,
		boolean:       false,
		singleKey:     map[string]string{"key": "default"},
		multiKeys:     map[string]string{"map.key1": "map.default1", "map.key2": "map.default2"},
		singleKeyJSON: map[string]string{"key": "defaultJSON"},
		multiKeysJSON: map[string]string{"map.key1": "map.defaultJSON1", "map.key2": "map.defaultJSON2"},
		strList:       []string{"list.default1", "list.default2", "list.default3"},
		duration:      time.Hour + 2*time.Minute + 3*time.Second,
	}

	diff := pretty.Compare(want, target)
	assert.Emptyf(t, diff, "want: %+v, got: %+v", want, target)
}

func TestInitFlags_ConfigFile(t *testing.T) {
	yamlConfig := []byte(`
str-reader: str_reader
deep:
  nested:
    str: nested_str_reader
number: 234
boolean: true
single-key:
  key: reader
single-key-json:
  key: readerJSON
multi-keys:
  map.key1: map.reader1
  map.key2: map.reader2
str-list:
  - list.reader1
  - list.reader2
  - list.reader3
duration: 2h3m4s
`)

	reader := bytes.NewBuffer(yamlConfig)

	target := new(targetVars)
	cmd := newCommand(t, target, reader)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	require.NoError(t, err)

	want := targetVars{
		strDefault:    "str_default",
		strReader:     "str_reader",
		strEnv:        "str_env_default",
		strFlag:       "str_flag_default",
		nested:        "nested_str_reader",
		number:        234,
		boolean:       true,
		singleKey:     map[string]string{"key": "reader"},
		multiKeys:     map[string]string{"map.key1": "map.reader1", "map.key2": "map.reader2"},
		singleKeyJSON: map[string]string{"key": "readerJSON"},
		multiKeysJSON: map[string]string{"map.key1": "map.defaultJSON1", "map.key2": "map.defaultJSON2"},
		strList:       []string{"list.reader1", "list.reader2", "list.reader3"},
		duration:      2*time.Hour + 3*time.Minute + 4*time.Second,
	}

	diff := pretty.Compare(want, target)
	assert.Emptyf(t, diff, "want: %+v, got: %+v", want, target)
}

func TestInitFlags_Env(t *testing.T) {
	// Temporarily set environment for the duration of the test.
	t.Setenv("GL_STR_ENV", "str_env")
	t.Setenv("GL_NUMBER", "345")
	t.Setenv("GL_BOOLEAN", "true")
	t.Setenv("GL_DEEP_NESTED_STR", "nested_str_env")
	t.Setenv("GL_SINGLE_KEY", "key=env")
	t.Setenv("GL_SINGLE_KEY_JSON", `{"key": "envJSON"}`)
	t.Setenv("GL_MULTI_KEYS", `{"map.key1": "map.env1", "map.key2": "map.env2"}`)
	t.Setenv("GL_MULTI_KEYS_JSON", `{"map.key1": "map.envJSON1", "map.key2": "map.envJSON2"}`)
	t.Setenv("GL_STR_LIST", `"list.env1","list.env2","list.env3"`)
	t.Setenv("GL_DURATION", "3h4m5s")

	target := new(targetVars)
	cmd := newCommand(t, target, nil)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	require.NoError(t, err)

	want := targetVars{
		strDefault:    "str_default",
		strReader:     "str_reader_default",
		strEnv:        "str_env",
		strFlag:       "str_flag_default",
		nested:        "nested_str_env",
		number:        345,
		boolean:       true,
		singleKey:     map[string]string{"key": "env"},
		multiKeys:     map[string]string{"map.key1": "map.env1", "map.key2": "map.env2"},
		singleKeyJSON: map[string]string{"key": "envJSON"},
		multiKeysJSON: map[string]string{"map.key1": "map.envJSON1", "map.key2": "map.envJSON2"},
		strList:       []string{"list.env1", "list.env2", "list.env3"},
		duration:      3*time.Hour + 4*time.Minute + 5*time.Second,
	}

	diff := pretty.Compare(want, target)
	assert.Emptyf(t, diff, "want: %+v, got: %+v", want, target)
}

func TestInitFlags_Flags(t *testing.T) {
	target := new(targetVars)
	cmd := newCommand(t, target, nil)
	cmd.SetArgs([]string{
		"--str-flag=str_flag",
		"--deep.nested.str=nested_str_flag",
		"--number=456",
		"--boolean=true",
		"--single-key",
		"key=flag",
		"--single-key-json",
		"key=flagJSON",
		"--multi-keys",
		"map.key1=map.flag1",
		"--multi-keys",
		"map.key2=map.flag2",
		"--str-list=list.flag1",
		"--str-list=list.flag2",
		"--str-list=list.flag3",
		"--duration=4h5m6s",
	})

	err := cmd.Execute()
	require.NoError(t, err)

	want := targetVars{
		strDefault:    "str_default",
		strReader:     "str_reader_default",
		strEnv:        "str_env_default",
		strFlag:       "str_flag",
		nested:        "nested_str_flag",
		number:        456,
		boolean:       true,
		singleKey:     map[string]string{"key": "flag"},
		multiKeys:     map[string]string{"map.key1": "map.flag1", "map.key2": "map.flag2"},
		singleKeyJSON: map[string]string{"key": "flagJSON"},
		multiKeysJSON: map[string]string{"map.key1": "map.defaultJSON1", "map.key2": "map.defaultJSON2"},
		strList:       []string{"list.flag1", "list.flag2", "list.flag3"},
		duration:      4*time.Hour + 5*time.Minute + 6*time.Second,
	}

	diff := pretty.Compare(want, target)
	assert.Emptyf(t, diff, "want: %+v, got: %+v", want, target)
}

func TestInitFlags_Map(t *testing.T) {
	type tuple struct {
		key string
		val string
	}

	tests := []struct {
		description string
		args        []string
		env         *tuple
		want        map[string]string
		wantErr     bool
	}{
		{
			description: "Using flag with key=value format",
			args:        []string{"--arg", "key1=val1", "--arg", "key2=val2"},
			env:         nil,
			want:        map[string]string{"key1": "val1", "key2": "val2"},
			wantErr:     false,
		},
		{
			description: `Using flag with json format`,
			args:        []string{"--arg", `{"key1": "val1", "key2": "val2"}`},
			env:         nil,
			want:        map[string]string{"key1": "val1", "key2": "val2"},
			wantErr:     true,
		},
		{
			description: `Using env variable with json format`,
			args:        nil,
			env:         &tuple{key: "GL_ARG", val: `{"key1": "val1", "key2": "val2"}`},
			want:        map[string]string{"key1": "val1", "key2": "val2"},
			wantErr:     false,
		},
	}

	target := map[string]string{}

	bindings := []FlagBinding{
		{
			Usage:   "test map[string]string flag",
			Name:    "arg",
			Target:  &target,
			Default: map[string]string{},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			delete(target, "arg")

			if test.env != nil {
				t.Setenv(test.env.key, test.env.val)
			}

			v := NewViper("GL")
			cmd := cobra.Command{
				Run: func(cmd *cobra.Command, args []string) {},
			}
			err := InitFlags(v, cmd.Flags(), bindings)
			require.NoError(t, err)

			cmd.SetArgs(test.args)
			err = cmd.Execute()
			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			diff := pretty.Compare(test.want, target)
			assert.Emptyf(t, diff, "want: %+v, got: %+v", test.want, target)
		})
	}
}
