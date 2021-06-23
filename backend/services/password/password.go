package password

type Hasher interface {
	Hash(string) string
	Verify(hash string, plaintext string) error
}
