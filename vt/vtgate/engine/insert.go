// ...existing code...

func (ins *Insert) validatePrimaryVindexValues(values map[string]interface{}) error {
    for col, val := range values {
        if val == nil {
            return fmt.Errorf("cannot insert NULL value into primary vindex column: %s", col)
        }
    }
    return nil
}

// ...existing code...

func (ins *Insert) Execute(ctx context.Context, vcursor VCursor, bindVars map[string]*querypb.BindVariable, wantfields bool) (*sqltypes.Result, error) {
    // ...existing code...

    // Validate primary vindex values
    if err := ins.validatePrimaryVindexValues(bindVars); err != nil {
        return nil, err
    }

    // ...existing code...
}
