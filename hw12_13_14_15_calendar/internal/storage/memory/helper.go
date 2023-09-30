package memorystorage

import "github.com/google/uuid"

func (s *Storage) contains(id uuid.UUID) bool {
	_, ok := s.events[id]
	return ok
}

func min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}
