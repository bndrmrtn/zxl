package pkgman

// Package represents a package in the package manager
type Package struct {
	// Url is the package repository url
	Url string `yaml:"url"`

	// Author is the author of the package
	Author string `yaml:"-"`
	// Package is the name of the package
	Package string `yaml:"-"`
	// Version is the version of the package
	Version string `yaml:"version"`
}
