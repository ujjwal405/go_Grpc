package service

import "sync"

type RatingStore interface {
	Add(laptopid string, score float64) (*Rating, error)
}
type Rating struct {
	Count uint32
	Sum   float64
}
type Ratingmemory struct {
	mutex  sync.RWMutex
	rating map[string]*Rating
}

func NewRatingMemory() *Ratingmemory {
	return &Ratingmemory{
		rating: make(map[string]*Rating),
	}
}
func (store *Ratingmemory) Add(laptopid string, score float64) (*Rating, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	rating := store.rating[laptopid]
	if rating == nil {
		rating = &Rating{
			Count: 1,
			Sum:   score,
		}
	} else {
		rating.Count++
		rating.Sum += score

	}
	store.rating[laptopid] = rating
	return rating, nil

}
