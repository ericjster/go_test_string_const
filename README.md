# try-string-const

This project tries to find a way to determine if a string is a const, in an efficient way.

The simple idea is to create a long string which is a concatenation of all strings. Then we create the individual strings from it.

In practice we would want to use `go generate` to create the separate strings, very similar to `stringer`.
See this blog post by Rob Pike for details on `go generate` and `stringer`: https://blog.golang.org/generate

Then to determine if a string is a const, we check if the memory for the string is inside a larger string. We expect that can be done in a few machine instructions.

