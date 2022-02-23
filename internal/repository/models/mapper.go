package models

type Mapper interface {
	MapProtoToModel(obj interface{}) (interface{}, error)
	MapModelToProto(obj interface{}) (interface{}, error)
}
