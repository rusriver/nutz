Some notes:

    Should we use NSTIC() or SpanBegin()/SpanEnd()?

        If your subject is potentially creating parallel processes, or you're sure it does,
        or it's some third-party lib and you're unsure what it does underneath - then the
        NSTIC() would be the choice.

        If you want to make a sub-span in the code under your control, then use SpanBegin()/SpanEnd().

        The Span*() methods modify existing tracer, the same one. This is not intended for concurrent use.
        The NSTIC() creates new tracer, which can be freely used concurrently.

    If we have a recursive algorithm, should we make sub-span or NSTIC inside the called subroutines,
    or just outside?

        Both options are okay, but it's important to stick with one, and not mix them, otherwise you'll
        get a decent portion of mess in your tagging. The recommended way is to make sub-spans and NSTICs
        on the calling side of the subroutines, the practice showed this is better. Because? Take nutz/up
        as an example, if you make a sub in the callee, then you have to do that in every user-defined
        method WindUp(), and if the user forgets that - you tagging gets broken. If you do this on the
        calling side, a user may forget to make a sub at the root, but otherwise the risks of bad tagging
        are minimized.

        The exception is constructors with explicitly specified object name - in that case it looks okay
        if it makes a sub tracer on its own.

    JUST REMEMBER:

        One tracer - one thread.

        NSTIC() creates another tracer. SpanBegin() doesn not.

===
