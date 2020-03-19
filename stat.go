package sqlstat

type Stat interface {
}

func New() *Stat {}

type collectorKeeper interface {
	GetCollectors()
}

type databaseRegister interface {
	RegisterDB(*sql.DB) error
}

type stat struct{}
