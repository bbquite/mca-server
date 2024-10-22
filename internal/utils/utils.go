package utils

import (
	"fmt"
	"github.com/bbquite/mca-server/internal/model"
	"reflect"
	"runtime"
)

var floatType = reflect.TypeOf(float64(0))

func GetFloatFromInterface(unk interface{}) (float64, error) {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(floatType)
	return fv.Float(), nil
}

func GetFieldFromMemStats(e *runtime.MemStats, field string) (model.Gauge, error) {
	r := reflect.ValueOf(e)
	f := reflect.Indirect(r).FieldByName(field)
	res, err := GetFloatFromInterface(f.Interface())
	if err != nil {
		return 0, err
	}
	return model.Gauge(res), nil
}
