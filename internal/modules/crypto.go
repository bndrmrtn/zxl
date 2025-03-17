package modules

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"github.com/bndrmrtn/zxl/lang"
	"golang.org/x/crypto/bcrypt"
)

type Crypto struct{}

func NewCrypto() *Crypto {
	return &Crypto{}
}

func (*Crypto) Namespace() string {
	return "crypto"
}

func (c *Crypto) Objects() map[string]lang.Object {
	return map[string]lang.Object{}
}

func (c *Crypto) Methods() map[string]lang.Method {
	return map[string]lang.Method{
		// MD5 hash
		"md5": lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			data := args[0].Value().(string)
			hash := md5.Sum([]byte(data))
			return lang.NewString("md5", hex.EncodeToString(hash[:]), nil), nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "data", Type: lang.TString}),

		// SHA1 hash
		"sha1": lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			data := args[0].Value().(string)
			hash := sha1.Sum([]byte(data))
			return lang.NewString("sha1", hex.EncodeToString(hash[:]), nil), nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "data", Type: lang.TString}),

		// SHA256 hash
		"sha256": lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			data := args[0].Value().(string)
			hash := sha256.Sum256([]byte(data))
			return lang.NewString("sha256", hex.EncodeToString(hash[:]), nil), nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "data", Type: lang.TString}),

		// SHA512 hash
		"sha512": lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			data := args[0].Value().(string)
			hash := sha512.Sum512([]byte(data))
			return lang.NewString("sha512", hex.EncodeToString(hash[:]), nil), nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "data", Type: lang.TString}),

		// BCrypt hash generation
		"bcrypt": lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			password := args[0].Value().(string)
			cost := bcrypt.DefaultCost

			variadicArg := args[1].Value().([]lang.Object)
			if len(variadicArg) > 0 {
				if variadicArg[0].Type() != lang.TInt {
					return nil, fmt.Errorf("invalid argument type for cost: %s, expected int", variadicArg[0].Type())
				}
				cost = int(variadicArg[0].Value().(int))
			}
			hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
			if err != nil {
				return nil, fmt.Errorf("failed to generate bcrypt hash: %w", err)
			}
			return lang.NewString("bcrypt", string(hash), nil), nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "password", Type: lang.TString}).
			WithVariadicArg("cost"),

		// BCrypt verification
		"bcryptCompareHashAndPassword": lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			password := args[0].Value().(string)
			hash := args[1].Value().(string)
			err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
			return lang.NewBool("bcryptVerify", err == nil, nil), nil
		}).WithTypeSafeArgs(
			lang.TypeSafeArg{Name: "password", Type: lang.TString},
			lang.TypeSafeArg{Name: "hash", Type: lang.TString},
		),

		// General hash function that allows selecting algorithm
		"hash": lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			data := args[0].Value().(string)
			algorithm := args[1].Value().(string)

			var hashStr string
			switch algorithm {
			case "md5":
				hash := md5.Sum([]byte(data))
				hashStr = hex.EncodeToString(hash[:])
			case "sha1":
				hash := sha1.Sum([]byte(data))
				hashStr = hex.EncodeToString(hash[:])
			case "sha256":
				hash := sha256.Sum256([]byte(data))
				hashStr = hex.EncodeToString(hash[:])
			case "sha512":
				hash := sha512.Sum512([]byte(data))
				hashStr = hex.EncodeToString(hash[:])
			default:
				return nil, fmt.Errorf("unsupported hash algorithm: %s", algorithm)
			}

			return lang.NewString("hash", hashStr, nil), nil
		}).WithTypeSafeArgs(
			lang.TypeSafeArg{Name: "data", Type: lang.TString},
			lang.TypeSafeArg{Name: "algorithm", Type: lang.TString},
		),
	}
}
