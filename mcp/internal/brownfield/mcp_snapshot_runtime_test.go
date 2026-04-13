package brownfield

import (
	"database/sql"
	"os"
	"strings"
	"testing"
)

func runtimeSQLiteTableMetadata(t *testing.T, s *Store, name string) RuntimeSQLiteTableMetadata {
	t.Helper()

	metadata, err := s.RuntimeSQLiteTableMetadata(name)
	if err != nil {
		t.Fatalf("RuntimeSQLiteTableMetadata(%q): %v", name, err)
	}
	return metadata
}

func runtimeSQLiteTableSchema(t *testing.T, s *Store, name string) SQLiteTableSchema {
	t.Helper()
	return runtimeSQLiteTableMetadata(t, s, name).Table
}

func assertRuntimeSQLiteTableExists(t *testing.T, s *Store, name string) {
	t.Helper()

	if !runtimeSQLiteTableSchema(t, s, name).Exists {
		t.Fatalf("expected runtime table %q to exist", name)
	}
}

func runtimeSQLiteMCPEntryNames(t *testing.T, s *Store) []string {
	t.Helper()

	rows, err := s.db.Query(`
		SELECT name
		FROM brownfield_entries
		WHERE type = 'mcp'
		ORDER BY name ASC
	`)
	if err != nil {
		t.Fatalf("query runtime mcp entry names: %v", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("scan runtime mcp entry name: %v", err)
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate runtime mcp entry names: %v", err)
	}
	return names
}

func assertRuntimeSQLiteMCPEntryNames(t *testing.T, s *Store, want []string) {
	t.Helper()

	got := runtimeSQLiteMCPEntryNames(t, s)
	if len(got) != len(want) {
		t.Fatalf("runtime mcp entry count = %d, want %d (%v)", len(got), len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("runtime mcp entry names = %v, want %v", got, want)
		}
	}

	var duplicateNames int
	if err := s.db.QueryRow(`
		SELECT COUNT(*)
		FROM (
			SELECT name
			FROM brownfield_entries
			WHERE type = 'mcp'
			GROUP BY name
			HAVING COUNT(*) > 1
		)
	`).Scan(&duplicateNames); err != nil {
		t.Fatalf("count duplicate runtime mcp entry names: %v", err)
	}
	if duplicateNames != 0 {
		t.Fatalf("duplicate runtime mcp entry names = %d, want 0", duplicateNames)
	}
}

// Legacy aliases for tests that haven't been updated yet.
var (
	assertRuntimeSQLiteSnapshotNames = assertRuntimeSQLiteMCPEntryNames
)

func runtimeSQLiteMCPEntryByName(t *testing.T, s *Store, name string) MCPServerSnapshot {
	t.Helper()

	var (
		row    MCPServerSnapshot
		dbPath sql.NullString
	)
	if err := s.db.QueryRow(`
		SELECT name, path, desc, is_default, registered_at
		FROM brownfield_entries
		WHERE type = 'mcp' AND key = ?
	`, name).Scan(&row.Name, &dbPath, &row.Desc, &row.IsDefault, &row.RegisteredAt); err != nil {
		if err == sql.ErrNoRows {
			t.Fatalf("runtime mcp entry %q not found", name)
		}
		t.Fatalf("query runtime mcp entry %q: %v", name, err)
	}
	if dbPath.Valid {
		p := dbPath.String
		row.Path = &p
	}
	return row
}

// Legacy alias
var runtimeSQLiteSnapshotRowByName = runtimeSQLiteMCPEntryByName

func assertRuntimeSQLiteMCPDefaultsAllFalse(t *testing.T, s *Store) {
	t.Helper()

	rows, err := s.db.Query(`
		SELECT name, is_default
		FROM brownfield_entries
		WHERE type = 'mcp'
		ORDER BY name ASC
	`)
	if err != nil {
		t.Fatalf("query runtime mcp defaults: %v", err)
	}
	defer rows.Close()

	var rowCount int
	for rows.Next() {
		var (
			name      string
			isDefault bool
		)
		if err := rows.Scan(&name, &isDefault); err != nil {
			t.Fatalf("scan runtime mcp default row: %v", err)
		}
		rowCount++
		if isDefault {
			t.Fatalf("runtime mcp entry %q has is_default=true, want false for new MCP entries", name)
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate runtime mcp defaults: %v", err)
	}
}

// Legacy alias
var assertRuntimeSQLiteSnapshotDefaultsAllFalse = assertRuntimeSQLiteMCPDefaultsAllFalse

func assertRuntimeSQLiteMCPSharedRegisteredAt(t *testing.T, s *Store, wantCount int) string {
	t.Helper()

	rows, err := s.db.Query(`
		SELECT name, registered_at
		FROM brownfield_entries
		WHERE type = 'mcp'
		ORDER BY name ASC
	`)
	if err != nil {
		t.Fatalf("query runtime mcp registered_at rows: %v", err)
	}
	defer rows.Close()

	var (
		rowCount           int
		sharedRegisteredAt string
	)
	for rows.Next() {
		var (
			name         string
			registeredAt string
		)
		if err := rows.Scan(&name, &registeredAt); err != nil {
			t.Fatalf("scan runtime mcp registered_at row: %v", err)
		}
		rowCount++
		registeredAt = strings.TrimSpace(registeredAt)
		if registeredAt == "" {
			t.Fatalf("runtime mcp entry %q missing registered_at", name)
		}
		if sharedRegisteredAt == "" {
			sharedRegisteredAt = registeredAt
			continue
		}
		if registeredAt != sharedRegisteredAt {
			t.Fatalf("runtime mcp entry %q registered_at = %q, want shared scan timestamp %q", name, registeredAt, sharedRegisteredAt)
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate runtime mcp registered_at rows: %v", err)
	}
	if rowCount != wantCount {
		t.Fatalf("runtime mcp registered_at row count = %d, want %d", rowCount, wantCount)
	}
	if sharedRegisteredAt == "" && wantCount > 0 {
		t.Fatal("expected non-empty shared runtime mcp registered_at")
	}
	return sharedRegisteredAt
}

// Legacy alias
var assertRuntimeSQLiteSnapshotSharedRegisteredAt = assertRuntimeSQLiteMCPSharedRegisteredAt

func assertRuntimeSQLiteMetadataUsesDatabasePath(t *testing.T, metadata RuntimeSQLiteTableMetadata, wantPath string) {
	t.Helper()

	gotPath := strings.TrimSpace(metadata.DatabasePath)
	if gotPath == "" {
		t.Fatal("expected runtime sqlite metadata to include a database path")
	}

	gotInfo, err := os.Stat(gotPath)
	if err != nil {
		t.Fatalf("stat runtime sqlite metadata path %q: %v", gotPath, err)
	}
	wantInfo, err := os.Stat(wantPath)
	if err != nil {
		t.Fatalf("stat expected runtime sqlite path %q: %v", wantPath, err)
	}
	if !os.SameFile(gotInfo, wantInfo) {
		t.Fatalf("runtime sqlite metadata path = %q, want same file as %q", gotPath, wantPath)
	}
}
