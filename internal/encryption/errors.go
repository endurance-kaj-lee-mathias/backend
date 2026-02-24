package encryption

import "errors"

var ErrInvalidMasterKey = errors.New("master key must be exactly 32 bytes")
var ErrEncryptionFailed = errors.New("encryption failed")
var ErrDecryptionFailed = errors.New("decryption failed")
var ErrCiphertextTooShort = errors.New("ciphertext is too short")
