package account

// Throw by the InventoryClient on validation issues done before sending any http requests
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
