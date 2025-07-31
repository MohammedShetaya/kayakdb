package test_data

import (
	"github.com/MohammedShetaya/kayakdb/types"
)

var GetPayload types.Payload = types.Payload{
	Headers: types.Headers{
		Path: "/get",
	},
	Data: []types.Type{
		types.KeyValue{Key: types.Number([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F})},
	},
}

var PutPayload types.Payload = types.Payload{
	Headers: types.Headers{
		Path: "/put",
	},
	Data: []types.Type{
		types.KeyValue{Key: types.Number([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F}), Value: types.String("hello")},
	},
}
