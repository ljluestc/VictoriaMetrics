// ...existing code...

func TestInsertNullPrimaryVindex(t *testing.T) {
    // Setup test environment
    // ...existing setup code...

    // Attempt to insert NULL into primary vindex column
    _, err := executor.Execute(ctx, "TestInsert", session, "insert into t1 (c1) values (null)", nil)
    if err == nil || !strings.Contains(err.Error(), "cannot insert NULL value into primary vindex column") {
        t.Errorf("expected error for NULL primary vindex, got: %v", err)
    }

    // ...existing code...
}
