package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func Load(cfg any, envFiles ...string) error {
	if len(envFiles) == 0 {
		envFiles = []string{".env"}
	}

	for _, f := range envFiles {
		if err := godotenv.Load(f); err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("config - Load - godotenv.Load(%s): %w", f, err)
			}
		}
	}

	if err := env.Parse(cfg); err != nil {
		return fmt.Errorf("config - Load - env.Parse: %w", err)
	}

	trimQuoteFields(cfg)
	return nil
}

func trimQuoteFields(cfg any) {
	v := reflect.ValueOf(cfg)

	if v.Kind() != reflect.Pointer || v.IsNil() {
		return
	}
	walkAndTrim(v.Elem())
}

func walkAndTrim(v reflect.Value) {
	switch v.Kind() {
	case reflect.Struct:
		for _, f := range v.Fields() {
			walkAndTrim(f)
		}
	case reflect.String:
		if v.CanSet() {
			v.SetString(strings.Trim(strings.TrimSpace(v.String()), `"'`))
		}
	}
}
