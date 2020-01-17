package copier_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"

	"github.com/kmpm/copier"
)

type TsStruct1 struct {
	Field1 int
	Field2 timestamp.Timestamp
}

type TsStruct2 struct {
	Field1 int
	Field2 time.Time
}

// Convert from time.Time to Protobuf Timestamp
func fromTime2Timestamp(to, from reflect.Value) (err error) {
	// log.Println(to.Addr().Type(), "->", from.Addr().Type())

	if _, ok := to.Addr().Interface().(*timestamp.Timestamp); ok {
		if fromTime, ok2 := from.Addr().Interface().(*time.Time); ok2 {
			var ts *timestamp.Timestamp
			ts, err = ptypes.TimestampProto(*fromTime)
			to.Set(reflect.Indirect(reflect.ValueOf(ts)))
		} else {
			err = errors.New("not from time.Time")
		}
	} else {
		err = errors.New("not to timestamp.Timestamp")
	}
	return err
}

// Convert from protobuf Timestamp to time.Time
func fromTimestamp2Time(to, from reflect.Value) (err error) {
	// log.Println(to.Addr().Type(), "->", from.Addr().Type())

	if _, ok := to.Addr().Interface().(*time.Time); ok {
		if fromTimestamp, ok2 := from.Addr().Interface().(**timestamp.Timestamp); ok2 {
			var t time.Time
			t, err = ptypes.Timestamp(*fromTimestamp)

			to.Set(reflect.Indirect(reflect.ValueOf(t)))
		} else if fromTimestamp, ok2 := from.Addr().Interface().(*timestamp.Timestamp); ok2 {
			var t time.Time
			t, err = ptypes.Timestamp(fromTimestamp)
			to.Set(reflect.Indirect(reflect.ValueOf(t)))
		} else {
			err = errors.New("not from timestamp.Timestamp")
		}
	} else {
		err = errors.New("not to time.Time")
	}
	return err
}

func TestCopyFuncTimestamp(t *testing.T) {
	copier.RegisterCopyFunc(
		copier.CopierFunc{
			ToType:   reflect.TypeOf(timestamp.Timestamp{}),
			FromType: reflect.TypeOf(time.Time{}),
			CopyFunc: fromTime2Timestamp,
		},
		copier.CopierFunc{
			ToType:   reflect.TypeOf(time.Time{}),
			FromType: reflect.TypeOf(timestamp.Timestamp{}),
			CopyFunc: fromTimestamp2Time,
		},
	)
}
