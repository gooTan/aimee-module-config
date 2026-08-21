# Aimee module: config

This is the independent `config` source-ownership repository.

It builds `aimee-module-config` as a pure-Go process for the
server, kb bus. The exported repository includes the
exact canonical Go bus client/runtime snapshot, caller contract, handler,
validation, and persistent store. It contains no C implementation or bridge.

The daemon admits the process only when its installed absolute executable path,
UID, principal class, principal reference, and event-kind grants match the
installed `.grant` file. Copy that generated grant into each declared daemon
policy directory under `modules.d`.

The descriptor-owned production sources, tests, contracts,
and documentation are preserved at their canonical paths so their migration
history remains auditable.
