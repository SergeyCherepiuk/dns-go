package dns

type Query struct {
	Name  string
	Type  uint16
	Class uint16
}

type QueryPacket struct {
	Header  Header
	Queries []Query
}

// TODO: Not implemented
func MarshalQueryPacket(packet QueryPacket) []byte {
	return nil
}

// TODO: Not implemented
func UnmarshalQeuryPacket(bytes []byte) (QueryPacket, error) {
	return QueryPacket{}, nil
}
