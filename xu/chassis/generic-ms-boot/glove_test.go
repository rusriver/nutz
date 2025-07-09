package genericmsboot_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rusriver/configdb/v3"
	"github.com/rusriver/nutz/controlflow"
	loggeterrzl "github.com/rusriver/nutz/loggeterr-zl"
	genericmsboot "github.com/rusriver/nutz/xu/chassis/generic-ms-boot"
	"github.com/stretchr/testify/require"
)

func Test_Glove_001(t *testing.T) {
	config_text := `
		aaa:
			sss:
				ddd: 123
	`

	tracer := loggeterrzl.New(nil)
	ic := (&configdb.InitContext{}).WithTracer(tracer)
	conf := ic.FromBytes([]byte(strings.ReplaceAll(config_text, "\t", "    "))).Load()

	bb := configdb.MustRenderYaml(conf.Data)
	fmt.Println(string(bb))

	b := genericmsboot.NewBoot([]func(b *genericmsboot.Boot){
		func(b *genericmsboot.Boot) {
			b.L1ConfigTextLiteral = config_text
		},
	})

	bb = configdb.MustRenderYaml(b.ConfDB.Data)
	fmt.Println(string(bb))

	controlflow.Try(func() (err any) {
		g := b.GetGloveAt(nil, "aaa.sss")
		val := g.ConfDB.PP("ddd").String()
		require.Equal(t, "123", val)

		g2 := b.GetGloveAt(nil, "aaa.xxx")
		_ = g2

		return nil
	}).Catch(func(e *controlflow.Event) {
		require.Equal(t, `{: FORWARDTRACE {: ETC [2 4], TAGS [BOOT-CONFIG-DB]}, BACKTRACE { [e6269413b51b "nonexistent map key 'aaa.xxx'"] [7e7e8964 "generic error handler"]}}`, e.String())
	})
}
