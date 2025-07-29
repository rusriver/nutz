package railroad_test

import (
	"fmt"
	"testing"

	"github.com/rusriver/nutz/controlflow/railroad"
)

type RailroadContext struct {
	Raiload *railroad.Raiload[StationFunc]
}

type StationFunc func(rc *RailroadContext)

func Test_Railroad_2_01(t *testing.T) {

	rc := &RailroadContext{
		Raiload: &railroad.Raiload[StationFunc]{},
	}

	rc.Raiload.AddStation(&railroad.Station[StationFunc]{V: stationFunc_1})

	for {
		st := rc.Raiload.NextStation()
		if st == nil {
			fmt.Println("FINISHED")
			break
		}
		sf := st.V
		sf(rc)
	}

}

func stationFunc_1(rc *RailroadContext) {
	fmt.Println("STATION 01")
	rc.Raiload.AddStation(&railroad.Station[StationFunc]{V: stationFunc_2})
}
func stationFunc_2(rc *RailroadContext) {
	fmt.Println("STATION 02")
	rc.Raiload.AddStation(&railroad.Station[StationFunc]{V: stationFunc_3})
}
func stationFunc_3(rc *RailroadContext) {
	fmt.Println("STATION 03")
	rc.Raiload.AddStation(&railroad.Station[StationFunc]{V: stationFunc_4})
}
func stationFunc_4(rc *RailroadContext) {
	fmt.Println("STATION 04")
	rc.Raiload.AddStation(&railroad.Station[StationFunc]{V: stationFunc_5})
}
func stationFunc_5(rc *RailroadContext) {
	fmt.Println("TERMINAL STATION 05")
}
