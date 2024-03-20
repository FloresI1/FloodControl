fpackage main

import (
	"context"
	"sync"
	"time"
)

func main() {

}

type FloodControl interface {
	Check(ctx context.Context, userID int64) (bool, error)
}
type MemoryFloodControl struct {
	mu            sync.Mutex
	requestCounts map[int64][]time.Time // Хранение времени запросов для каждого пользователя
	Interval      time.Duration         // Интервал времени для проверки флуда
	MaxRequests   int                   // Максимальное количество запросов за интервал
}
func (m *MemoryFloodControl) Check(ctx context.Context, userID int64) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.requestCounts == nil {
		m.requestCounts = make(map[int64][]time.Time)
	}
	request := m.requestCounts[userID]
	if len(request) > m.MaxRequests {
		return true, nil
	}
	for _, t := range request {
		if time.Since(t) < m.Interval {
			return true, nil
		}
	}
	now := time.Now()
	for i := len(request) - 1; i >= 0; i-- {
		if now.Sub(request[i]) > m.Interval {
			m.requestCounts[userID] = request[i+1:]
			break
		}
	}
	m.requestCounts[userID] = append(m.requestCounts[userID], time.Now())
	return false, nil
}