package crypt

import "testing"

func TestEncryptDecrypt(t *testing.T) {
	e := BasicEncryption{}
	data := []byte{0xAA, 0x30, 0xFC}
	key := "Ø"
	encrypted := e.Encrypt(data, key)
	decrypted := e.Decrypt(encrypted, key)

	dataStr := string(data[:])
	decryptStr := string(decrypted[:])
	if decryptStr != dataStr {
		t.Fatalf("wrong value - expected: %s got: %s", data, decrypted)
	}
}
