package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"go.starlark.net/starlark"

	"tidbyt.dev/pixlet/runtime"
)

const PublicKeysetJSON = `{
  "primaryKeyId": 1589560679,
  "key": [
    {
      "keyData": {
        "typeUrl": "type.googleapis.com/google.crypto.tink.EciesAeadHkdfPublicKey",
        "value": "ElwKBAgCEAMSUhJQCjh0eXBlLmdvb2dsZWFwaXMuY29tL2dvb2dsZS5jcnlwdG8udGluay5BZXNDdHJIbWFjQWVhZEtleRISCgYKAggQEBASCAoECAMQEBAgGAEYARogLGtas20og5yP8/g9mCNLNCWTDeLUdcHH7o9fbzouOQoiIBIth4hdVF5A2sztwfW+hNoZ0ht/HNH3dDTEBPW3GXA2",
        "keyMaterialType": "ASYMMETRIC_PUBLIC"
      },
      "status": "ENABLED",
      "keyId": 1589560679,
      "outputPrefixType": "TINK"
    }
  ]
}`

var EncryptCmd = &cobra.Command{
	Use:     "encrypt [app name] [secret value]...",
	Short:   "Encrypts secrets for use in an app that will be submitted to the Tidbyt community repo",
	Example: "encrypt weather my-top-secretweather-api-key-123456",
	Args:    cobra.MinimumNArgs(2),
	Run:     encrypt,
}

func encrypt(cmd *cobra.Command, args []string) {
	sek := &runtime.SecretEncryptionKey{
		PublicKeysetJSON: []byte(PublicKeysetJSON),
	}

	appName := args[0]
	encrypted := make([]string, len(args)-1)

	for i, val := range args[1:] {
		var err error
		encrypted[i], err = sek.Encrypt(appName, val)
		if err != nil {
			log.Fatalf("encrypting value: %v", err)
		}
	}

	for _, val := range encrypted {
		fmt.Println(starlark.String(val).String())
	}
}
