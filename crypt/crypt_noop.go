package crypt

type NoOpEncryption struct {
}

func (e NoOpEncryption) Encrypt(data []byte, pvtKey string) []byte {
	return data
}

func (e NoOpEncryption) Decrypt(data []byte, pubKey string) []byte {
	return data
}

func (e NoOpEncryption) Sign(data []byte, pvtKey string) string {
	return "not impl"
}

func (e NoOpEncryption) VerifySiganture(data []byte, pubKey string) bool {
	return false
}
