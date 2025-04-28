package eventbus

type Event interface{}

type Listener func(Event)

var listeners = make(map[string][]Listener)

func Subscribe(eventName string, listener Listener) {
	listeners[eventName] = append(listeners[eventName], listener)
}

func Publish(eventName string, event Event) {
	if ls, ok := listeners[eventName]; ok {
		for _, listener := range ls {
			go listener(event)
		}
	}
}
