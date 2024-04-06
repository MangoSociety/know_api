package questions

type UseCase interface {
	GetQuestion() (string, error)
}
