package runtime

import (
	"bytes"
	"encoding/base64"
	"strings"

	"github.com/google/tink/go/hybrid"
	"github.com/google/tink/go/keyset"
	"github.com/google/tink/go/tink"
	"github.com/pkg/errors"
)

// SecretDecryptionKey is a key that can be used to decrypt secrets.
type SecretDecryptionKey struct {
	// EncryptedKeysetJSON is the encrypted JSON representation of a Tink keyset.
	EncryptedKeysetJSON []byte

	// KeyEncryptionKey is a Tink key that can be used to decrypt the keyset.
	KeyEncryptionKey tink.AEAD
}

// SecretEncryptionKey is a key that can be used to encrypt secrets,
// but not decrypt them.
type SecretEncryptionKey struct {
	// PublicKeysetJSON is the serialized JSON representation of a Tink keyset.
	PublicKeysetJSON []byte
}

func (sdk *SecretDecryptionKey) decrypt(a *Applet) error {
	if a.schema == nil || len(a.schema.Secrets) == 0 {
		// nothing to do
		return nil
	}

	r := bytes.NewReader(sdk.EncryptedKeysetJSON)
	kh, err := keyset.Read(keyset.NewJSONReader(r), sdk.KeyEncryptionKey)
	if err != nil {
		return errors.Wrap(err, "reading keyset JSON")
	}

	dec, err := hybrid.NewHybridDecrypt(kh)
	if err != nil {
		return errors.Wrap(err, "NewHybridDecrypt")
	}

	context := []byte(strings.TrimSuffix(a.Filename, ".star"))

	a.decryptedSecrets = make(map[string]string, len(a.schema.Secrets))
	for k, v := range a.schema.Secrets {
		ciphertext, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return errors.Wrapf(err, "base64 decoding of secret '%s'", k)
		}

		cleartext, err := dec.Decrypt(ciphertext, context)
		if err != nil {
			return errors.Wrapf(err, "decrypting secret '%s'", k)
		}

		a.decryptedSecrets[k] = string(cleartext)
	}

	return nil
}

// Encrypt encrypts a value for use as a secret in an app. Provide both a value
// and the name of the app the encrypted secret will be used in. The value will
// only be usable with the specified app.
func (sek *SecretEncryptionKey) Encrypt(appName, plaintext string) (string, error) {
	r := bytes.NewReader(sek.PublicKeysetJSON)
	kh, err := keyset.ReadWithNoSecrets(keyset.NewJSONReader(r))
	if err != nil {
		return "", errors.Wrap(err, "reading keyset JSON")
	}

	enc, err := hybrid.NewHybridEncrypt(kh)
	if err != nil {
		return "", errors.Wrap(err, "NewHybridEncrypt")
	}

	context := []byte(strings.TrimSuffix(appName, ".star"))
	ciphertext, err := enc.Encrypt([]byte(plaintext), context)
	if err != nil {
		return "", errors.Wrap(err, "encrypting secret")
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
