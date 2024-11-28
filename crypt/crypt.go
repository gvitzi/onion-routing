package crypt

type Encryption interface {
	Encrypt(data []byte, pubKey string) []byte       // Encrypt using foreign pub key
	Decrypt(data []byte, pvtKey string) []byte       // Decrypt using own pvt key
	Sign(data []byte, pvtKey string) string          // Sign using own pvt key
	VerifySiganture(data []byte, pubKey string) bool // Verify sign using foreing pub key
}
