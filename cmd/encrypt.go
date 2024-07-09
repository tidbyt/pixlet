package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.starlark.net/starlark"

	"tidbyt.dev/pixlet/manifest"
	"tidbyt.dev/pixlet/runtime"
	"tidbyt.dev/pixlet/tools/repo"
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
	Use:     "encrypt [app ID] [secret value]...",
	Short:   "Encrypt a secret for use in the Tidbyt community repo",
	Example: "encrypt weather my-top-secretweather-api-key-123456",
	Args:    cobra.MinimumNArgs(2),
	Run:     encrypt,
}

func encrypt(cmd *cobra.Command, args []string) {
	sek := &runtime.SecretEncryptionKey{
		PublicKeysetJSON: []byte(PublicKeysetJSON),
	}

	appID := args[0]
	if err := validateAppID(args[0]); err != nil {
		log.Fatalf("Cannot encrypt with appID '%s': %v", appID, err)
	}

	encrypted := make([]string, len(args)-1)

	for i, val := range args[1:] {
		var err error
		encrypted[i], err = sek.Encrypt(appID, val)
		if err != nil {
			log.Fatalf("encrypting value: %v", err)
		}
	}

	for _, val := range encrypted {
		fmt.Println(starlark.String(val).String())
	}
}

func validateAppID(appID string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("something went wrong with your local filesystem: %v", err)
	}
	if !repo.IsInRepo(cwd, "community") {
		log.Printf("Skipping validation of appID since command is not being run from inside the community repo.")
		return nil // Can only apply check if running from the community repo
	}
	root, err := repo.RepoRoot(cwd)
	if err != nil {
		return fmt.Errorf("something went wrong with your community repo: %v", err)
	}
	entries, err := os.ReadDir(filepath.Join(root, "apps"))
	if err != nil {
		return fmt.Errorf("something went wrong listing existing apps: %v", err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		f, err := os.Open(filepath.Join(root, "apps", e.Name(), manifest.ManifestFileName))
		if err != nil {
			log.Printf("Skipping %s/%s/%s/%s: %v", root, "apps", e.Name(), manifest.ManifestFileName, err)
			continue
		}
		m, err := manifest.LoadManifest(f)
		if err != nil {
			return fmt.Errorf("something went wrong loading manifest for %s: %v", e.Name(), err)
		}
		if m.ID == appID {
			return nil
		}
	}
	return fmt.Errorf("does not match manifest ID for any app in the community repo")
}
