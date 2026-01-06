package stats

type StatsRepository interface {
	UserCountInWindow(minutes int) (int, error)
}
