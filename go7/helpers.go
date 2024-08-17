package main

func isAlreadyQueued(datagram *Datagram) (bool, error) {
  return CheckCounterParity(datagram)
}
