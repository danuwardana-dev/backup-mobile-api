package helpers

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

func CustomValidatePayload(err error, req interface{}) error {
	var requiredText string
	var wrongFormatText string
	ve, fe := err.(validator.ValidationErrors)
	if fe {
		for i, fieldError := range ve {
			jsonField := getJsonName(fieldError, req)
			if i == 0 {
				if fieldError.Tag() == "required" {
					requiredText = jsonField
					continue
				}
				wrongFormatText = jsonField

			} else {
				if jsonField != requiredText {
					requiredText = requiredText + ", " + jsonField
					continue
				}
				wrongFormatText = wrongFormatText + ", " + jsonField
			}
		}
		if len(wrongFormatText) > 0 && len(requiredText) > 0 {
			return errors.New(fmt.Sprintf("%s %s & %s %s", requiredText, "is require", wrongFormatText, "is wrong format"))
		}
		if len(requiredText) > 0 {
			return errors.New(fmt.Sprintf("%s %s", requiredText, "is required"))
		}
		if len(wrongFormatText) > 0 {
			return errors.New(fmt.Sprintf("%s %s", wrongFormatText, "is wrong format"))
		}
	}
	return err
}
func getJsonName(fieldError validator.FieldError, structModel interface{}) string {
	path := strings.Split(fieldError.StructNamespace(), ".")
	if len(path) > 0 {
		path = path[1:] // Hilangkan nama struct paling atas
	}

	rt := reflect.TypeOf(structModel)
	return findJsonTagName(rt, path)

}
func findJsonTagName(t reflect.Type, path []string) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if len(path) == 0 {
		return ""
	}
	fieldName := path[0]
	field, ok := t.FieldByName(fieldName)
	if !ok {
		return fieldName
	}
	jsonTag := field.Tag.Get("json")
	jsonName := strings.Split(jsonTag, ",")[0]
	if jsonName == "" || jsonName == "-" {
		jsonName = field.Name
	}
	if len(path) == 1 {
		return jsonName
	}
	return findJsonTagName(field.Type, path[1:])
}
