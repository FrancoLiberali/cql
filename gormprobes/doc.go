// Package gormprobes holds probe (characterization) tests that pin down the
// behavior of gorm + SQLite that CQL's fast scanner and cql-gen rely on.
//
// These tests do NOT test CQL code. They exercise gorm/database-sql directly to
// lock in third-party/database semantics that CQL silently depends on:
//
//   - a LEFT JOIN with no matching child row returns all-NULL joined columns
//     (the fast scanner's "no child → skip mount" signal);
//   - gorm leaves a value-type relation as the zero struct on no match
//     (so the scanner's Mount is a no-op for value relations);
//   - gorm refuses uintptr / complex64 / complex128 at AutoMigrate (so cql-gen
//     is right to hard-fail on those field types instead of emitting code that
//     would crash at query time).
//
// If a gorm upgrade changed any of these, the corresponding probe fails here —
// loudly and in isolation — instead of surfacing as a confusing failure deep
// inside the scanner or generator.
package gormprobes
