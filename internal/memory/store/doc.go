/*
Package store provides persistent storage drivers for the Memory Engine.

The primary abstraction is the contracts.Store interface, which allows the Memory Engine
to interact with different storage backends (defaulting to SQLite) without coupling
to a specific database implementation.

The default implementation uses modernc.org/sqlite, a pure Go SQLite implementation.
*/
package store
