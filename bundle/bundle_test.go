package bundle_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/bundle"
)

func TestBundleWriteAndLoad(t *testing.T) {
	// Ensure we can load the bundle from an app.
	ab, err := bundle.InitFromPath("testdata/testapp")
	assert.NoError(t, err)
	assert.Equal(t, "test-app", ab.Manifest.ID)
	assert.True(t, len(ab.Source) > 0)

	// Create a temp directory.
	dir, err := os.MkdirTemp("", "")
	assert.NoError(t, err)

	// Write bundle to the temp directory.
	err = ab.WriteBundleToPath(dir)
	assert.NoError(t, err)

	// Ensure we can load up the bundle just created.
	path := filepath.Join(dir, bundle.AppBundleName)
	f, err := os.Open(path)
	assert.NoError(t, err)
	defer f.Close()
	newBun, err := bundle.LoadBundle(f)
	assert.NoError(t, err)
	assert.Equal(t, "test-app", newBun.Manifest.ID)
	assert.True(t, len(ab.Source) > 0)
}
func TestLoadBundle(t *testing.T) {
	f, err := os.Open("testdata/bundle.tar.gz")
	assert.NoError(t, err)
	defer f.Close()
	ab, err := bundle.LoadBundle(f)
	assert.NoError(t, err)
	assert.Equal(t, "test-app", ab.Manifest.ID)
	assert.True(t, len(ab.Source) > 0)
}
func TestLoadBundleExcessData(t *testing.T) {
	f, err := os.Open("testdata/excess-files.tar.gz")
	assert.NoError(t, err)
	defer f.Close()

	ab, err := bundle.LoadBundle(f)
	assert.NoError(t, err)
	assert.Equal(t, "test-app", ab.Manifest.ID)
	assert.True(t, len(ab.Source) > 0)
}
