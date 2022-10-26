package tcpx

import (
	"encoding/xml"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

type Marshaller interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	MarshalName() string
}

func GetMarshallerByMarshalName(marshalName string) (Marshaller, error) {
	switch marshalName {
	case "json":
		return JsonMarshaller{}, nil
	case "xml":
		return XmlMarshaller{}, nil
	default:
		return nil, errors.New(fmt.Sprintf("unknown marshalName %s,requires in [json,xml,toml,yaml,protobuf]", marshalName))
	}
}

type JsonMarshaller struct{}

func (js JsonMarshaller) Marshal(v interface{}) ([]byte, error) {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(v)
}
func (js JsonMarshaller) Unmarshal(data []byte, dest interface{}) error {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(data, dest)
}

func (js JsonMarshaller) MarshalName() string {
	return "json"
}

type XmlMarshaller struct{}

func (xm XmlMarshaller) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}
func (xm XmlMarshaller) Unmarshal(data []byte, dest interface{}) error {
	return xml.Unmarshal(data, dest)
}

func (xm XmlMarshaller) MarshalName() string {
	return "xml"
}
