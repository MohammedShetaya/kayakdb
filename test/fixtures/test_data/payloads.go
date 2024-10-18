package test_data

import "github.com/MohammedShetaya/kayakdb/api"

var GetPayload api.Payload = api.Payload{
	Headers: api.Headers{
		Path: "/get",
	},
	Data: []api.KeyValue{
		{Key: api.Number([]byte{0x00, 0x00, 0x00, 0x0F})},
	},
}

var PutPayload api.Payload = api.Payload{
	Headers: api.Headers{
		Path: "/put",
	},
	Data: []api.KeyValue{
		{Key: api.Number([]byte{0x00, 0x00, 0x00, 0x0F}), Value: api.String("hello")},
	},
}
