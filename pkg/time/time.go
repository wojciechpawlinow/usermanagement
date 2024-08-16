package time

import (
	"time"

	"github.com/wojciechpawlinow/usermanagement/internal/domain"
)

type timeService struct{}

func NewTimeService() domain.TimeProvider {
	return &timeService{}
}

func (*timeService) UtcNow() time.Time {
	loc, _ := time.LoadLocation("")
	now := time.Now()
	now = now.In(loc)

	return now
}
