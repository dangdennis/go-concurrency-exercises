//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import "time"

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

// Dude. It's impossible to cancel the process if process arg doesn't change.
// We can use a cancellation context but our process fn doesn't accept one.
// So for now, we'll just let the process fn run and

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if !u.IsPremium && u.TimeUsed >= 10 {
		return false
	}

	quit := make(chan bool, 1)
	processTime := int64(0)

	go func() {
		startTime := time.Now()
		process()
		elapsedTime := time.Since(startTime)
		processTime = int64(elapsedTime.Seconds())
		quit <- true
	}()

	select {
	case <-quit:
		u.TimeUsed += processTime
		return true
	case <-time.After(time.Second * 10):
		u.TimeUsed += 10
		return false
	}
}

func main() {
	RunMockServer()
}
