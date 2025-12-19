package assert

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func DeepEqual[T any](t *testing.T, want, got T, errMsg ...string) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		if len(errMsg) > 0 {
			t.Fatal(errMsg[0])
		}
		t.Fatalf("want %v got %v", want, got)
	}
}

func Equal[T comparable](t *testing.T, want, got T, errMsg ...string) {
	t.Helper()
	if want != got {
		if len(errMsg) > 0 {
			t.Fatal(errMsg[0])
		}
		t.Fatalf("want %v got %v", want, got)
	}
}

func NotEqual[T comparable](t *testing.T, notWant, got T, errMsg ...string) {
	t.Helper()
	if notWant == got {
		if len(errMsg) > 0 {
			t.Fatal(errMsg[0])
		}
		t.Fatalf("want value to not be equal, got them equal %v", got)
	}
}

func Nil(t *testing.T, v any, errMsg ...string) {
	t.Helper()
	if !isNil(v) {
		if len(errMsg) > 0 {
			t.Fatal(errMsg[0])
		}

		t.Fatalf("want nil, got %#v", v)
	}
}

func NotNil(t *testing.T, v any, errMsg ...string) {
	t.Helper()
	if isNil(v) {
		if len(errMsg) > 0 {
			t.Fatal(errMsg[0])
		}
		t.Fatalf("want not nil, got nil")
	}
}

func True(t *testing.T, cond bool, errMsg ...string) {
	t.Helper()
	if !cond {
		if len(errMsg) > 0 {
			t.Fatal(errMsg[0])
		}
		t.Fatalf("want true, got false")
	}
}

func False(t *testing.T, cond bool, errMsg ...string) {
	t.Helper()
	if cond {
		if len(errMsg) > 0 {
			t.Fatal(errMsg[0])
		}
		t.Fatalf("want false, got true")
	}
}

func isNil(v any) bool {
	if v == nil {
		return true
	}

	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan:
		return val.IsNil()
	}

	return false
}

func Contains[T any, C comparable](t *testing.T, needle C, stash T, errMsg ...string) {
	t.Helper()
	containsValue, err := contains(needle, stash)
	if !containsValue {
		if len(errMsg) > 0 {
			t.Fatal(errMsg[0])
		}
		t.Fatal(err.Error())
	}
}

func NotContain[T any, C comparable](t *testing.T, needle C, stash T, errMsg ...string) {
	t.Helper()
	conatinsValue, err := contains(needle, stash)
	if conatinsValue {
		if len(errMsg) > 0 {
			t.Fatal(errMsg[0])
		}
		t.Fatal(err.Error())
	}
}

func contains[T any, C comparable](needle C, stash T) (bool, error) {
	stashValue := reflect.ValueOf(stash)

	switch stashValue.Kind() {
	case reflect.String:
		if stashString, ok := any(stash).(string); ok {
			needleString, ok2 := any(needle).(string)
			if !ok2 {
				return false, fmt.Errorf("Contains failed: needle must be string when stash is string")
			}
			if !strings.Contains(stashString, needleString) {
				return false, fmt.Errorf("Contains failed: '%s' does not contain '%s'", stashString, needleString)
			}

			return true, fmt.Errorf("Contains failed: '%s' contains '%s'", stashString, needleString)
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < stashValue.Len(); i++ {
			elem := stashValue.Index(i).Interface()
			if reflect.DeepEqual(elem, needle) {
				return true, fmt.Errorf("slice/array contains %#v", needle)
			}
		}

		return false, fmt.Errorf("slice/array does not contain %#v", needle)
	case reflect.Map:
		for _, key := range stashValue.MapKeys() {
			val := stashValue.MapIndex(key).Interface()
			if reflect.DeepEqual(val, needle) {
				return true, fmt.Errorf("Contains failed: map contains key %#v", needle)
			}
		}
		return false, fmt.Errorf("Contains failed: map does not contain key %#v", needle)

	default:
		return false, fmt.Errorf("Contains failed: unsupported type %T", stash)
	}

	return false, fmt.Errorf("Contains failed: unsupported type %T", stash)
}

func Count[T any](t *testing.T, want int, got T, errMsg ...string) {
	t.Helper()

	v := reflect.ValueOf(got)

	switch v.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String, reflect.Chan:
		if v.Len() != want {
			if len(errMsg) > 0 {
				t.Fatal(errMsg[0])
			}
			t.Fatalf("Len failed: want=%d, got=%d", want, v.Len())
		}
	default:
		t.Fatalf("Len failed: type %T does not support len()", got)
	}
}
