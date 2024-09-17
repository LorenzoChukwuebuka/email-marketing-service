package cronjobs

type Job interface {
	Run()
	Schedule() string
}
