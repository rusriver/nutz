This is a complete list, from docs, 27 ones:

	Invalid	x
	Bool	
	Int	
	Int8	
	Int16	
	Int32	
	Int64	
	Uint	
	Uint8	
	Uint16	
	Uint32	
	Uint64	
	Uintptr	x
	Float32	
	Float64	
	Complex64	x
	Complex128	x
	Array	
	Chan	x
	Func	x
	Interface	
	Map	
	Ptr	
	Slice	
	String	
	Struct	
	UnsafePointer	x

The x mark those we don't use / don't need, at least for now.

So I do copy-paste this list, and here's what is left, to make it tidy and clear:

	Bool			a	
	Int			a	
	Int8			a	
	Int16			a	
	Int32			a	
	Int64			a	
	Uint			a	
	Uint8			a	
	Uint16			a	
	Uint32			a	
	Uint64			a	
	Float32			a	
	Float64			a	
	Array		s		=
	Interface	z	s		=
	Map	z	s		=
	Ptr	z	s		=
	Slice	z	s		=
	String	z		a	=
	Struct	z	s		=

20 ones left. We need to handle these at least.

The z marks those which are handled in the gist of Heye Vöcking, out of order (that's why I do the copy-paste, to preserve the order here).

The s marks those who need some kind of recursive handling. The a marks those which need the attempt to convert the type to it.


--
