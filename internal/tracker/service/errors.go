package service

type appError string

func (e appError) Error() string {
	return string(e)
}

const (
	ErrSubscriptionExists appError = "Subscription already exists"
)


