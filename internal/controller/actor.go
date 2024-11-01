package controller

type actor struct {
	service serviceActor
}

func NewActor(service serviceActor) *actor {
	return &actor{service: service}
}
