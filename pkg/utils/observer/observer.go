package observer

type Observer interface {
	Update()
	GetId() string
}

type Subject interface {
	RegisterObserver(observer Observer)
	DeregisterObserver(observer Observer)
	notifyAllObserver()
}
