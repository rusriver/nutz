This has absolutely nothing to do with the bullshit "factory pattern" popularized by
the book "Design Patterns: Elements of Reusable Object-Oriented Software", by
Gamma, Erich; Helm, Richard; Johnson, Ralph; Vlissides, John (1995), Addison-Wesley,
ISBN 0-201-63361-2. Nor is it the https://www.google.com/search?q=oop+factory+pattern.
It has nothing to do with Gang of Four, or Gang of Five, or any other Gang of Bullshit.

Instead, it specifically implements this:
https://www.google.com/search?q=factory+production+line&udm=2
https://en.wikipedia.org/wiki/Production_line
https://en.wikipedia.org/wiki/Assembly_line

This is also close conceptually to the https://www.google.com/search?q=state+machine,
the main difference is that here all states (stations) are grouped in lines. Or you
can think of lines as states, and have more granularity inside, and persistize the
line names, as state IDs, i.e., the persistization granularity would be lines.

Principle of operation:

        factory {:}
            line-A [s1 s2 s3 s4]
            line-B [s1 s2 s3 s4]
            line-C [s1 s2 s3]

	- when we finish last station of line, and there's no intent to jump to another line - stop;
	- we can stop explicitly;
	- can jump to another line from any station;

	Lines are like labels in a program, kind of, if all stations were one single list. This setup
	simplifies management of it and is simpler conceptually.

What this can be used for:

    - For pipelines/sequences of processing operations;

    - For business logic state machines, including persistent ones;

Unrelated links:
    https://refactoring.guru/design-patterns/history
    https://perl.plover.com/yak/design/

