package exception

import "errors"

// used for validate users only
var (
	ErrIncorrectPassword error = errors.New("email atau password salah")
	ErrNoSessionFound error= errors.New("no session found,login required")
)

var (
	ErrNotFound error = errors.New("product yang anda minta tidak ditemukan")
)

var (
	ErrNoOrderRow error= errors.New("orderan yang kamu minta tidak ada")
)