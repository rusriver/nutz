package genericmsboot

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rusriver/configdb/v3"
	"github.com/rusriver/nutz/controlflow"
	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
	"github.com/rusriver/nutz/loggeterr-api/trace"
	loggeterrzl "github.com/rusriver/nutz/loggeterr-zl"
)

type Boot struct {
	RootTracer            loggeterrapi.ITracer
	DiagnosticFlags       *map[string]bool // to allow RCU/APR
	ConfDB                *configdb.Root
	RootContext           context.Context
	RootContextCancel     context.CancelFunc
	L1ConfigEffectiveName string
	//
	ProductName                  string
	Copyright                    string
	CallerSkipFrames             int            // used in FILE:LINE diagnostic mode
	L1ConfigFileNameOverride     string         // overrides the arg-1 value
	EnvsPrefix                   string         // can be used to distinguish envs specific to this service
	PathToTheTracerConfigSection string         // in the config file
	L1ConfigFileRequired         bool           //
	LocalMapOverride             map[string]any //
	L1ConfigTextLiteral          string         // can be used e.g. in tests
	TracerDisableOpt01_1         bool           // disable logging of xu
	TracerDisableOpt02_1         bool           // disable logging of xu-ver
	TracerDisableOpt03_1         bool           // disable logging of xu-fingerprint
	TracerDisableOpt04_1         bool           // disable logging of L1-conf-eff-name
}

