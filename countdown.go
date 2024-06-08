package main

import (
	"time"
)

type CountdownState uint8

const (
	Stopped CountdownState = iota
	Running
	Finished
)

type Countdown struct {
	remaining   time.Duration
	increment   time.Duration
	status      CountdownState
	callback    func()
	startedAt   time.Duration
	stopChannel chan int
}

func NewCountdown(remaining, increment int64, callback func()) *Countdown {
	return &Countdown{
		remaining: time.Duration(remaining) * time.Second,
		increment: time.Duration(increment) * time.Second,
		status:    Stopped, // stopped, running, finished
		callback:  callback,
		// StopChannel: make(chan int),
	}
}

func (countdown *Countdown) Start() {
	if countdown.status != Stopped {
		return
	}
	countdown.startedAt = time.Duration(time.Now().UnixNano()) * time.Nanosecond
	countdown.status = Running
	countdown.stopChannel = make(chan int)
	timer := time.NewTimer(countdown.remaining)
	go func() {
		select {
		case <-timer.C:
			countdown.callback()
			countdown.status = Finished
		case <-countdown.stopChannel:
			timer.Stop()
		}
	}()
}

func (countdown *Countdown) Stop() {
	if countdown.status != Running {
		return
	}
	now := time.Duration(time.Now().UnixNano()) * time.Nanosecond
	countdown.remaining -= countdown.startedAt - now
	countdown.remaining += countdown.increment
	close(countdown.stopChannel)
	countdown.status = Stopped
}

func GetRemaining(countdown Countdown) int {
	return int(countdown.remaining.Seconds())
}
