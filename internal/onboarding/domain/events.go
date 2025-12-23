package onboarding_domain


type Event interface {
	Name() string
	Data() map[string]interface{}
}


