package domain

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/wojciechpawlinow/usermanagement/internal/domain"
)

type TimeProviderMock struct {
	mock.Mock
}

var _ domain.TimeProvider = (*TimeProviderMock)(nil)

func (m *TimeProviderMock) UtcNow() time.Time {
	args := m.Called()

	return args.Get(0).(time.Time)
}
