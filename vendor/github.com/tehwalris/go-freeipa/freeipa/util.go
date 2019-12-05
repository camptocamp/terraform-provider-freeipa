package freeipa

// String is a helper to fill *string fields in request options.
func String(v string) *string {
	return &v
}

// Int is a helper to fill *int fields in request options.
func Int(v int) *int {
	return &v
}

// Bool is a helper to fill *bool fields in request options.
func Bool(v bool) *bool {
	return &v
}
