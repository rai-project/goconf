// Read configuration from files which may be overriden by envar value
// with prefix. To override nested yaml value with envar, replace
// its path with an underscore.
package goconf

import (
	"github.com/spf13/viper"
	"path/filepath"
	"reflect"
	"strings"
)

const tagName = "config"

var config *viper.Viper

func buildConfiguration(reflection reflect.Value, configValue interface{}) {
	if reflection.IsValid() && reflection.Kind() == reflect.Struct {
		//fmt.Println("Instance can be set :", reflection.CanSet())
		for i := 0; i < reflection.NumField(); i++ {
			structField := reflection.Type().Field(i)
			valueField := reflect.Indirect(reflection.Field(i))
			//fmt.Printf("`%s` can be set : %t\n", structField.Name, reflection.Field(i).CanSet())
			tag := structField.Tag
			configTag := tag.Get("config")
			if configTag == "" {
				continue
			}
			var conf interface{}
			if configValue == nil {
				conf = config.Get(configTag)
			} else {
				conf = configValue
			}
			switch t := conf.(type) {
			default:
				confValue := reflect.ValueOf(conf)
				//fmt.Printf("valueField can be set : %t\n", valueField.CanSet())
				if valueField.CanSet() && confValue.IsValid() {
					valueField.Set(confValue)
				}
			case map[interface{}]interface{}:
				valueField = reflection.Field(i)
				for k, v := range t {
					if k.(string) == configTag {
						confValue := reflect.ValueOf(v)
						//fmt.Printf("valueField can be set: %t\n", valueField.CanSet())
						if valueField.CanSet() && confValue.IsValid() {
							valueField.Set(confValue)
						}
					}
				}
			case []interface{}:
				slice := reflect.MakeSlice(structField.Type, len(t), len(t))
				slices := reflect.New(slice.Type())
				slices.Elem().Set(slice)
				for i, v := range t {
					child := reflect.New(slices.Elem().Index(0).Type())
					buildConfiguration(child.Elem(), v)
					slices.Elem().Index(i).Set(child.Elem())
				}
				//fmt.Printf("valueField can be set: %t\n", valueField.CanSet())
				if valueField.CanSet() {
					valueField.Set(slices.Elem())
				}
			}
		}
	}
}

func GetConfig(class interface{}) {
	err := config.ReadInConfig()
	if err != nil {
		panic("Cannot find configuration file (config.yaml) " + err.Error())
	}
	r := reflect.ValueOf(class)
	conf := reflect.Indirect(r)
	buildConfiguration(conf, nil)
}

func Setup(fileType string, prefix string) {
	config = viper.New()
	config.SetConfigName("config")
	config.AddConfigPath(".")
	config.AddConfigPath("conf")
	config.AddConfigPath("config")
	config.SetConfigType(fileType)
	config.AutomaticEnv()
	config.SetEnvPrefix(prefix)
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func FromFile(filename string, prefix string) {
	Setup(strings.Replace(filepath.Ext(filename), ".", "", 1), prefix)
	config.SetConfigFile(filename)
}
