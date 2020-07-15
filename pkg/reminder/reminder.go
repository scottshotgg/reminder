package reminder

type Reminder interface {
	// encoding.BinaryMarshaler
	// encoding.BinaryUnmarshaler

	// To() string
	// From() string
	// When() <some struct with all times>
	Message() string
}
