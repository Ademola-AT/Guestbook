# Design Decisions

This document outlines important design decisions made for this repository and
attempts to provide succinct rationales.

## Errors

-   The callee is expected to return `error`s with messages that include
    information about the particular call, as opposed to the caller adding this
    information. This aligns with common Go practice.

-   Prefer to keep details of returned `error`s unspecified. The most common
    case is that the caller will only care whether an operation succeeds or not.

-   If certain kinds of `error`s are interesting for the caller of a function or
    non-interface method to distinguish, prefer to expose additional information
    through the use of predicate functions like
    [`os.IsNotExist`](https://golang.org/pkg/os/#IsNotExist). This allows the
    internal representation of the `error` to change over time while being
    simple to use.

-   If it is important to distinguish different kinds of `error`s returned from
    an interface method, then the `error` should implement extra methods and the
    interface should document these assumptions. Just remember that each method
    can be implemented independently: if one method is mutually exclusive with
    another, it would be better to return a more complicated data type from one
    methodthan to have separate methods.

-   Transient network errors should be handled by an interface's implementation
    and not bubbled up as a distinguishible error through a generic interface.
    Retry logic is best handled as low in the stack as possible to avoid
    [cascading failure][]. APIs should try to surface "permanent" errors (e.g.
    malformed request, bad permissions) where appropriate so that application
    logic does not attempt to retry idempotent operations, but the
    responsibility is largely on the library, not on the application.

[cascading failure]: https://landing.google.com/sre/book/chapters/addressing-cascading-failures.html
