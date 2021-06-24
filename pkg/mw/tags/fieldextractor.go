package tags

import "reflect"

type RequestFieldExtractorFunc func(fullMethod string, req interface{}) map[string]interface{}

type requestFieldExtractor interface {
	ExtractRequestFields(m map[string]interface{})
}

func CodeGenRequestFieldExtractor(fullMethod string, req interface{}) map[string]interface{} {
	if inst, ok := req.(requestFieldExtractor); ok {
		m := make(map[string]interface{})
		inst.ExtractRequestFields(m)
		if len(m) == 0 {
			return nil
		}
	}
	return nil
}

func TagBasedRequestFieldExtractor(tagName string) RequestFieldExtractorFunc {
	return func(fullMethod string, req interface{}) map[string]interface{} {
		m := make(map[string]interface{})
		reflectMessageTags(req, tagName, m)
		if len(m) == 0 {
			return nil
		}
		return m
	}
}

func reflectMessageTags(msg interface{}, tagName string, m map[string]interface{}) {
	v := reflect.ValueOf(msg)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return
	}
	v = v.Elem()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		k := f.Kind()
		// ptr interface
		if (k == reflect.Ptr || k == reflect.Interface) && f.CanInterface() {
			reflectMessageTags(f.Interface(), tagName, m)
		}
		// array
		if k == reflect.Array || k == reflect.Slice {
			if f.Len() == 0 {
				continue
			}
			k = f.Index(0).Kind()
		}
		// handled types
		if (k >= reflect.Bool && k <= reflect.Float64) || k == reflect.String {
			if tag := v.Type().Field(i).Tag.Get(tagName); tag != "" {
				m[tag] = f.Interface()
			}
		}
	}
}
