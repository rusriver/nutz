package railroad

type Station[T any] struct {
	NextStation *Station[T]
	LastStation *Station[T]
	V           T
}

func (st *Station[T]) RouteTo(newLastStation *Station[T]) *Station[T] {
	st.NextStation = newLastStation
	if newLastStation.LastStation != nil {
		st.LastStation = newLastStation.LastStation // if we add a chain
	} else {
		st.LastStation = newLastStation // if we add a single one
	}
	return newLastStation
}

type Raiload[T any] struct {
	FirstStation   *Station[T]
	CurrentStation *Station[T]
}

func (rr *Raiload[T]) NextStation() *Station[T] {
	if rr.CurrentStation == nil {
		nextOne := rr.FirstStation
		rr.CurrentStation = nextOne
		return nextOne
	} else {
		nextOne := rr.CurrentStation.NextStation
		rr.CurrentStation = rr.CurrentStation.NextStation
		return nextOne
	}
}

func (rr *Raiload[T]) AddStation(newLastStation *Station[T]) *Station[T] {
	if rr.FirstStation == nil {
		rr.FirstStation = newLastStation
		rr.CurrentStation = nil
	} else {
		lastStation := rr.FirstStation.LastStation
		if lastStation != nil {
			lastStation.RouteTo(newLastStation)
		} else {
			rr.FirstStation.RouteTo(newLastStation)
		}
		rr.FirstStation.LastStation = newLastStation
	}
	return newLastStation
}
