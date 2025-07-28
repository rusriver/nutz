package railroad

import (
	"fmt"
	"testing"
)

type StationObject struct {
	Id string
	F  func(*StationObject)
}

func Test_Railroad_01(t *testing.T) {

	rr := &Raiload[*StationObject]{}

	stFunc1 := func(so *StationObject) {
		fmt.Println("station", so.Id, "type 1")
	}
	stFunc2 := func(so *StationObject) {
		fmt.Println("station", so.Id, "type 2")
		rr.AddStation(&Station[*StationObject]{
			V: &StationObject{Id: "added-01", F: stFunc1},
		})
		rr.AddStation(&Station[*StationObject]{
			V: &StationObject{Id: "added-02", F: stFunc1},
		})
	}

	rr.AddStation(&Station[*StationObject]{
		V: &StationObject{Id: "01", F: stFunc1},
	})
	rr.AddStation(&Station[*StationObject]{
		V: &StationObject{Id: "02", F: stFunc2},
	})
	rr.AddStation(&Station[*StationObject]{
		V: &StationObject{Id: "03", F: stFunc1},
	})

	for {
		st := rr.NextStation()
		if st == nil {
			fmt.Println("FINISHED")
			break
		}
		so := st.V
		so.F(so)
	}

}
