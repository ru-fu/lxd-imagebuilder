package generators

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/canonical/lxd-imagebuilder/image"
	"github.com/canonical/lxd-imagebuilder/shared"
)

func TestHostsGeneratorRunLXC(t *testing.T) {
	cacheDir, err := os.MkdirTemp(os.TempDir(), "lxd-imagebuilder-test-")
	require.NoError(t, err)

	rootfsDir := filepath.Join(cacheDir, "rootfs")

	setup(t, cacheDir)
	defer teardown(cacheDir)

	generator, err := Load("hosts", nil, cacheDir, rootfsDir, shared.DefinitionFile{Path: "/etc/hosts"}, shared.Definition{})
	require.IsType(t, &hosts{}, generator)
	require.NoError(t, err)

	definition := shared.Definition{
		Image: shared.DefinitionImage{
			Distribution: "ubuntu",
			Release:      "artful",
		},
	}

	image := image.NewLXCImage(context.TODO(), cacheDir, "", cacheDir, definition)

	err = os.MkdirAll(filepath.Join(cacheDir, "rootfs", "etc"), 0755)
	require.NoError(t, err)

	createTestFile(t, filepath.Join(cacheDir, "rootfs", "etc", "hosts"),
		"127.0.0.1\tlocalhost\n127.0.0.1\tlxd-imagebuilder\n")

	err = generator.RunLXC(image, shared.DefinitionTargetLXC{})
	require.NoError(t, err)

	validateTestFile(t, filepath.Join(cacheDir, "rootfs", "etc", "hosts"),
		"127.0.0.1\tlocalhost\n127.0.0.1\tLXC_NAME\n")
}

func TestHostsGeneratorRunLXD(t *testing.T) {
	cacheDir, err := os.MkdirTemp(os.TempDir(), "lxd-imagebuilder-test-")
	require.NoError(t, err)

	rootfsDir := filepath.Join(cacheDir, "rootfs")

	setup(t, cacheDir)
	defer teardown(cacheDir)

	generator, err := Load("hosts", nil, cacheDir, rootfsDir, shared.DefinitionFile{Path: "/etc/hosts"}, shared.Definition{})
	require.IsType(t, &hosts{}, generator)
	require.NoError(t, err)

	definition := shared.Definition{
		Image: shared.DefinitionImage{
			Distribution: "ubuntu",
			Release:      "artful",
		},
	}

	image := image.NewLXDImage(context.TODO(), cacheDir, "", cacheDir, definition)

	err = os.MkdirAll(filepath.Join(cacheDir, "rootfs", "etc"), 0755)
	require.NoError(t, err)

	createTestFile(t, filepath.Join(cacheDir, "rootfs", "etc", "hosts"),
		"127.0.0.1\tlocalhost\n127.0.0.1\tlxd-imagebuilder\n")

	err = generator.RunLXD(image, shared.DefinitionTargetLXD{})
	require.NoError(t, err)

	validateTestFile(t, filepath.Join(cacheDir, "templates", "hosts.tpl"),
		"127.0.0.1\tlocalhost\n127.0.0.1\t{{ container.name }}\n")
}
