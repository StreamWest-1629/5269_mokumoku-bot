package apis

var Finalizer = []func(){}

func Finalize() {
	for i := range Finalizer {
		Finalizer[i]()
	}
}
