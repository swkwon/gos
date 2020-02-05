package tcp

type Packet struct {
	BodyLength uint32
	Body       []byte
}
