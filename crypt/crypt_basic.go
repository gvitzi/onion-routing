package crypt

type BasicEncryption struct {
}

func (e BasicEncryption) Encrypt(data []byte, pvtKey string) []byte {
	firstChar := pvtKey[0]

	encrypted := make([]byte, len(data))
	for i, b := range data {
		encrypted[i] = (b + firstChar) % 255
	}

	return encrypted
}

func (e BasicEncryption) Decrypt(data []byte, pubKey string) []byte {
	firstChar := pubKey[0]

	for i, b := range data {
		data[i] = (b - firstChar) % 255
	}

	return data
}

func (e BasicEncryption) Sign(data []byte, pvtKey string) string {
	return "not impl"
}

func (e BasicEncryption) VerifySiganture(data []byte, pubKey string) bool {
	return false
}
