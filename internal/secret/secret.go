package secret

import (
	"os"
	"reflect"
	"strings"
)

func ReplaceSecret(input string) string {
	if strings.HasPrefix(input, "$") {
		return os.Getenv(input[1:])
	}
	return input
}

func replaceReflected(reflectedValue reflect.Value) reflect.Value {
	switch reflectedValue.Kind() {
	case reflect.String:
		return reflect.ValueOf(ReplaceSecret(reflectedValue.String()))
	case reflect.Slice:
		for i := 0; i < reflectedValue.Len(); i++ {
			reflectedValue.Index(i).Set(replaceReflected(reflectedValue.Index(i)))
		}
	case reflect.Map:
		keys := reflectedValue.MapKeys()
		for _, key := range keys {
			reflectedValue.SetMapIndex(key, replaceReflected(reflectedValue.MapIndex(key)))
		}
	case reflect.Struct:
		for i := 0; i < reflectedValue.NumField(); i++ {
			field := reflectedValue.Field(i)
			field.Set(replaceReflected(field))
		}
	}
	return reflectedValue
}

func ReplaceSecrets[T any](input T) {
	reflectedInput := reflect.ValueOf(input)
	replaceReflected(reflectedInput.Elem())
}
