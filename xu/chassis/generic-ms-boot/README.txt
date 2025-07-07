Generic micro-service boot boilerplate, wrapped as a package.

What does it do:

    - Inits a tracer, primary one;
        Uses loggeterrzl implementation, zerolog-based, with legacy log-levels and JSON output.

    - Loads the config database:

            - Loads L1 config from the file specified in the 1st CLI argument, with parenting

            - Optionally loads the L2 config from the file specified, if any, in 2nd CLI arg, with parenting

            - Extends L1 with L2

            - Extends the result with envs

            - Executes the $isa processing

            - Executes the DDash DSL processor

    - Inits/fixes a new tracer, working one, parameterized by config section "tracer"

    - Loads the diagnostic flags from the rest of the CLI arguments

===
