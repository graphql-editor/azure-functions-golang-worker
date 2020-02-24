package converters_test

import (
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/graphql-editor/azure-functions-golang-worker/converters"
	NullableTypes "github.com/graphql-editor/azure-functions-golang-worker/rpc/shared"
	"github.com/stretchr/testify/assert"
)

func strptr(s string) *string        { return &s }
func boolptr(b bool) *bool           { return &b }
func float64ptr(f float64) *float64  { return &f }
func timeptr(t time.Time) *time.Time { return &t }

func TestNullable(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		data := []struct {
			input    *string
			expected *NullableTypes.NullableString
		}{
			{},
			{input: strptr(""), expected: &NullableTypes.NullableString{
				String_: &NullableTypes.NullableString_Value{
					Value: "",
				},
			}},
			{input: strptr("data"), expected: &NullableTypes.NullableString{
				String_: &NullableTypes.NullableString_Value{
					Value: "data",
				},
			}},
		}
		for _, tc := range data {
			assert.Equal(t, tc.expected, converters.EncodeNullableString(tc.input))
		}
	})
	t.Run("Bool", func(t *testing.T) {
		data := []struct {
			input    *bool
			expected *NullableTypes.NullableBool
		}{
			{},
			{input: boolptr(true), expected: &NullableTypes.NullableBool{
				Bool: &NullableTypes.NullableBool_Value{
					Value: true,
				},
			}},
			{input: boolptr(false), expected: &NullableTypes.NullableBool{
				Bool: &NullableTypes.NullableBool_Value{
					Value: false,
				},
			}},
		}
		for _, tc := range data {
			assert.Equal(t, tc.expected, converters.EncodeNullableBool(tc.input))
		}
	})
	t.Run("Double", func(t *testing.T) {
		data := []struct {
			input    *float64
			expected *NullableTypes.NullableDouble
		}{
			{},
			{input: float64ptr(0.0), expected: &NullableTypes.NullableDouble{
				Double: &NullableTypes.NullableDouble_Value{
					Value: 0.0,
				},
			}},
			{input: float64ptr(1.0), expected: &NullableTypes.NullableDouble{
				Double: &NullableTypes.NullableDouble_Value{
					Value: 1.0,
				},
			}},
		}
		for _, tc := range data {
			assert.Equal(t, tc.expected, converters.EncodeNullableDouble(tc.input))
		}
	})
	t.Run("Timestamp", func(t *testing.T) {
		now := time.Now()
		data := []struct {
			input    *time.Time
			expected *NullableTypes.NullableTimestamp
			err      assert.ErrorAssertionFunc
		}{
			{},
			{
				input: timeptr(time.Time{}.AddDate(2014, 2, 1)),
				expected: &NullableTypes.NullableTimestamp{
					Timestamp: &NullableTypes.NullableTimestamp_Value{
						Value: &timestamp.Timestamp{
							Seconds: time.Time{}.AddDate(2014, 2, 1).Unix(),
							Nanos:   int32(time.Time{}.AddDate(2014, 2, 1).UnixNano() - time.Time{}.AddDate(2014, 2, 1).Unix()*1e9),
						},
					},
				},
			},
			{
				input: &now,
				expected: &NullableTypes.NullableTimestamp{
					Timestamp: &NullableTypes.NullableTimestamp_Value{
						Value: &timestamp.Timestamp{
							Seconds: now.Unix(),
							Nanos:   int32(now.UnixNano() - now.Unix()*1e9),
						},
					},
				},
			},
			{
				// Test negative timestamp behaviour
				input: timeptr(time.Unix(0, 0).Add(time.Second*(-1) + time.Nanosecond*(-1))),
				expected: &NullableTypes.NullableTimestamp{
					Timestamp: &NullableTypes.NullableTimestamp_Value{
						Value: &timestamp.Timestamp{
							Seconds: -2,
							Nanos:   1e9 - 1,
						},
					},
				},
			},
			{
				// Expect error on invalid timestamp, for example before before 0001-01-01
				input: timeptr(time.Time{}.Add(time.Nanosecond * (-1))),
				err:   assert.Error,
			},
		}
		for i := range data {
			if data[i].err == nil {
				data[i].err = assert.NoError
			}
		}
		for _, tc := range data {
			ts, err := converters.EncodeNullableTimestamp(tc.input)
			tc.err(t, err)
			assert.Equal(t, tc.expected, ts)
		}
	})
}
