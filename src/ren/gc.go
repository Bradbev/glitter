package ren

var toCleanUp []func()

func onMainThread(f func()) {
	toCleanUp = append(toCleanUp, f)
}

func GarbageCollect() {
	for _, c := range toCleanUp {
		c()
	}
	toCleanUp = nil
}
