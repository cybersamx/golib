package cli

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/cybersamx/golib/stringsutils"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	envDelimiter = "_"
)

var ErrFlagBinding = errors.New("failed to bind flag")

func flagBindingError(flagName string, err error) error {
	return fmt.Errorf("flagBindingError - flag=%s; root_err=%v; %w",
		flagName, err, ErrFlagBinding)
}

type FlagBindingParser func(v *viper.Viper, flags *pflag.FlagSet, binding *FlagBinding) error

// FlagBinding represents the info for simplify the setup of the popular command line management
// packages such as spf13/pflag and spf13/viper.
//  1. Set up the name (posix flag aka pflag), shorthand (arg), and usage (description) of a flag.
//  2. Define the default value of a flag.
//  3. Once a flag is set by the user, where to bind the value of a flag to a target (variable or a
//     field in a struct object).
type FlagBinding struct {
	Usage     string
	Name      string
	Shorthand rune // One character for a shorthand.
	Target    any
	Default   any
	Parser    FlagBindingParser
}

func NewViper(envPrefix string) *viper.Viper {
	v := viper.New()
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()

	// Substitute STORE.DB-URL to STORE_DB_URL
	envReplacer := strings.NewReplacer(".", envDelimiter, "-", envDelimiter)
	v.SetEnvKeyReplacer(envReplacer)

	return v
}

// InitFlags sets up the accepted flags in a command line program and then bind the values set in the
// flags by the user to the target variable or field in a struct object.
func InitFlags(v *viper.Viper, flags *pflag.FlagSet, bindings []FlagBinding) error {
	for _, binding := range bindings {
		if binding.Parser != nil {
			if err := binding.Parser(v, flags, &binding); err != nil {
				return flagBindingError(binding.Name, err)
			}

			continue
		}

		val := v.Get(binding.Name)
		shorthand := stringsutils.RuneToString(binding.Shorthand)

		switch target := binding.Target.(type) {
		case *string:
			var def string
			if binding.Default != nil {
				// No ok checking needed, if it fails, the zero value will be assigned.
				def = binding.Default.(string)
			}

			flags.StringVarP(target, binding.Name, shorthand, def, binding.Usage)
			if err := v.BindPFlag(binding.Name, flags.Lookup(binding.Name)); err != nil {
				return flagBindingError(binding.Name, err)
			}

			if val != nil {
				if str, err := cast.ToStringE(val); err == nil {
					*target = str
				}
			}
		case *bool:
			var def bool
			if binding.Default != nil {
				def = binding.Default.(bool)
			}

			flags.BoolVarP(target, binding.Name, shorthand, def, binding.Usage)
			if err := v.BindPFlag(binding.Name, flags.Lookup(binding.Name)); err != nil {
				return flagBindingError(binding.Name, err)
			}

			if val != nil {
				if b, err := cast.ToBoolE(val); err == nil {
					*target = b
				}
			}
		case *int:
			var def int
			if binding.Default != nil {
				def = binding.Default.(int)
			}

			flags.IntVarP(target, binding.Name, shorthand, def, binding.Usage)
			if err := v.BindPFlag(binding.Name, flags.Lookup(binding.Name)); err != nil {
				return flagBindingError(binding.Name, err)
			}

			if val != nil {
				if i, err := cast.ToIntE(val); err == nil {
					*target = i
				}
			}
		case *time.Duration:
			var def time.Duration
			if binding.Default != nil {
				def = binding.Default.(time.Duration)
			}
			flags.DurationVarP(target, binding.Name, shorthand, def, binding.Usage)
			if err := v.BindPFlag(binding.Name, flags.Lookup(binding.Name)); err != nil {
				return flagBindingError(binding.Name, err)
			}

			if val != nil {
				if t, err := cast.ToDurationE(val); err == nil {
					*target = t
				}
			}
		case *[]string:
			var def []string
			if binding.Default != nil {
				def = binding.Default.([]string)
			}

			flags.StringSliceVarP(target, binding.Name, shorthand, def, binding.Usage)
			if err := v.BindPFlag(binding.Name, flags.Lookup(binding.Name)); err != nil {
				return flagBindingError(binding.Name, err)
			}

			// An environment variable, unlike a flag, is a singleton. Any subsequent set env will just override
			// the previous value. So the set string value should be encoded and the following section decodes
			// an env value. For env, viper returns a string (not []string) so we can decode the string.
			// For 1 item slice, the value can be encoded in key=value.
			// For 1+ item slice, the value must be encoded in csv item1,item2,item3 format.
			// We may have comma in the value as long as it is enclosed by "" - standard csv.
			// For a map[string]string typed flag, the flag value must be formatted as key=value only.
			if val != nil {
				s, err := cast.ToStringSliceE(val)
				if err != nil || len(s) == 0 {
					break
				}

				if len(s) == 1 {
					// If the env variable is set, we get an 1-item slice. Treat the string
					// as a csv and parse accordingly.
					sreader := strings.NewReader(s[0])
					creader := csv.NewReader(sreader)
					if valSlice, err := creader.Read(); err == nil {
						s = valSlice
					}
				}

				*target = s
			}
		case *map[string]string:
			var def map[string]string
			if binding.Default != nil {
				def = binding.Default.(map[string]string)
			}

			flags.StringToStringVarP(target, binding.Name, shorthand, def, binding.Usage)
			if err := v.BindPFlag(binding.Name, flags.Lookup(binding.Name)); err != nil {
				return flagBindingError(binding.Name, err)
			}

			if val != nil {
				m, err := cast.ToStringMapStringE(val)
				if err == nil {
					*target = m
					break
				}

				// An environment variable, unlike a flag, is a singleton. Any subsequent set env will just override
				// the previous value. So the set string value should be encoded and the following section decodes
				// an env value. For env, viper returns a string (not map[string]string) so we can decode the string.
				// For 1 key map, the value can be encoded in key=value or json {"key": "value"} format.
				// For 1+ key map, the value must be encoded in json {"key": "value"} format.
				str, err := cast.ToStringE(val)
				if err != nil {
					break
				}

				segments := strings.Split(str, "=")
				if len(segments) >= 2 {
					var valMap map[string]string
					if err := json.Unmarshal([]byte(segments[1]), &valMap); err == nil {
						*target = valMap
						break
					}

					valMap = map[string]string{segments[0]: segments[1]}
					*target = valMap
				}
			}
		default:
			if target == nil {
				break
			}

			typ := reflect.TypeOf(target).String()
			panic(fmt.Sprintf("need to implement logic to bind flag to type %s", typ))
		}
	}

	return nil
}
