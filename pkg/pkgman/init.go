package pkgman

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

func NewInitializer(root string) (*PackageManager, error) {
	if info, err := os.Stat(filepath.Join(root, PkgFile)); err == nil && !info.IsDir() {
		return nil, errors.New("package file already exists")
	}

	name, err := getPackageName()
	if err != nil {
		return nil, err
	}

	typ, err := getPackageType()
	if err != nil {
		return nil, err
	}

	pm, err := New(root)
	if err != nil {
		return nil, err
	}

	pm.PackageName = name
	pm.PackageType = typ

	if typ != TypeModule {
		var entry = "main.fl"
		if typ == TypeWeb {
			entry = "public/"
		}

		entryPoint, err := getEntryPoint(entry, typ == TypeCLI)
		if err != nil {
			return nil, err
		}
		pm.PackageConfig["entry"] = entryPoint
	}

	if typ == TypeWeb {
		pm.PackageConfig["web"] = map[string]any{
			"host": "127.0.0.1",
			"port": 3000,
		}
	}

	return pm, nil
}

func getPackageName() (string, error) {
	prompt := promptui.Prompt{
		Label: "Package Name",
		Validate: func(s string) error {
			if len(s) == 0 {
				return errors.New("package name cannot be empty")
			}
			return nil
		},
	}

	return prompt.Run()
}

func getPackageType() (Type, error) {
	items := []Type{TypeModule, TypeCLI, TypeWeb}

	prompt := promptui.Select{
		Label: "Package Type",
		Items: items,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return items[i], nil
}

func getEntryPoint(def string, mustFlare bool) (string, error) {
	prompt := promptui.Prompt{
		Label:   "Entry Point",
		Default: def,
		Validate: func(s string) error {
			if len(s) == 0 {
				return errors.New("entry point cannot be empty")
			}
			if mustFlare && !strings.HasSuffix(s, ".fl") {
				return errors.New("entry point must end with .fl")
			}
			return nil
		},
	}

	return prompt.Run()
}
