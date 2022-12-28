package types

type Bytes []byte

func (b Bytes) MarshalBinary() (data []byte, err error) {
	return b, nil
}
