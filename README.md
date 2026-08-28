# Aimee module: config

This is the independent `config` source-ownership repository.

This component executes inside the trusted C core. It
does not receive a bus principal, a process shim, or a grant. The repository is
the independently versioned source owner consumed by the server/KB core build.


The descriptor-owned production sources, headers, tests, and documentation are
preserved at their canonical paths so their migration history remains auditable.
