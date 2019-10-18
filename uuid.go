package uuid

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
)

var random io.Reader = rand.Reader

func MakeV3(ns UUID, name string) string {
	return NewV3(ns, name).String()
}

func MakeV4() string {
	return NewV4().String()
}

func MakeV5(ns UUID, name string) string {
	return NewV5(ns, name).String()
}

func NewV3(ns UUID, name string) UUID {
	uuid := newFromHash(md5.New(), ns, name)
	uuid[6] = (uuid[6] & 0x0f) | 0x30
	uuid[8] = (uuid[8] & 0xbf) | 0x80
	return uuid
}

func NewV4() UUID {
	buf := make([]byte, 16)
	random.Read(buf)
	buf[6] = (buf[6] & 0x0f) | 0x40
	buf[8] = (buf[8] & 0xbf) | 0x80
	var uuid UUID
	copy(uuid[:], buf[:])
	return uuid
}

func NewV5(ns UUID, name string) UUID {
	uuid := newFromHash(sha1.New(), ns, name)
	uuid[6] = (uuid[6] & 0x0f) | 0x50
	uuid[8] = (uuid[8] & 0xbf) | 0x80
	return uuid
}

var (
	Nil = UUID{}
)

var (
	NamespaceDNS  = Must(Parse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"))
	NamespaceURL  = Must(Parse("6ba7b811-9dad-11d1-80b4-00c04fd430c8"))
	NamespaceOID  = Must(Parse("6ba7b812-9dad-11d1-80b4-00c04fd430c8"))
	NamespaceX500 = Must(Parse("6ba7b814-9dad-11d1-80b4-00c04fd430c8"))
)

type Variant uint8

// UUID layout variants
const (
	VariantNCS Variant = iota
	VariantRFC4122
	VariantMicrosoft
	VariantFuture
)

func Must(u UUID, err error) UUID {
	if err != nil {
		panic(err)
	}
	return u
}

func Parse(raw string) (UUID, error) {
	u := UUID{}
	err := u.UnmarshalText([]byte(raw))
	return u, err
}

func Equal(a UUID, b UUID) bool {
	return bytes.Equal(a[:], b[:])
}

// The UUID represents a Universally Unique Identifier.
type UUID [16]byte

func (u UUID) Version() uint8 {
	return uint8(u[6] >> 4)
}

func (u UUID) Variant() Variant {
	switch {
	case (u[8] & 0x80) == 0x00:
		return VariantNCS
	case (u[8]&0xc0)|0x80 == 0x80:
		return VariantRFC4122
	case (u[8]&0xe0)|0xc0 == 0xc0:
		return VariantMicrosoft
	}
	return VariantFuture
}

// Returns canonical string representation of UUID
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func (u UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

func (u UUID) Bytes() []byte {
	return u[:]
}

// MarshalText implements the encoding.TextMarshaler interface.
// The encoding is the same as returned by String.
func (u UUID) MarshalText() (text []byte, err error) {
	text = []byte(u.String())
	return
}

var (
	urnPrefix  = []byte("urn:uuid:")
	byteGroups = []int{8, 4, 4, 4, 12}
)

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// Following formats are supported:
// "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
// "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
// "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8"
func (u *UUID) UnmarshalText(text []byte) (err error) {
	if len(text) < 32 {
		err = fmt.Errorf("uuid: UUID string too short: %s", text)
		return
	}

	t := text[:]
	braced := false

	if bytes.Equal(t[:9], urnPrefix) {
		t = t[9:]
	} else if t[0] == '{' {
		braced = true
		t = t[1:]
	}

	b := u[:]

	for i, byteGroup := range byteGroups {
		if i > 0 {
			if t[0] != '-' {
				err = fmt.Errorf("uuid: invalid string format")
				return
			}
			t = t[1:]
		}

		if len(t) < byteGroup {
			err = fmt.Errorf("uuid: UUID string too short: %s", text)
			return
		}

		if i == 4 && len(t) > byteGroup &&
			((braced && t[byteGroup] != '}') || len(t[byteGroup:]) > 1 || !braced) {
			err = fmt.Errorf("uuid: UUID string too long: %s", text)
			return
		}

		_, err = hex.Decode(b[:byteGroup/2], t[:byteGroup])
		if err != nil {
			return
		}

		t = t[byteGroup:]
		b = b[byteGroup/2:]
	}

	return
}

func newFromHash(hash hash.Hash, ns UUID, name string) UUID {
	hash.Write(ns[:])
	hash.Write([]byte(name[:]))

	uuid := UUID{}
	copy(uuid[:], hash.Sum(nil))
	return uuid
}
