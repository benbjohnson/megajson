package generator

import "errors"

// Used for marking fields as not written. Not actually returned to user.
var unsupportedTypeError = errors.New("Unsupported type")
var ignoreFieldError = errors.New("Ignore field")
