// Package cryptoblob defines the versioned envelope-encryption blob layout
// shared by the encrypt and decrypt packages.
//
// A blob produced by encrypt is an opaque, versioned byte sequence:
//
//	[version 1B][wrapped AES key 256B][nonce 12B][ciphertext + tag ...]
//
// The leading version byte lets decrypt reject blobs it does not understand
// instead of misparsing them.
package cryptoblob

// Version is the leading byte of the current blob layout.
const Version byte = 1
