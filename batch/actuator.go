package batch

type Actuator interface {
	Bulk(data []any) error
	Immediate(data any) error
}