func NewBoot(of []func(b *Boot)) *Boot {
	var ok bool

	diagnostic_map := make(map[string]bool, 0xF)
	b := &Boot{
		DiagnosticFlags: &diagnostic_map,
	}
	// the defaults below
	b.ProductName = ""
	b.Copyright = ""
	b.RootContext, b.RootContextCancel = context.WithCancel(context.Background())
	b.CallerSkipFrames = -1
	b.L1ConfigFileNameOverride = ""
	b.EnvsPrefix = "71f2-0197e6694528dc94ddc5"
	b.PathToTheTracerConfigSection = "tracer"
	b.L1ConfigFileRequired = false
	b.LocalMapOverride = nil
	b.L1ConfigTextLiteral = ""
	b.TracerDisableOpt01_1 = false
	b.TracerDisableOpt02_1 = false
	b.TracerDisableOpt03_1 = false
	b.TracerDisableOpt04_1 = false
	//
	for _, f := range of {
		f(b)
	}

	// init the primary tracer
	b.RootTracer = loggeterrzl.New(nil)

	ev := b.RootTracer.InfoEvent("BOOT-STARTED")
	if len(b.ProductName) > 0 {
		ev = ev.Str("PRODUCT", b.ProductName)
	}
	if len(b.Copyright) > 0 {
		ev = ev.Str("COPYRIGHT", b.Copyright)
	}
	ev.Send()

	b.RootTracer.SpanBegin("CONFIGDB-INIT")

	// load the L1 config
	controlflow.Try(func() (err any) {
		if len(os.Args) >= 2 {
			b.L1ConfigEffectiveName = os.Args[1]
		}
		if len(b.L1ConfigFileNameOverride) > 0 {
			b.L1ConfigEffectiveName = b.L1ConfigFileNameOverride
		}
		switch {

		case len(b.L1ConfigEffectiveName) == 0 && len(b.L1ConfigTextLiteral) == 0:
			if b.L1ConfigFileRequired {
				controlflow.RaiseCause(nil, &trace.BacktracePoint{"aa9741598829", "L1 config file name is not specified"})
			} else {
				b.RootTracer.InfoEvent("ae6a84057237").SendMsg("missing L1 config, going with empty ConfDB")
				ic := (&configdb.InitContext{}).WithTracer(b.RootTracer)
				b.ConfDB = &configdb.Root{
					InitContext: ic,
				}
				b.ConfDB = b.ConfDB.DontPanic()
			}

		case len(b.L1ConfigTextLiteral) > 0:
			// this takes priority over the CLI args
			ic := (&configdb.InitContext{}).WithTracer(b.RootTracer)
			b.ConfDB = ic.FromBytes([]byte(
				strings.ReplaceAll(b.L1ConfigTextLiteral, "\t", "    "),
			)).Load().DontPanic()

		default:
			ic := (&configdb.InitContext{}).WithTracer(b.RootTracer)
			b.ConfDB = ic.FromFile(b.L1ConfigEffectiveName).LoadWithParenting().DontPanic()

		}
		return nil
	}).Catch(func(e *controlflow.Event) {
		e.Reraise(nil)
	})

	// load the L2 config
	controlflow.Try(func() (err any) {
		if len(b.L1ConfigTextLiteral) > 0 {
			return nil
		}
		var L2_name string
		if len(os.Args) >= 3 {
			L2_name = os.Args[2]
		}
		if len(L2_name) == 0 || L2_name == "." {
			return nil
		}
		ic := (&configdb.InitContext{}).WithTracer(b.RootTracer)
		L2_conf := ic.FromFile(L2_name).LoadWithParenting().DontPanic()
		b.ConfDB.ExtendBy_v2(L2_conf)
		return nil
	}).Catch(func(e *controlflow.Event) {
		e.Reraise(nil)
	})

	b.ConfDB.ExtendByEnvs_WithPrefix(b.EnvsPrefix)

	b.ConfDB.ExtendByEnvsV2_WithPrefix(b.EnvsPrefix, func(s string) string {
		return configdb.TranslateEnvs_KeySuffix(s)
	})

	if b.LocalMapOverride != nil {
		b.ConfDB.ExtendByLocalMap(b.LocalMapOverride)
	}

	b.ConfDB = b.ConfDB.DoPanic()
	{
		b.ConfDB.TheIsa()

		b.ConfDB.DDash()
	}
	b.ConfDB = b.ConfDB.DontPanic()

	// diagnostic flags init
	if len(os.Args) >= 4 {
		diagnosticFlags := os.Args[3:]
		for _, flag := range diagnosticFlags {
			(*b.DiagnosticFlags)[flag] = true
		}
	}

	b.RootTracer.SpanEnd()

	// init/fix the new main working tracer
	{
		// set the legacy log-level if any
		legacy_log_level_str := b.ConfDB.PP(b.PathToTheTracerConfigSection + ".legacy-log-level").String()
		if len(legacy_log_level_str) > 0 {
			switch legacy_log_level_str {
			case "trace", "debug", "info", "warn", "error", "fatal", "panic":
				legacy_log_level, err := zerolog.ParseLevel(legacy_log_level_str)
				if err != nil {
					controlflow.RaiseCause(nil, &trace.BacktracePoint{"821ba70f19e3", err})
				}
				// some implementation-specific magic - to set the legacy level
				loggeterrzl_tracer := b.RootTracer.(*loggeterrzl.Tracer)
				zero_logger := loggeterrzl_tracer.ZeroLogger
				new_zero_logger := zero_logger.Level(legacy_log_level)
				loggeterrzl_tracer.ZeroLogger = &new_zero_logger

			default:
				controlflow.RaiseCause(nil, trace.BacktracePointf(
					"03e5680bf8b7",
					"unknown legacy log-level, must be either [%v]",
					strings.Join([]string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}, " "),
				))
			}
		}

		{
			nstic := b.RootTracer.NSTIC("mwt") // mwt - main working tracer

			// optionally enable the FILE:LINE info
			if (*b.DiagnosticFlags)["FILE:LINE"] {
				nstic = nstic.Caller(b.CallerSkipFrames)
			}

			// set service attributes from the config
			ok = true
			tracerConf := b.ConfDB.WithOk(&ok).P(b.PathToTheTracerConfigSection)
			if ok && !(*b.DiagnosticFlags)["SHORT-192"] {

				xu := tracerConf.P("executable-unit-name").String()
				if len(xu) > 0 && !b.TracerDisableOpt01_1 {
					nstic = nstic.Str("xu", xu)
				}
				xu_ver := tracerConf.P("executable-unit-version").String()
				if len(xu_ver) > 0 && !b.TracerDisableOpt02_1 {
					nstic = nstic.Str("xu-ver", xu_ver)
				}
				xu_fp := tracerConf.P("executable-unit-fingerprint").String() // or git commit hash
				if len(xu_fp) > 0 && !b.TracerDisableOpt03_1 {
					nstic = nstic.Str("xu-fingerprint", xu_fp)
				}
				if len(b.L1ConfigEffectiveName) > 0 && !b.TracerDisableOpt04_1 {
					nstic = nstic.Str("L1-conf-eff-name", b.L1ConfigEffectiveName)
				}
			}

			b.RootTracer = nstic.ITracer()
		}
	} // init new main working tracer

	var MustExit bool

	if (*b.DiagnosticFlags)["envs-printout-and-exit"] {
		fmt.Println("envs-printout-and-exit:")
		for _, pair := range os.Environ() {
			fmt.Println("    ", pair)
		}
		fmt.Println("-- --")
		MustExit = true // yes, we must. It's mandatory here
	}

	if (*b.DiagnosticFlags)["full-configdb-printout-and-exit"] {
		fmt.Println("full-configdb-printout-and-exit:")
		bb := configdb.MustRenderYaml(b.ConfDB)
		fmt.Println(string(bb))
		fmt.Println("-- --")
		MustExit = true // yes, we must. It's mandatory here
	}

	if MustExit {
		os.Exit(1) // we must then
	}

	return b
}

func (b *Boot) BlockAndWait() {
	b.RootTracer.InfoEvent("BAW-227").SendMsg("block and wait")

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)

	// main blocking point
	select {
	case <-chSignal:
		b.RootTracer.InfoEvent("BAW-235").SendMsg("got system interrupt signal")
	case <-b.RootContext.Done():
		b.RootTracer.InfoEvent("BAW-237").SendMsg("root context cancelled")
	}
}
