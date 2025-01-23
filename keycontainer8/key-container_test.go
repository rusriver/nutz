package keycontainer8_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/rusriver/nutz/controlflow"
	"github.com/rusriver/nutz/keycontainer8"
	"github.com/stretchr/testify/require"
)

func Test_01(t *testing.T) {
	kc1 := &keycontainer8.KeyContainer8{
		Type:         "public",
		Id:           "name-e7d",
		Alg:          "ES256",
		KeyX509Bytes: []byte(`9af083b1-2bb2-4fb8-9351-c4cbcb1c2a0ac25a1183-cc65-4278-b9c7-7575775779f2-just-example-data`),
	}

	gobed := keycontainer8.NewGobED()

	kcEncoded := kc1.SerializeBase64(gobed)
	fmt.Println("1-0", len(kcEncoded), kcEncoded)

	kc2, err := keycontainer8.ParseBase64(gobed, kcEncoded)
	controlflow.Checkerr2("df65fc0c9292", err)

	bb, _ := json.MarshalIndent(kc2, "", "    ")
	fmt.Println("2-0", string(bb))

	require.Equal(t, kc1.KeyX509Bytes, kc2.KeyX509Bytes)
	require.Equal(t, kc1.Type, kc2.Type)
	require.Equal(t, kc1.Id, kc2.Id)
	require.Equal(t, kc1.Alg, kc2.Alg)
}
