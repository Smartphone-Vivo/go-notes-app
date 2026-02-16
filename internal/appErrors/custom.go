package errors

type NotFoundError struct {
	Entity string
	ID     string
}

func (e NotFoundError) Error() string {
	return e.Entity + "with id" + e.ID + "not found"
}

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return "validation failed on field '" + e.Field + "': " + e.Message
}

type DatabaseError struct {
	Err error
	Op  string
}

func (e DatabaseError) Error() string {
	return "database error during" + e.Op + ": " + e.Err.Error()
}

func (e DatabaseError) Unwrap() error {
	return e.Err
}
