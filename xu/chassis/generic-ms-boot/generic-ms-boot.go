package genericmsboot

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rusriver/configdb/v3"
	"github.com/rusriver/nutz/controlflow"
	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
	"github.com/rusriver/nutz/loggeterr-api/trace"
	loggeterrzl "github.com/rusriver/nutz/loggeterr-zl"
)

type Boot struct {
	Tracer          loggeterrapi.ITracer
	DiagnosticFlags map[string]bool
	ConfDB          *configdb.Root
	//
	CallerSkipFrames             int    // used in FILE:LINE diagnostic mode
	L1ConfigFileNameOverride     string // overrides the arg-1 value
	EnvsPrefix                   string // can be used to distinguish envs specific to this service
	PathToTheTracerConfigSection string // in the config file
	L1ConfigFileRequired         bool   //
}

func NewBoot(of []func(b *Boot)) *Boot {
	var ok bool

	b := &Boot{
		DiagnosticFlags: make(map[string]bool, 0xF),
	}
	// the defaults
	b.CallerSkipFrames = -1
	b.L1ConfigFileNameOverride = ""
	b.EnvsPrefix = ""
	b.PathToTheTracerConfigSection = "tracer"
	b.L1ConfigFileRequired = false
	//
	for _, f := range of {
		f(b)
	}

	// init the primary tracer
	b.Tracer = loggeterrzl.New(nil)

	// load the L1 config
	controlflow.Try(func() (err any) {
		var L1_name string
		if len(os.Args) >= 2 {
			L1_name = os.Args[1]
		}
		if len(b.L1ConfigFileNameOverride) > 0 {
			L1_name = b.L1ConfigFileNameOverride
		}
		if len(L1_name) == 0 {
			if b.L1ConfigFileRequired {
				controlflow.RaiseCause(nil, &trace.BacktracePoint{"aa9741598829", "L1 config file name is not specified"})
			} else {
				b.Tracer.InfoEvent("ae6a84057237").SendMsg("missing L1 config, going with empty ConfDB")
				ic := (&configdb.InitContext{}).WithTracer(b.Tracer)
				b.ConfDB = &configdb.Root{
					InitContext: ic,
				}
				b.ConfDB = b.ConfDB.DontPanic()
				return nil
			}
		}
		ic := (&configdb.InitContext{}).WithTracer(b.Tracer)
		b.ConfDB = ic.FromFile(L1_name).LoadWithParenting().DontPanic()
		return nil
	}).Catch(func(e *controlflow.Event) {
		e.Reraise(nil)
	})

	// load the L2 config
	controlflow.Try(func() (err any) {
		var L2_name string
		if len(os.Args) >= 3 {
			L2_name = os.Args[2]
		}
		if len(L2_name) == 0 || L2_name == "." {
			return nil
		}
		ic := (&configdb.InitContext{}).WithTracer(b.Tracer)
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
			b.DiagnosticFlags[flag] = true
		}
	}

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
				loggeterrzl_tracer := b.Tracer.(*loggeterrzl.Tracer)
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
			nstic := b.Tracer.NSTIC("MWT") // MWT - main working tracer

			// optionally enable the FILE:LINE info
			if b.DiagnosticFlags["FILE:LINE"] {
				nstic = nstic.Caller(b.CallerSkipFrames)
			}

			// set service attributes from the config
			ok = true
			tracerConf := b.ConfDB.WithOk(&ok).P(b.PathToTheTracerConfigSection)
			if ok {

				xu := tracerConf.P("executable-unit-name").String()
				if len(xu) > 0 {
					nstic = nstic.Str("xu", xu)
				}
				xu_ver := tracerConf.P("executable-unit-version").String()
				if len(xu_ver) > 0 {
					nstic = nstic.Str("xu-ver", xu_ver)
				}
				xu_fp := tracerConf.P("executable-unit-fingerprint").String() // or git commit hash
				if len(xu_fp) > 0 {
					nstic = nstic.Str("xu-fingerprint", xu_fp)
				}
			}

			b.Tracer = nstic.ITracer()
		}
	} // init new main working tracer

	var MustExit bool

	if b.DiagnosticFlags["envs-printout-and-exit"] {
		fmt.Println("envs-printout-and-exit:")
		for _, pair := range os.Environ() {
			fmt.Println("    ", pair)
		}
		fmt.Println("-- --")
		MustExit = true // yes, we must. It's mandatory here
	}

	if b.DiagnosticFlags["full-configdb-printout-and-exit"] {
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
