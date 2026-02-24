package encryption

type Service interface {
	GenerateUserEncryptionKey() ([]byte, error)
	EncryptUserKey(userKey []byte) ([]byte, error)
	DecryptUserKey(encrypted []byte) ([]byte, error)
	Encrypt(plaintext []byte, userKey []byte) ([]byte, error)
	Decrypt(ciphertext []byte, userKey []byte) ([]byte, error)
	Hash(value string) string
}

type service struct {
	masterKey []byte
}

func NewService(masterKey []byte) (Service, error) {
	if len(masterKey) != 32 {
		return nil, ErrInvalidMasterKey
	}

	return &service{masterKey: masterKey}, nil
}
