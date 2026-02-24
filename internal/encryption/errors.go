package encryption

import "errors"

var InvalidMasterKey = errors.New("master key must be exactly 32 bytes")
var EncryptionFailed = errors.New("encryption failed")
var DecryptionFailed = errors.New("decryption failed")
var CiphertextTooShort = errors.New("ciphertext is too short")
