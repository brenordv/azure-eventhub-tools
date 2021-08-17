package utils

// TODO: add summary
func FreeSlotOnQueue(c chan bool) {
	<- c
}
