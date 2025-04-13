package pkgman

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

const PackageDirectory = ".flmod"

// installPackage installs a package from a given URL and version.
func (pm *PackageManager) installPackage(pkg *Package) error {
	dest := filepath.Join(pm.root, PackageDirectory, pkg.Author, pkg.Package)

	if info, err := os.Stat(dest); err == nil && info.IsDir() {
		fmt.Printf("✅ Package %s already installed\n", pkg.Author+":"+pkg.Package)
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
		if _, err := repo.Reference(ref, true); err != nil {
			// Try as short commit hash
			h, err := repo.ResolveRevision(plumbing.Revision(pkg.Version))
			if err != nil {
				return err
			}
			w, err := repo.Worktree()
			if err != nil {
				return err
			}
			return w.Checkout(&git.CheckoutOptions{
				Hash: *h,
			})
		}
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

	fmt.Printf("✅ Package %s installed successfully\n", pkg.Author+":"+pkg.Package)
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
	pkgDest := filepath.Join(pm.root, PackageDirectory)
	authorDest := filepath.Join(pkgDest, pkg.Author)
	dest := filepath.Join(authorDest, pkg.Package)

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

	// Remove author folder if empty
	empty, err := isDirEmpty(authorDest)
	if err == nil && empty {
		_ = os.Remove(authorDest)
	}

	// Remove main pkg directory if empty
	empty, err = isDirEmpty(pkgDest)
	if err == nil && empty {
		_ = os.Remove(pkgDest)
	}

	return nil
}

func isDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}

	return false, err
}
