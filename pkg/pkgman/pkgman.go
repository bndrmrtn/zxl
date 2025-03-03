package pkgman

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const PkgFile = "zxpkg.yaml"

type PackageManager struct {
	PackageName string `yaml:"packageName"`

	Packages []*Package `yaml:"packages"`

	root        string
	packageFile string
}

func New(root string) (*PackageManager, error) {
	var m PackageManager
	m.root = filepath.Clean(root)
	m.packageFile = filepath.Join(m.root, PkgFile)

	f, err := os.Open(m.packageFile)
	if err == nil {
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(data, &m); err != nil {
			return nil, err
		}
	}

	return &m, nil
}

// Add adds a package to pkg.yaml
func (pm *PackageManager) Add(packageUrl string) error {
	pkg, err := pm.parseUrl(packageUrl)
	if err != nil {
		return err
	}

	for _, pack := range pm.Packages {
		if pack.Author == pkg.Author && pack.Package == pkg.Package {
			return fmt.Errorf("package already exists")
		}
	}

	if err := pm.installPackage(pkg); err != nil {
		return err
	}

	pm.Packages = append(pm.Packages, pkg)

	return pm.save()
}

// Remove removes a package from pkg.yaml
func (pm *PackageManager) Remove(packageUrl string) error {
	pkg, err := pm.parseUrl(packageUrl)
	if err != nil {
		return err
	}

	for i, pack := range pm.Packages {
		if pack.Author == pkg.Author && pack.Package == pkg.Package {
			if err := pm.deletePackage(pack); err != nil {
				return err
			}
			pm.Packages = append(pm.Packages[:i], pm.Packages[i+1:]...)
			return pm.save()
		}
	}

	return fmt.Errorf("package not found")
}

// parseUrl parses a package URL and returns a Package struct
func (pm *PackageManager) parseUrl(packageUrl string) (*Package, error) {
	if !strings.HasPrefix(packageUrl, "https://") && !strings.HasPrefix(packageUrl, "git@") {
		return nil, errors.New("invalid package URL")
	}

	// Trim the Git or SSH prefix
	trimmed := strings.TrimPrefix(packageUrl, "git@")
	trimmed = strings.TrimPrefix(trimmed, "https://")
	trimmed = strings.TrimSuffix(trimmed, ".git") // remove .git suffix

	// Split the trimmed URL into parts
	parts := strings.Split(trimmed, "@")
	repoPath := parts[0]
	version := "latest"
	if len(parts) > 1 && parts[1] != "" {
		version = parts[1]
	}

	// Split the repo path into author and package name
	subParts := strings.Split(repoPath, "/")
	if len(subParts) < 2 {
		return nil, errors.New("invalid repository format")
	}
	author := subParts[len(subParts)-2]
	pkgName := subParts[len(subParts)-1]

	return &Package{
		Url:     packageUrl,
		Author:  author,
		Package: pkgName,
		Version: version,
	}, nil
}

// save saves the package manager state to pkg.yaml
func (pm *PackageManager) save() error {
	if pm.PackageName == "" {
		pm.PackageName = filepath.Base(pm.root)
	}

	b, err := yaml.Marshal(pm)
	if err != nil {
		return err
	}

	return os.WriteFile(pm.packageFile, b, os.ModePerm)
}
