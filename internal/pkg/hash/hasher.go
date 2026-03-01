package hash

type Hasher interface {
	Hash(data []byte) ([]byte, error)
	Verify(data, hash []byte) bool
}
