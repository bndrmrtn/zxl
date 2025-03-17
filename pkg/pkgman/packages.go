package pkgman

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

const PackageDirectory = ".zxpack"

// installPackage installs a package from a given URL and version.
func (pm *PackageManager) installPackage(pkg *Package) error {
	dest := filepath.Join(pm.root, PackageDirectory, pkg.Author, pkg.Package)

	if _, err := os.Stat(dest); err == nil {
		return nil
	}

	err := os.MkdirAll(dest, 0755)
	if err != nil {
		return err
	}

	repo, err := git.PlainClone(dest, false, &git.CloneOptions{
		URL:          pkg.Url,
		SingleBranch: false,
		Depth:        1,
	})

	if err != nil {
		return err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	var ref plumbing.ReferenceName
	if pkg.Version == "latest" {
		head, err := repo.Head()
		if err != nil {
			return err
		}
		ref = head.Name()
	} else {
		ref = plumbing.NewTagReferenceName(pkg.Version)
	}

	// Checkout to the specified version
	err = wt.Checkout(&git.CheckoutOptions{
		Branch: ref,
	})
	if err != nil {
		return err
	}

	// Check if package descriptor exists
	if _, err := os.Stat(filepath.Join(dest, PkgFile)); err == os.ErrNotExist {
		return nil
	}

	manager, err := New(dest)
	if err != nil {
		return err
	}

	for _, pack := range manager.Packages {
		if err := pm.installPackage(pack); err != nil {
			return err
		}
	}

	return nil
}

// isOtherPackageUsing checks if another package is using the specified package.
func (pm *PackageManager) isOtherPackageUsing(pkg *Package) bool {
	for _, pack := range pm.Packages {
		if pack.Author == pkg.Author && pack.Package == pkg.Package {
			continue
		}

		pkgFile := filepath.Join(pm.root, PackageDirectory, pkg.Author, pkg.Package, PkgFile)
		manager, err := New(pkgFile)
		if err != nil {
			continue
		}

		for _, pack := range manager.Packages {
			if pack.Author == pkg.Author && pack.Package == pkg.Package {
				return true
			}
		}
	}
	return false
}

// deletePackage removes an installed package from the package manager.
func (pm *PackageManager) deletePackage(pkg *Package) error {
	dest := filepath.Join(pm.root, PackageDirectory, pkg.Author, pkg.Package)

	if _, err := os.Stat(dest); os.IsNotExist(err) {
		return fmt.Errorf("package not found: %s", pkg.Package)
	}

	pkgFile := filepath.Join(dest, PkgFile)
	manager, err := New(pkgFile)
	if err != nil {
		return err
	}

	for _, pack := range manager.Packages {
		if err := pm.Remove(pack.Url); err != nil {
			return err
		}
	}

	if err = os.RemoveAll(dest); err != nil {
		return fmt.Errorf("failed to delete package: %s, error: %v", pkg.Package, err)
	}

	return nil
}
