Possible sources of literals:

    1) Inline literal;

    2) A file embedded as a variable, that variable is used as literal;

    3) Storage is nil, using ParsegetS function loads from the default FS using the os.ReadFile();

    4) Storage is set, using ParsegetS function loads from that using the ReadFile() method;
        The embed.FS is compatible with this.


