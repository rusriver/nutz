# The Chassis 🧬💚

## The config 📗

The chassis can load either of two types of configs - the structural config, and the free-form config.

The structural config, is the config strictly defined by the struct types by at the user side. The free-form config
 specifically is the http://github.com/rusriver/config/v2.

Which path to choose is defined by presense of the StructConfig property - if it's nil, then
the chassis will load the free-form config, else go the struct config initialization path.

Both approaches are quite good, and it's probably more of a taste, which one to use.

### L1 and L2 configs (free-form config only) 📗📔

The XU can accept one or two configs, L1 and L2 respectively. The L1 is most general, the L2 is most specific.
This can be used to have one config with most of the configuration, and another one with some secrets, applied
on top of L1 one.

The L2 config is only loaded if you use the free-form config path, else it doesn't make sense,
and only the L1 is used.

### Passing values from envs 📌

...

## Diagnostic flags 🔨

You can pass extra arguments to the CLI command when executing the XU, and those arguments will
be set in the DiagnosticFlags map. Then, your program can check for them, and do something
diagnostic in nature, if those flags are set. This is indispensable tool for diagnostigs, for
config debug, and for automated testing.

For example, you can output some extra data to be analyzed by auto-testers, or trigger intentionally wrong execution to do test-testing, etc.

## Logger 🐘

The chassis uses the http://github.com/rusriver/nutz/logger, the ILogger interface.

## Metrics 💹

The `/telemetry/prometh-abs` metrics interface is recommended. The drop-in implementation is `/telemetry/prometheus-x`,
which uses the standard Prometheus package for Go.

If you want to use the conventional http://github.com/prometheus/client_golang/prometheus directly, then it's impossible to put that in here;
instead, the metrics must be put as a facet, in `/axu/facets/metrics`.

## USD - User-Side Data 📎📎📎💲

The context: The chassis can be passed from the core down to facets and subsystems, and used there. Sometimes,
those facets and subsystems get the DI (dependency injections) from the core, and inside those injections you might want
to access some data struct from the core, of the core.

Of course you can capture that in a closure. But an alternative is to pass that as the USD, and assert its type
upon use. Sometimes this is very handy.

#### IMPORTANT NOTE:

    You can't pass the kernel as a USD - though that might be ideal, but that's impossible, because it'll inevitably lead to "circular import" error.

    Instead, define these type in the `/axu/core-types` package, and then you'll be able to pass them all over.

## Inventory 📋 & health-checks 💚 & readiness-sensing 👽

...

