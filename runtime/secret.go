package runtime

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"regexp"
	"sync"

	"github.com/google/tink/go/hybrid"
	"github.com/google/tink/go/keyset"
	"github.com/google/tink/go/tink"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

const (
	threadDecrypterKey = "tidbyt.dev/pixlet/runtime/decrypter"
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

// Encrypt encrypts a value for use as a secret in an app. Provide both a value
// and the ID of the app the encrypted secret will be used in. The value will
// only be usable with the specified app.
func (sek *SecretEncryptionKey) Encrypt(appID, plaintext string) (string, error) {
	r := bytes.NewReader(sek.PublicKeysetJSON)
	kh, err := keyset.ReadWithNoSecrets(keyset.NewJSONReader(r))
	if err != nil {
		return "", fmt.Errorf("%s: %w", "reading keyset JSON", err)
	}

	enc, err := hybrid.NewHybridEncrypt(kh)
	if err != nil {
		return "", fmt.Errorf("%s: %w", "NewHybridEncrypt", err)
	}

	context := []byte(appID)
	ciphertext, err := enc.Encrypt([]byte(plaintext), context)
	if err != nil {
		return "", fmt.Errorf("%s: %w", "encrypting secret", err)
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

var (
	secretOnce   sync.Once
	secretModule starlark.StringDict
)

func LoadSecretModule() (starlark.StringDict, error) {
	secretOnce.Do(func() {
		secretModule = starlark.StringDict{
			"secret": &starlarkstruct.Module{
				Name: "secret",
				Members: starlark.StringDict{
					"decrypt": starlark.NewBuiltin("decrypt", secretDecrypt),
				},
			},
		}
	})

	return secretModule, nil
}

type decrypter func(starlark.String) (starlark.String, error)

func (sdk *SecretDecryptionKey) decrypterForApp(a *Applet) (decrypter, error) {
	r := bytes.NewReader(sdk.EncryptedKeysetJSON)
	kh, err := keyset.Read(keyset.NewJSONReader(r), sdk.KeyEncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "reading keyset JSON", err)
	}

	dec, err := hybrid.NewHybridDecrypt(kh)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "NewHybridDecrypt", err)
	}

	context := []byte(a.ID)

	return func(s starlark.String) (starlark.String, error) {
		v := regexp.MustCompile(`\s`).ReplaceAllString(s.GoString(), "")
		ciphertext, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return "", fmt.Errorf("base64 decoding of secret: %s: %w", s, err)
		}

		cleartext, err := dec.Decrypt(ciphertext, context)
		if err != nil {
			return "", fmt.Errorf("decrypting secret %s: %w", s, err)
		}

		return starlark.String(cleartext), nil
	}, nil
}

func (d decrypter) attachToThread(t *starlark.Thread) {
	t.SetLocal(threadDecrypterKey, d)
}

func decrypterForThread(t *starlark.Thread) decrypter {
	d, ok := t.Local(threadDecrypterKey).(decrypter)
	if ok {
		return d
	} else {
		return nil
	}
}

func secretDecrypt(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var encryptedVal starlark.String

	if err := starlark.UnpackPositionalArgs(
		"decrypt",
		args, kwargs,
		0, &encryptedVal,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for secret.decrypt: %v", err)
	}

	dec := decrypterForThread(thread)

	if dec == nil {
		// no decrypter configured
		return starlark.None, nil
	}

	return dec(encryptedVal)
}
