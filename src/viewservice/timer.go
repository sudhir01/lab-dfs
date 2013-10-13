package viewservice

import "time"

type Timer interface {
	 Now() time.Time
}

type DefaultTimer struct {
}

func (this *DefaultTimer) Now() time.Time {
	 return time.Now()
}
