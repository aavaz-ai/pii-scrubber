package piiscrubber

import (
	"reflect"
)

// Borrowed from: https://gist.github.com/hvoecking/10772475

const (
	_piiTag          = "pii"
	_piiTagValueTrue = "true"
)

type taggedField struct {
	text string
	ref  reflect.Value
}

func (s *scrubber) parse(obj interface{}) (interface{}, error) {
	// Wrap the original in a reflect.Value
	original := reflect.ValueOf(obj)

	copy := reflect.New(original.Type()).Elem()
	if err := s.parseRecursive(copy, original, false); err != nil {
		return nil, err
	}

	// Remove the reflection wrapper
	return copy.Interface(), nil
}

func (s *scrubber) parseRecursive(copy, original reflect.Value, hasPIITag bool) error {

	switch original.Kind() {
	// The first cases handle nested structures and parse them recursively

	// If it is a pointer we need to unwrap and call once again
	case reflect.Ptr:
		// To get the actual value of the original we have to call Elem()
		// At the same time this unwraps the pointer so we don't end up in
		// an infinite recursion
		originalValue := original.Elem()
		// Check if the pointer is nil
		if !originalValue.IsValid() {
			return nil
		}
		// Allocate a new object and set the pointer to it
		copy.Set(reflect.New(originalValue.Type()))
		// Unwrap the newly created pointer
		if err := s.parseRecursive(copy.Elem(), originalValue, hasPIITag); err != nil {
			return err
		}

	// If it is an interface (which is very similar to a pointer), do basically the
	// same as for the pointer. Though a pointer is not the same as an interface so
	// note that we have to call Elem() after creating a new object because otherwise
	// we would end up with an actual pointer
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()
		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		if err := s.parseRecursive(copyValue, originalValue, hasPIITag); err != nil {
			return err
		}
		copy.Set(copyValue)

	// If it is a struct we parse each field
	case reflect.Struct:
		t := original.Type()

		for i := 0; i < original.NumField(); i++ {
			piiTagForField := hasPIITag
			tagVal := t.Field(i).Tag.Get(_piiTag)
			if tagVal == "true" {
				piiTagForField = true
			}
			if err := s.parseRecursive(copy.Field(i), original.Field(i), piiTagForField); err != nil {
				return err
			}
		}

	// If it is a slice we create a new slice and parse each element
	case reflect.Slice:
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i++ {
			if err := s.parseRecursive(copy.Index(i), original.Index(i), hasPIITag); err != nil {
				return err
			}
		}

	// If it is a map we create a new map and parse each value
	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			// New gives us a pointer, but again we want the value
			copyValue := reflect.New(originalValue.Type()).Elem()
			if err := s.parseRecursive(copyValue, originalValue, hasPIITag); err != nil {
				return err
			}
			copy.SetMapIndex(key, copyValue)
		}

	// Otherwise we cannot traverse anywhere so this finishes the the recursion

	// If it is a string parse it (yay finally we're doing what we came for)
	case reflect.String:
		// TODO: this is not optimal way to do it
		// In future figure out way to update all the string fields at once
		text := original.String()
		if hasPIITag {
			scrubbedTexts, err := s.ScrubTexts([]string{text})
			if err != nil {
				return err
			}
			text = scrubbedTexts[0]
		}
		copy.SetString(text)

	// And everything else will simply be taken from the original
	default:
		copy.Set(original)
	}

	return nil
}
