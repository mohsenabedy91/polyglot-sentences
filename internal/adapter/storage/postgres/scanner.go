package postgres

type Scanner interface {
	Scan(dest ...any) error
}
