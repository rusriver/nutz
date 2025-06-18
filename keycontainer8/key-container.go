package keycontainer8

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"
)

// This is a replacement for PEM block, used to pass all necessary metainfo as a single
// yaml-friendly string.
type KeyContainer8 struct {
	Type         string // public | private
	Id           string
	Alg          string // ES256 | EdDSA | etc
	Attrs        map[string]string
	KeyX509Bytes []byte // https://www.itu.int/rec/T-REC-X.509-201910-I/en
}

var ContainerDescriptor = "KEYCONTAINER8"

type GobED struct {
	pipe    bytes.Buffer
	Encoder *gob.Encoder
	Decoder *gob.Decoder
}

func NewGobED() (gobed *GobED) {
	gobed = &GobED{}
	gobed.Encoder = gob.NewEncoder(&gobed.pipe)
	gobed.Decoder = gob.NewDecoder(&gobed.pipe)
	kc := &KeyContainer8{}
	gobed.Encoder.Encode(kc)
	gobed.Decoder.Decode(kc)
	gobed.pipe.Reset()
	return gobed
}

func (kc *KeyContainer8) SerializeBase64(gobed *GobED) (str string) {
	gobed.Encoder.Encode(kc)
	str = fmt.Sprintf("%v--%v--%v--", ContainerDescriptor, kc.Type, kc.Id)
	str += base64.StdEncoding.EncodeToString(gobed.pipe.Bytes())
	gobed.pipe.Reset()
	return
}

func ParseBase64(gobed *GobED, str string) (kc *KeyContainer8, err error) {
	parts := strings.Split(str, "--")
	if len(parts) < 4 || parts[0] != ContainerDescriptor {
		err = errors.New("eebf17d3 wrong format")
		return
	}
	bb, err := base64.StdEncoding.DecodeString(parts[3])
	if err != nil {
		return
	}
	gobed.pipe.Write(bb)
	kc = &KeyContainer8{}
	err = gobed.Decoder.Decode(kc)
	gobed.pipe.Reset()
	return
}
