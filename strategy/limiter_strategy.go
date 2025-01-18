package strategy

type LimiterStrategy interface {
	Increment(key string, window int) (int, error)
	Block(key string, duration int) error
	IsBlocked(key string) (bool, error)
}
