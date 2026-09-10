// Package mocked holds integration-style tests that drive CQL's full query
// pipeline (build SQL, execute, scan) against a mocked SQL driver (go-sqlmock)
// rather than a real database. They assert the exact SQL CQL emits and feed
// controlled rows back, so they run deterministically with no database — which
// is what the fast-scanner and logging tests here need.
//
// They complement the sibling test/ package, which runs CQL's behavior against
// real databases across every supported dialect.
package mocked
