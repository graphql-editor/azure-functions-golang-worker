package converters

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	NullableTypes "github.com/graphql-editor/azure-functions-golang-worker/rpc/shared"
)

// EncodeNullableString supports encoding string to a value that supports null values for given type
func EncodeNullableString(nullable *string) (rpcNullable *NullableTypes.NullableString) {
	if nullable != nil {
		rpcNullable = &NullableTypes.NullableString{
			String_: &NullableTypes.NullableString_Value{
				Value: *nullable,
			},
		}
	}
	return
}

// EncodeNullableBool supports encoding bool to a value that supports null values for given type
func EncodeNullableBool(nullable *bool) (rpcNullable *NullableTypes.NullableBool) {
	if nullable != nil {
		rpcNullable = &NullableTypes.NullableBool{
			Bool: &NullableTypes.NullableBool_Value{
				Value: *nullable,
			},
		}
	}
	return
}

// EncodeNullableDouble supports encoding double to a value that supports null values for given type
func EncodeNullableDouble(nullable *float64) (rpcNullable *NullableTypes.NullableDouble) {
	if nullable != nil {
		rpcNullable = &NullableTypes.NullableDouble{
			Double: &NullableTypes.NullableDouble_Value{
				Value: *nullable,
			},
		}
	}
	return
}

// EncodeNullableTimestamp supports encoding time to a value that supports null values for given type
func EncodeNullableTimestamp(nullable *time.Time) (rpcNullable *NullableTypes.NullableTimestamp, err error) {
	if nullable != nil {
		var t *timestamp.Timestamp
		t, err = ptypes.TimestampProto(*nullable)
		if err == nil {
			rpcNullable = &NullableTypes.NullableTimestamp{
				Timestamp: &NullableTypes.NullableTimestamp_Value{
					Value: t,
				},
			}
		}
	}
	return
}
