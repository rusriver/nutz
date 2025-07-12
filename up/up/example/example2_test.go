package example_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"sync"
	"testing"
	"time"

	"github.com/rusriver/nutz/controlflow"
	loggeterrzl "github.com/rusriver/nutz/loggeterr-zl"
	"github.com/rusriver/nutz/up/fork"
	"github.com/rusriver/nutz/up/reactor"
	"github.com/rusriver/nutz/up/up"
	upcommon "github.com/rusriver/nutz/up/up-common"
	"github.com/rusriver/nutz/up/up/example"
	"github.com/rusriver/nutz/up/watcher"
	"github.com/stretchr/testify/require"
)

// Here we will create a topology with reactors, and will try to use the RecursiveWindUp().
// Consider this a more of a real-life example.

// See also the xu/chassis/generic-ms-boot/idiomatic-simple-example, for a more extended
// form of initialization of sub-systems, though with a rather simple single OTDG.

func Test_DI_Container_Example2_00(t *testing.T) {
	DI_Container_Example2_Main(t, 0)

	/*
		TEST OUTPUT:
			Running tool: /usr/bin/go test -timeout 30s -run ^Test_DI_Container_Example2_00$ github.com/rusriver/nutz/up/up/example -count=1

			=== RUN   Test_DI_Container_Example2_00
			{"level":"INFO","tags":["BEGIN"],"etc":[1],"ts":"20250704-180040"}
			{"level":"INFO","tags":["W01","e325ac80ec09"],"etc":[15,1],"msg":"the watcher W01 starting","ts":"20250704-180040"}
			{"level":"INFO","tags":["W02","e325ac80ec09"],"etc":[16,1],"msg":"the watcher W02 starting","ts":"20250704-180040"}
			{"level":"INFO","tags":["W03","e325ac80ec09"],"etc":[17,1],"msg":"the watcher W03 starting","ts":"20250704-180040"}
			{"level":"INFO","tags":["S_rastan","288"],"etc":[2,1],"msg":"starting","ts":"20250704-180040"}
			{"level":"INFO","tags":["S_rastan","aslast","START"],"etc":[2,2,1],"ts":"20250704-180042"}
			{"level":"INFO","tags":["S_rastan","uphili","START"],"etc":[2,3,1],"ts":"20250704-180042"}
			{"level":"INFO","tags":["S_rastan","389"],"etc":[2,4],"msg":"done","ts":"20250704-180042"}
			{"level":"INFO","tags":["S_myinei","288"],"etc":[4,1],"msg":"starting","ts":"20250704-180042"}
			{"level":"INFO","tags":["S_dodism","288"],"etc":[8,1],"msg":"starting","ts":"20250704-180042"}
			{"level":"INFO","tags":["S_xylane","288"],"etc":[6,1],"msg":"starting","ts":"20250704-180042"}
			{"level":"INFO","tags":["S_obdulu","288"],"etc":[7,1],"msg":"starting","ts":"20250704-180042"}
			{"level":"INFO","tags":["S_armati","288"],"etc":[5,1],"msg":"starting","ts":"20250704-180042"}
			{"level":"INFO","tags":["S_dodism","fopicu","START"],"etc":[8,2,1],"ts":"20250704-180043"}
			{"level":"INFO","tags":["S_dodism","ejacti","START"],"etc":[8,3,1],"ts":"20250704-180043"}
			{"level":"INFO","tags":["S_dodism","hoorgo","START"],"etc":[8,4,1],"ts":"20250704-180043"}
			{"level":"INFO","tags":["S_dodism","volath","START"],"etc":[8,5,1],"ts":"20250704-180043"}
			{"level":"INFO","tags":["S_dodism","vexati","START"],"etc":[8,6,1],"ts":"20250704-180043"}
			{"level":"INFO","tags":["S_dodism","ontras","START"],"etc":[8,7,1],"ts":"20250704-180043"}
			{"level":"INFO","tags":["S_dodism","389"],"etc":[8,8],"msg":"done","ts":"20250704-180043"}
			{"level":"INFO","tags":["S_xylane","requac","START"],"etc":[6,2,1],"ts":"20250704-180043"}
			{"level":"INFO","tags":["S_xylane","rebrei","START"],"etc":[6,3,1],"ts":"20250704-180043"}
			{"level":"INFO","tags":["S_xylane","389"],"etc":[6,4],"msg":"done","ts":"20250704-180043"}
			{"level":"INFO","tags":["S_obdulu","unborr","START"],"etc":[7,2,1],"ts":"20250704-180043"}
			{"level":"INFO","tags":["S_obdulu","axseps","START"],"etc":[7,3,1],"ts":"20250704-180043"}
			{"level":"INFO","tags":["S_obdulu","asdeed","START"],"etc":[7,4,1],"ts":"20250704-180043"}
			{"level":"INFO","tags":["S_obdulu","itchot","START"],"etc":[7,5,1],"ts":"20250704-180043"}
			{"level":"INFO","tags":["S_obdulu","389"],"etc":[7,6],"msg":"done","ts":"20250704-180043"}
			{"level":"INFO","tags":["S_myinei","sostie","START"],"etc":[4,2,1],"ts":"20250704-180043"}
			{"level":"INFO","tags":["S_myinei","razinn","START"],"etc":[4,3,1],"ts":"20250704-180043"}
			{"level":"INFO","tags":["S_myinei","seigul","START"],"etc":[4,4,1],"ts":"20250704-180043"}
			{"level":"INFO","tags":["S_myinei","kaffin","START"],"etc":[4,5,1],"ts":"20250704-180043"}
			{"level":"INFO","tags":["S_myinei","sparip","START"],"etc":[4,6,1],"ts":"20250704-180043"}
			{"level":"INFO","tags":["S_myinei","389"],"etc":[4,7],"msg":"done","ts":"20250704-180043"}
			{"level":"INFO","tags":["S_uromer","288"],"etc":[12,1],"msg":"starting","ts":"20250704-180044"}
			{"level":"INFO","tags":["S_armati","fatoca","START"],"etc":[5,2,1],"ts":"20250704-180044"}
			{"level":"INFO","tags":["S_armati","paudat","START"],"etc":[5,3,1],"ts":"20250704-180044"}
			{"level":"INFO","tags":["S_armati","teashi","START"],"etc":[5,4,1],"ts":"20250704-180044"}
			{"level":"INFO","tags":["S_armati","389"],"etc":[5,5],"msg":"done","ts":"20250704-180044"}
			{"level":"INFO","tags":["S_nutari","288"],"etc":[10,1],"msg":"starting","ts":"20250704-180044"}
			{"level":"INFO","tags":["S_nutari","robles","START"],"etc":[10,2,1],"ts":"20250704-180045"}
			{"level":"INFO","tags":["S_nutari","389"],"etc":[10,3],"msg":"done","ts":"20250704-180045"}
			{"level":"INFO","tags":["S_uromer","petorn","START"],"etc":[12,2,1],"ts":"20250704-180045"}
			{"level":"INFO","tags":["S_uromer","histil","START"],"etc":[12,3,1],"ts":"20250704-180045"}
			{"level":"INFO","tags":["S_uromer","389"],"etc":[12,4],"msg":"done","ts":"20250704-180045"}
			{"level":"INFO","tags":["A1","288"],"etc":[14,1],"msg":"starting","ts":"20250704-180045"}
			{"level":"INFO","tags":["A1","a1_nguloo","START"],"etc":[14,2,1],"ts":"20250704-180046"}
			{"level":"INFO","tags":["A1","a1_fakint","START"],"etc":[14,3,1],"ts":"20250704-180046"}
			{"level":"INFO","tags":["A1","a1_mestic","START"],"etc":[14,4,1],"ts":"20250704-180046"}
			{"level":"INFO","tags":["A1","389"],"etc":[14,5],"msg":"done","ts":"20250704-180046"}
			{"level":"INFO","tags":["W01","5ad2754f035d"],"etc":[15,2],"msg":"the watcher W01 reached STOP/DONE state (OK)","ts":"20250704-180047"}
			{"level":"INFO","tags":["W03","5ad2754f035d"],"etc":[17,2],"msg":"the watcher W03 reached STOP/DONE state (OK)","ts":"20250704-180047"}
			{"level":"INFO","tags":["W02","5ad2754f035d"],"etc":[16,2],"msg":"the watcher W02 reached STOP/DONE state (OK)","ts":"20250704-180047"}
			{"level":"INFO","tags":["UP"],"etc":[18],"ts":"20250704-180048"}
			{
				"Components": {
					"System01": 12,
					"System02": 7,
					"System03": 5,
					"aslast": 1,
					"itchot": 1,
					"ontras": 1,
					"uphili": 1
				}
			}
			OK, 28 components of 7 types are UP
			{"level":"INFO","tags":["GO-DOWN"],"etc":[19],"ts":"20250704-180053"}
			{"level":"INFO","tags":["recursive-wind-up","A1","RWB-33"],"etc":[20,1,1],"msg":"rec-windup-begin","ts":"20250704-180053"}
			{"level":"INFO","tags":["recursive-wind-up","A1","a1_mestic","26"],"etc":[20,1,2,1],"msg":"removed","ts":"20250704-180053"}
			{"level":"INFO","tags":["recursive-wind-up","A1","a1_fakint","26"],"etc":[20,1,3,1],"msg":"removed","ts":"20250704-180053"}
			{"level":"INFO","tags":["recursive-wind-up","A1","a1_nguloo","26"],"etc":[20,1,4,1],"msg":"removed","ts":"20250704-180053"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","RWB-184"],"etc":[20,1,5,1],"msg":"rec-windup-begin","ts":"20250704-180053"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","RWB-33"],"etc":[20,1,5,2,1],"msg":"rec-windup-begin","ts":"20250704-180053"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_myinei","RWB-33"],"etc":[20,1,5,3,1],"msg":"rec-windup-begin","ts":"20250704-180053"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_nutari","RWB-33"],"etc":[20,1,5,4,1],"msg":"rec-windup-begin","ts":"20250704-180053"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_nutari","robles","26"],"etc":[20,1,5,4,2,1],"msg":"removed","ts":"20250704-180054"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","histil","26"],"etc":[20,1,5,2,2,1],"msg":"removed","ts":"20250704-180054"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_myinei","sparip","26"],"etc":[20,1,5,3,2,1],"msg":"removed","ts":"20250704-180054"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_nutari","R3","RWB-184"],"etc":[20,1,5,4,3,1],"msg":"rec-windup-begin","ts":"20250704-180054"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_nutari","R3","S_xylane","RWB-33"],"etc":[20,1,5,4,3,2,1],"msg":"rec-windup-begin","ts":"20250704-180054"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_nutari","R3","S_armati","RWB-33"],"etc":[20,1,5,4,3,3,1],"msg":"rec-windup-begin","ts":"20250704-180054"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_nutari","R3","S_armati","teashi","26"],"etc":[20,1,5,4,3,3,2,1],"msg":"removed","ts":"20250704-180054"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","petorn","26"],"etc":[20,1,5,2,3,1],"msg":"removed","ts":"20250704-180054"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","RWB-184"],"etc":[20,1,5,2,4,1],"msg":"rec-windup-begin","ts":"20250704-180054"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_nutari","R3","S_xylane","rebrei","26"],"etc":[20,1,5,4,3,2,2,1],"msg":"removed","ts":"20250704-180054"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_myinei","kaffin","26"],"etc":[20,1,5,3,3,1],"msg":"removed","ts":"20250704-180054"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_dodism","RWB-33"],"etc":[20,1,5,2,4,2,1],"msg":"rec-windup-begin","ts":"20250704-180054"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_obdulu","RWB-33"],"etc":[20,1,5,2,4,3,1],"msg":"rec-windup-begin","ts":"20250704-180054"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_obdulu","itchot","26"],"etc":[20,1,5,2,4,3,2,1],"msg":"removed","ts":"20250704-180054"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_myinei","seigul","26"],"etc":[20,1,5,3,4,1],"msg":"removed","ts":"20250704-180054"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_dodism","ontras","26"],"etc":[20,1,5,2,4,2,2,1],"msg":"removed","ts":"20250704-180054"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_nutari","R3","S_armati","paudat","26"],"etc":[20,1,5,4,3,3,3,1],"msg":"removed","ts":"20250704-180054"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_nutari","R3","S_xylane","requac","26"],"etc":[20,1,5,4,3,2,3,1],"msg":"removed","ts":"20250704-180054"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_nutari","R3","S_xylane","F1","f33"],"etc":[20,1,5,4,3,2,4,1],"msg":"rec-windup-call","ts":"20250704-180054"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_nutari","R3","S_xylane","RWE-34"],"etc":[20,1,5,4,3,2,5],"msg":"rec-windup-end","ts":"20250704-180054"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_myinei","razinn","26"],"etc":[20,1,5,3,5,1],"msg":"removed","ts":"20250704-180055"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_nutari","R3","S_armati","fatoca","26"],"etc":[20,1,5,4,3,3,4,1],"msg":"removed","ts":"20250704-180055"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_nutari","R3","S_armati","F1","f33"],"etc":[20,1,5,4,3,3,5,1],"msg":"rec-windup-call","ts":"20250704-180055"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_obdulu","asdeed","26"],"etc":[20,1,5,2,4,3,3,1],"msg":"removed","ts":"20250704-180055"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_dodism","vexati","26"],"etc":[20,1,5,2,4,2,3,1],"msg":"removed","ts":"20250704-180055"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_nutari","R3","S_armati","RWE-34"],"etc":[20,1,5,4,3,3,6],"msg":"rec-windup-end","ts":"20250704-180055"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_nutari","R3","RWE-185"],"etc":[20,1,5,4,3,4],"msg":"rec-windup-end","ts":"20250704-180055"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_nutari","RWE-34"],"etc":[20,1,5,4,4],"msg":"rec-windup-end","ts":"20250704-180055"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_obdulu","axseps","26"],"etc":[20,1,5,2,4,3,4,1],"msg":"removed","ts":"20250704-180055"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_myinei","sostie","26"],"etc":[20,1,5,3,6,1],"msg":"removed","ts":"20250704-180055"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_myinei","F1","f33"],"etc":[20,1,5,3,7,1],"msg":"rec-windup-call","ts":"20250704-180055"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_myinei","RWE-34"],"etc":[20,1,5,3,8],"msg":"rec-windup-end","ts":"20250704-180055"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_dodism","volath","26"],"etc":[20,1,5,2,4,2,4,1],"msg":"removed","ts":"20250704-180055"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_dodism","hoorgo","26"],"etc":[20,1,5,2,4,2,5,1],"msg":"removed","ts":"20250704-180055"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_obdulu","unborr","26"],"etc":[20,1,5,2,4,3,5,1],"msg":"removed","ts":"20250704-180055"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_obdulu","F1","f33"],"etc":[20,1,5,2,4,3,6,1],"msg":"rec-windup-call","ts":"20250704-180055"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_obdulu","RWE-34"],"etc":[20,1,5,2,4,3,7],"msg":"rec-windup-end","ts":"20250704-180055"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_dodism","ejacti","26"],"etc":[20,1,5,2,4,2,6,1],"msg":"removed","ts":"20250704-180056"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_dodism","fopicu","26"],"etc":[20,1,5,2,4,2,7,1],"msg":"removed","ts":"20250704-180056"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_dodism","F1","f33"],"etc":[20,1,5,2,4,2,8,1],"msg":"rec-windup-call","ts":"20250704-180056"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_dodism","F1","RWB-39"],"etc":[20,1,5,2,4,2,8,2],"msg":"rec-windup-begin","ts":"20250704-180056"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_dodism","F1","F1","S_rastan","RWB-33"],"etc":[20,1,5,2,4,2,8,3,1,1],"msg":"rec-windup-begin","ts":"20250704-180056"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_dodism","F1","F1","S_rastan","uphili","26"],"etc":[20,1,5,2,4,2,8,3,1,2,1],"msg":"removed","ts":"20250704-180056"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_dodism","F1","F1","S_rastan","aslast","26"],"etc":[20,1,5,2,4,2,8,3,1,3,1],"msg":"removed","ts":"20250704-180056"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_dodism","F1","F1","S_rastan","RWE-34"],"etc":[20,1,5,2,4,2,8,3,1,4],"msg":"rec-windup-end","ts":"20250704-180056"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_dodism","F1","RWE-40"],"etc":[20,1,5,2,4,2,8,4],"msg":"rec-windup-end","ts":"20250704-180056"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","S_dodism","RWE-34"],"etc":[20,1,5,2,4,2,9],"msg":"rec-windup-end","ts":"20250704-180056"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","R2","RWE-185"],"etc":[20,1,5,2,4,4],"msg":"rec-windup-end","ts":"20250704-180056"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","S_uromer","RWE-34"],"etc":[20,1,5,2,5],"msg":"rec-windup-end","ts":"20250704-180056"}
			{"level":"INFO","tags":["recursive-wind-up","A1","R1","RWE-185"],"etc":[20,1,5,5],"msg":"rec-windup-end","ts":"20250704-180056"}
			{"level":"INFO","tags":["recursive-wind-up","A1","RWE-34"],"etc":[20,1,6],"msg":"rec-windup-end","ts":"20250704-180056"}
			{
				"Components": {
					"System01": 0,
					"System02": 0,
					"System03": 0,
					"aslast": 0,
					"itchot": 0,
					"ontras": 0,
					"uphili": 0
				}
			}
			OK, all components are DOWN
			{"level":"INFO","tags":["END"],"etc":[21],"ts":"20250704-180056"}
			--- PASS: Test_DI_Container_Example2_00 (16.01s)
			PASS
			ok      github.com/rusriver/nutz/up/up/example  16.011s
	*/
}

func Test_DI_Container_Example2_01(t *testing.T) {
	DI_Container_Example2_Main(t, 1)

	/*
		TEST OUTPUT:
			Running tool: /usr/bin/go test -timeout 30s -run ^Test_DI_Container_Example2_01$ github.com/rusriver/nutz/up/up/example -count=1

			=== RUN   Test_DI_Container_Example2_01
			{"level":"INFO","tags":["BEGIN"],"etc":[1],"ts":"20250704-180157"}
			{"level":"INFO","tags":["W01","e325ac80ec09"],"etc":[15,1],"msg":"the watcher W01 starting","ts":"20250704-180157"}
			{"level":"INFO","tags":["W02","e325ac80ec09"],"etc":[16,1],"msg":"the watcher W02 starting","ts":"20250704-180157"}
			{"level":"INFO","tags":["W03","e325ac80ec09"],"etc":[17,1],"msg":"the watcher W03 starting","ts":"20250704-180157"}
			{"level":"INFO","tags":["S_rastan","288"],"etc":[2,1],"msg":"starting","ts":"20250704-180157"}
			{"level":"INFO","tags":["S_rastan","aslast","START"],"etc":[2,2,1],"ts":"20250704-180158"}
			{"level":"INFO","tags":["S_rastan","uphili","START"],"etc":[2,3,1],"ts":"20250704-180158"}
			{"level":"INFO","tags":["S_rastan","389"],"etc":[2,4],"msg":"done","ts":"20250704-180158"}
			{"level":"INFO","tags":["S_dodism","288"],"etc":[8,1],"msg":"starting","ts":"20250704-180158"}
			{"level":"INFO","tags":["S_armati","288"],"etc":[5,1],"msg":"starting","ts":"20250704-180158"}
			{"level":"INFO","tags":["S_myinei","288"],"etc":[4,1],"msg":"starting","ts":"20250704-180158"}
			{"level":"INFO","tags":["S_xylane","288"],"etc":[6,1],"msg":"starting","ts":"20250704-180158"}
			{"level":"INFO","tags":["S_obdulu","288"],"etc":[7,1],"msg":"starting","ts":"20250704-180158"}
			{"level":"INFO","tags":["S_obdulu","unborr","START"],"etc":[7,2,1],"ts":"20250704-180159"}
			{"level":"INFO","tags":["S_obdulu","axseps","START"],"etc":[7,3,1],"ts":"20250704-180159"}
			{"level":"INFO","tags":["S_obdulu","asdeed","START"],"etc":[7,4,1],"ts":"20250704-180159"}
			{"level":"INFO","tags":["S_obdulu","itchot","START"],"etc":[7,5,1],"ts":"20250704-180159"}
			{"level":"INFO","tags":["S_obdulu","389"],"etc":[7,6],"msg":"done","ts":"20250704-180159"}
			{"level":"INFO","tags":["S_armati","fatoca","START"],"etc":[5,2,1],"ts":"20250704-180159"}
			{"level":"INFO","tags":["S_armati","paudat","START"],"etc":[5,3,1],"ts":"20250704-180159"}
			{"level":"INFO","tags":["S_armati","teashi","START"],"etc":[5,4,1],"ts":"20250704-180159"}
			{"level":"INFO","tags":["S_armati","389"],"etc":[5,5],"msg":"done","ts":"20250704-180159"}
			{"level":"INFO","tags":["S_myinei","sostie","START"],"etc":[4,2,1],"ts":"20250704-180200"}
			{"level":"INFO","tags":["S_myinei","razinn","START"],"etc":[4,3,1],"ts":"20250704-180200"}
			{"level":"INFO","tags":["S_myinei","seigul","START"],"etc":[4,4,1],"ts":"20250704-180200"}
			{"level":"INFO","tags":["S_myinei","kaffin","START"],"etc":[4,5,1],"ts":"20250704-180200"}
			{"level":"INFO","tags":["S_myinei","sparip","START"],"etc":[4,6,1],"ts":"20250704-180200"}
			{"level":"INFO","tags":["S_myinei","389"],"etc":[4,7],"msg":"done","ts":"20250704-180200"}
			{"level":"INFO","tags":["S_dodism","fopicu","START"],"etc":[8,2,1],"ts":"20250704-180200"}
			{"level":"INFO","tags":["S_dodism","ejacti","START"],"etc":[8,3,1],"ts":"20250704-180200"}
			{"level":"INFO","tags":["S_dodism","hoorgo","START"],"etc":[8,4,1],"ts":"20250704-180200"}
			{"level":"INFO","tags":["S_dodism","volath","START"],"etc":[8,5,1],"ts":"20250704-180200"}
			{"level":"INFO","tags":["S_dodism","vexati","START"],"etc":[8,6,1],"ts":"20250704-180200"}
			{"level":"INFO","tags":["S_dodism","ontras","START"],"etc":[8,7,1],"ts":"20250704-180200"}
			{"level":"INFO","tags":["S_dodism","389"],"etc":[8,8],"msg":"done","ts":"20250704-180200"}
			{"level":"INFO","tags":["S_xylane","requac","START"],"etc":[6,2,1],"ts":"20250704-180200"}
			{"level":"INFO","tags":["S_xylane","rebrei","START"],"etc":[6,3,1],"ts":"20250704-180200"}
			{"level":"INFO","tags":["S_xylane","389"],"etc":[6,4],"msg":"done","ts":"20250704-180200"}
			{"level":"INFO","tags":["S_uromer","288"],"etc":[12,1],"msg":"starting","ts":"20250704-180200"}
			{"level":"INFO","tags":["S_nutari","288"],"etc":[10,1],"msg":"starting","ts":"20250704-180200"}
			{"level":"INFO","tags":["S_nutari","robles","START"],"etc":[10,2,1],"ts":"20250704-180201"}
			{"level":"INFO","tags":["S_nutari","389"],"etc":[10,3],"msg":"done","ts":"20250704-180201"}
			{"level":"INFO","tags":["S_uromer","petorn","START"],"etc":[12,2,1],"ts":"20250704-180201"}
			{"level":"INFO","tags":["S_uromer","histil","START"],"etc":[12,3,1],"ts":"20250704-180201"}
			{"level":"INFO","tags":["S_uromer","389"],"etc":[12,4],"msg":"done","ts":"20250704-180201"}
			{"level":"INFO","tags":["W02","b7d60e67e04c"],"etc":[16,2],"msg":"the watcher W02 timed out","ts":"20250704-180212"}
			{"level":"INFO","tags":["W01","b7d60e67e04c"],"etc":[15,2],"msg":"the watcher W01 timed out","ts":"20250704-180212"}
			{"level":"INFO","tags":["W03","b7d60e67e04c"],"etc":[17,2],"msg":"the watcher W03 timed out","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","RWB-33"],"etc":[18,1,1],"msg":"rec-windup-begin","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","RWB-184"],"etc":[18,1,2,1],"msg":"rec-windup-begin","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_myinei","RWB-33"],"etc":[18,1,2,2,1],"msg":"rec-windup-begin","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_nutari","RWB-33"],"etc":[18,1,2,3,1],"msg":"rec-windup-begin","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","RWB-33"],"etc":[18,1,2,4,1],"msg":"rec-windup-begin","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","histil","26"],"etc":[18,1,2,4,2,1],"msg":"removed","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_myinei","sparip","26"],"etc":[18,1,2,2,2,1],"msg":"removed","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_nutari","robles","26"],"etc":[18,1,2,3,2,1],"msg":"removed","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_nutari","R3","RWB-184"],"etc":[18,1,2,3,3,1],"msg":"rec-windup-begin","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_nutari","R3","S_xylane","RWB-33"],"etc":[18,1,2,3,3,2,1],"msg":"rec-windup-begin","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_nutari","R3","S_armati","RWB-33"],"etc":[18,1,2,3,3,3,1],"msg":"rec-windup-begin","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_nutari","R3","S_armati","teashi","26"],"etc":[18,1,2,3,3,3,2,1],"msg":"removed","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_myinei","kaffin","26"],"etc":[18,1,2,2,3,1],"msg":"removed","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_nutari","R3","S_xylane","rebrei","26"],"etc":[18,1,2,3,3,2,2,1],"msg":"removed","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","petorn","26"],"etc":[18,1,2,4,3,1],"msg":"removed","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","RWB-184"],"etc":[18,1,2,4,4,1],"msg":"rec-windup-begin","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_obdulu","RWB-33"],"etc":[18,1,2,4,4,2,1],"msg":"rec-windup-begin","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_dodism","RWB-33"],"etc":[18,1,2,4,4,3,1],"msg":"rec-windup-begin","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_dodism","ontras","26"],"etc":[18,1,2,4,4,3,2,1],"msg":"removed","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_myinei","seigul","26"],"etc":[18,1,2,2,4,1],"msg":"removed","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_nutari","R3","S_armati","paudat","26"],"etc":[18,1,2,3,3,3,3,1],"msg":"removed","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_obdulu","itchot","26"],"etc":[18,1,2,4,4,2,2,1],"msg":"removed","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_nutari","R3","S_xylane","requac","26"],"etc":[18,1,2,3,3,2,3,1],"msg":"removed","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_nutari","R3","S_xylane","F1","f33"],"etc":[18,1,2,3,3,2,4,1],"msg":"rec-windup-call","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_nutari","R3","S_xylane","RWE-34"],"etc":[18,1,2,3,3,2,5],"msg":"rec-windup-end","ts":"20250704-180217"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_nutari","R3","S_armati","fatoca","26"],"etc":[18,1,2,3,3,3,4,1],"msg":"removed","ts":"20250704-180218"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_nutari","R3","S_armati","F1","f33"],"etc":[18,1,2,3,3,3,5,1],"msg":"rec-windup-call","ts":"20250704-180218"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_nutari","R3","S_armati","RWE-34"],"etc":[18,1,2,3,3,3,6],"msg":"rec-windup-end","ts":"20250704-180218"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_nutari","R3","RWE-185"],"etc":[18,1,2,3,3,4],"msg":"rec-windup-end","ts":"20250704-180218"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_nutari","RWE-34"],"etc":[18,1,2,3,4],"msg":"rec-windup-end","ts":"20250704-180218"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_myinei","razinn","26"],"etc":[18,1,2,2,5,1],"msg":"removed","ts":"20250704-180218"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_obdulu","asdeed","26"],"etc":[18,1,2,4,4,2,3,1],"msg":"removed","ts":"20250704-180218"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_dodism","vexati","26"],"etc":[18,1,2,4,4,3,3,1],"msg":"removed","ts":"20250704-180218"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_myinei","sostie","26"],"etc":[18,1,2,2,6,1],"msg":"removed","ts":"20250704-180218"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_myinei","F1","f33"],"etc":[18,1,2,2,7,1],"msg":"rec-windup-call","ts":"20250704-180218"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_myinei","RWE-34"],"etc":[18,1,2,2,8],"msg":"rec-windup-end","ts":"20250704-180218"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_obdulu","axseps","26"],"etc":[18,1,2,4,4,2,4,1],"msg":"removed","ts":"20250704-180218"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_dodism","volath","26"],"etc":[18,1,2,4,4,3,4,1],"msg":"removed","ts":"20250704-180218"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_obdulu","unborr","26"],"etc":[18,1,2,4,4,2,5,1],"msg":"removed","ts":"20250704-180218"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_dodism","hoorgo","26"],"etc":[18,1,2,4,4,3,5,1],"msg":"removed","ts":"20250704-180218"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_obdulu","F1","f33"],"etc":[18,1,2,4,4,2,6,1],"msg":"rec-windup-call","ts":"20250704-180218"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_obdulu","RWE-34"],"etc":[18,1,2,4,4,2,7],"msg":"rec-windup-end","ts":"20250704-180218"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_dodism","ejacti","26"],"etc":[18,1,2,4,4,3,6,1],"msg":"removed","ts":"20250704-180219"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_dodism","fopicu","26"],"etc":[18,1,2,4,4,3,7,1],"msg":"removed","ts":"20250704-180219"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_dodism","F1","f33"],"etc":[18,1,2,4,4,3,8,1],"msg":"rec-windup-call","ts":"20250704-180219"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_dodism","F1","RWB-39"],"etc":[18,1,2,4,4,3,8,2],"msg":"rec-windup-begin","ts":"20250704-180219"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_dodism","F1","F1","S_rastan","RWB-33"],"etc":[18,1,2,4,4,3,8,3,1,1],"msg":"rec-windup-begin","ts":"20250704-180219"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_dodism","F1","F1","S_rastan","uphili","26"],"etc":[18,1,2,4,4,3,8,3,1,2,1],"msg":"removed","ts":"20250704-180219"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_dodism","F1","F1","S_rastan","aslast","26"],"etc":[18,1,2,4,4,3,8,3,1,3,1],"msg":"removed","ts":"20250704-180220"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_dodism","F1","F1","S_rastan","RWE-34"],"etc":[18,1,2,4,4,3,8,3,1,4],"msg":"rec-windup-end","ts":"20250704-180220"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_dodism","F1","RWE-40"],"etc":[18,1,2,4,4,3,8,4],"msg":"rec-windup-end","ts":"20250704-180220"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","S_dodism","RWE-34"],"etc":[18,1,2,4,4,3,9],"msg":"rec-windup-end","ts":"20250704-180220"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","R2","RWE-185"],"etc":[18,1,2,4,4,4],"msg":"rec-windup-end","ts":"20250704-180220"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","S_uromer","RWE-34"],"etc":[18,1,2,4,5],"msg":"rec-windup-end","ts":"20250704-180220"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","R1","RWE-185"],"etc":[18,1,2,5],"msg":"rec-windup-end","ts":"20250704-180220"}
			{"level":"INFO","tags":["recursive-wind-up--failure","A1","RWE-34"],"etc":[18,1,3],"msg":"rec-windup-end","ts":"20250704-180220"}
			{"level":"INFO","tags":["344"],"etc":[19],"msg":"FINISHED BY WATCHERS, TIMED OUT","ts":"20250704-180220"}
			{
				"Components": {
					"System01": 0,
					"System02": 0,
					"System03": 0,
					"aslast": 0,
					"itchot": 0,
					"ontras": 0,
					"uphili": 0
				}
			}
			OK, all components are DOWN
			--- PASS: Test_DI_Container_Example2_01 (23.01s)
			PASS
			ok      github.com/rusriver/nutz/up/up/example  23.015s
	*/
}

func DI_Container_Example2_Main(t *testing.T, mode int) {

	/*
		We will model this (the pseudocode) (see also the diagram picture):

			A1 otdg is {
				a1_nguloo := a1_nguloo()
				a1_fakint := a1_fakint()
				a1_mestic := a1_mestic(a1_nguloo, a1_fakint, sparip, robles, histil, petorn, seigul)
			} started by R1 reactor

			R1 depends on subsystems {
				S_myinei == up
				S_nutari == up
				S_uromer == up_or_ceded
			} starts A1

			S_myinei otdg is {
				sostie := new_sostie()
				razinn := new_razinn(sostie)
				seigul := new_seigul(sostie)
				kaffin := new_kaffin(sostie, seigul)
				sparip := new_sparip(sostie, kaffin)
				send R1 the signal we are up
			} started by F1

			S_nutari otdg is {
				robles := new_robles(rebrei, teashi)
				send R1 the signal we are up
			} started by R3 reactor

			S_uromer otdg is {
				petorn := new_petorn()
				histil := new_histil(petorn, itchot, ontras)
				send R1 the signal we are up
			} started by R2 reactor

			R2 depends on subsystems {
				S_obdulu == up
				S_dodism == up
			} starts S_uromer

			R3 depends on subsystems {
				S_armati == up
				S_xylane == up
			} starts S_nutari

			S_obdulu otdg is {
				unborr := new_unborr(aslast, uphili)
				axseps := new_axseps(unborr)
				asdeed := new_asdeed(unborr, axseps)
				itchot := new_itchot(unborr, axseps, asdeed)
				send R2 the signal we are up
			} started by F1

			S_dodism otdg is {
				ontras := new_ontras()
				vexati := new_vexati()
				fopicu := new_fopicu()
				ejacti := new_ejacti()
				hoorgo := new_hoorgo()
				volath := new_volath()
				send R2 the signal we are up
			} started by F1

			S_armati otdg is {
				fatoca := new_fatoca()
				paudat := new_paudat()
				teashi := new_teashi(fatoca, paudat, aslast, uphili)
				send R3 the signal we are up
			} started by F1

			S_xylane otdg is {
				requac := new_requac()
				rebrei := new_rebrei(requac, aslast, uphili)
				send R3 the signal we are up
			} started by F1

			F1 fork is {
				-- start in parallel
				S_obdulu
				S_dodism
				S_armati
				S_xylane
				S_myinei
			} started by S_rastan

			S_rastan otdg is {
				aslast := new_aslast()
				uphili := new_uphili()
				start F1
			} started by ENTRY POINT


		Note that because this model has no diamond topology, it is a DAG, we can safely use
		the RecursiveWindUp() to stop it and remove.
	*/

	/*
		Few words on how to draw it:
			1) draw reactors, forks, otdgs, connect them
			2) draw internals of each otdg in separation
			3) connect cross-otdg dependencies between components
	*/

	/*
		In code, it all starts with creation of reactors (joins) and forks and otdg
		objects, here in main code area (or in infrastructure, if used on a distributed system).

		The same time, define the otdgs constuctors. This is done by defining a subroutines
		(functions). When you will call them, you will actually bring up their contents to life.
		But now you just define those subroutines.

		Note, that otdg is used as a side-object when you explicitly start your components,
		and when you wind up you use the method on the otdg object, it's kind of non-symmetrical,
		worth noting.

		Once you have reactors and forks and otdg objects created and otdg constructors defined,
		you are ready to start it all up. This is done by manually invoking one (or several)
		entry point otdgs. In this example, we have one, the S_rastan.

		The procedure of coding it is somehow mimicking the drawing procedure - first define
		top-level things, create cross-otdg links, the inside of otdgs are coded after that.
	*/

	// Create infra prerequisite components
	tracer := loggeterrzl.New(nil) // this is special thing, never wound-up, so outside of otdg
	tracer.LogEvent("BEGIN").Send()

	ctx := context.Background()
	wg1 := sync.WaitGroup{}
	wg1.Add(1)

	// define the variables for cross-otdgs.
	// though these are defined here at the top, you typically add them at last stages of coding.
	// you don't need everything here, only those that are used in otdgs other than where these are created.
	var (
		aslast *example.Aslast
		uphili *example.Uphili
		itchot *example.Itchot
		ontras *example.Ontras
		petorn *example.System03
		histil *example.System02
		robles *example.System01
		seigul *example.System03
		sparip *example.System02
	)

	// Create reactors and forks and otdg objects, the top-level topology
	// Start from the entry point, and build from that forward
	S_rastan := up.NewOTDG("S_rastan", tracer, nil)

	F1 := fork.NewFork("F1", tracer, S_rastan)

	S_myinei := up.NewOTDG("S_myinei", tracer, F1)
	S_armati := up.NewOTDG("S_armati", tracer, F1)
	S_xylane := up.NewOTDG("S_xylane", tracer, F1)
	S_obdulu := up.NewOTDG("S_obdulu", tracer, F1)
	S_dodism := up.NewOTDG("S_dodism", tracer, F1)

	R3 := reactor.NewReactor("R3", ctx, tracer)
	R3.NeedsPattern = map[string]uint8{
		S_armati.Name: reactor.RequiredStatus_Up,
		S_xylane.Name: reactor.RequiredStatus_Up,
	}
	R3.Func = func(r *reactor.Reactor, msg *reactor.Msg) {

	}

	S_nutari := up.NewOTDG("S_nutari", tracer, R3)

	R2 := reactor.NewReactor("R2", ctx, tracer)
	R2.NeedsPattern = map[string]uint8{
		S_obdulu.Name: reactor.RequiredStatus_Up,
		S_dodism.Name: reactor.RequiredStatus_Up,
	}

	S_uromer := up.NewOTDG("S_uromer", tracer, R2)

	R1 := reactor.NewReactor("R1", ctx, tracer)
	R1.NeedsPattern = map[string]uint8{
		S_myinei.Name: reactor.RequiredStatus_Up,
		S_nutari.Name: reactor.RequiredStatus_Up,
		S_uromer.Name: reactor.RequiredStatus_Up,
	}

	R2.Func_CanReact = func(r *reactor.Reactor, msg *reactor.Msg, status *reactor.Status) {
		petorn, histil = Make_S_uromer(S_uromer, ontras, itchot, func(otdg *up.OTDG) {
			R1.SendStatus_Up(otdg, otdg.Name)
		})
	}
	R3.Func_CanReact = func(r *reactor.Reactor, msg *reactor.Msg, status *reactor.Status) {
		Make_S_nutari(S_nutari, func(otdg *up.OTDG) {
			switch mode {
			case 0:
				R1.SendStatus_Up(otdg, otdg.Name)
			case 1:
				R1.SendStatus_Fail(otdg, otdg.Name)
			}
		})
	}
	A1 := up.NewOTDG("A1", tracer, R1)
	R1.Func_CanReact = func(r *reactor.Reactor, msg *reactor.Msg, status *reactor.Status) {
		Make_A1(A1, seigul, sparip, robles, petorn, histil, func(otdg *up.OTDG) {
			wg1.Done()
		})
	}

	// create watchers - you can add them the very last in the coding process
	var finished_by_watchers bool
	W01 := &watcher.Watcher{
		Name:        "W01",
		Tracer:      tracer.NSTIC("W01").ITracer(),
		Context:     ctx,
		Timeout:     time.Second * 15,
		TimeoutFunc: func(w *watcher.Watcher) {},
	}
	W02 := &watcher.Watcher{
		Name:        "W02",
		Tracer:      tracer.NSTIC("W02").ITracer(),
		Context:     ctx,
		Timeout:     time.Second * 15,
		TimeoutFunc: func(w *watcher.Watcher) {},
	}
	W03 := &watcher.Watcher{
		Name:    "W03",
		Tracer:  tracer.NSTIC("W03").ITracer(),
		Context: ctx,
		Timeout: time.Second * 20,
		TimeoutFunc: func(w *watcher.Watcher) {
			A1.RecursiveWindUp(tracer.NSTIC("recursive-wind-up--failure").ITracer())
			finished_by_watchers = true
			wg1.Done()
		},
	}
	W01.Start()
	W02.Start()
	W03.Start()

	// reactors need explicit start, don't forget that
	R1.StartFresh()
	R2.StartFresh()
	R3.StartFresh()

	// the entry point to start the whole thing up
	// make the otdand initiate the whole thing start up process
	aslast, uphili = Make_S_rastan(S_rastan, func(otdg *up.OTDG) {
		F1.RunFuncInParallel(func() {
			seigul, sparip = Make_S_myinei(S_myinei, func(otdg *up.OTDG) {
				R1.SendStatus_Up(otdg, otdg.Name)
			})
		})
		F1.RunFuncInParallel(func() {
			Make_S_armati(S_armati, aslast, uphili, func(otdg *up.OTDG) {
				R3.SendStatus_Up(otdg, otdg.Name)
			})
		})
		F1.RunFuncInParallel(func() {
			Make_S_xylane(S_xylane, aslast, uphili, func(otdg *up.OTDG) {
				R3.SendStatus_Up(otdg, otdg.Name)
			})
		})
		F1.RunFuncInParallel(func() {
			itchot = Make_S_obdulu(S_obdulu, aslast, uphili, func(otdg *up.OTDG) {
				R2.SendStatus_Up(otdg, otdg.Name)
			})
		})
		F1.RunFuncInParallel(func() {
			ontras = Make_S_dodism(S_dodism, func(otdg *up.OTDG) {
				R2.SendStatus_Up(otdg, otdg.Name)
			})
		})
	})

	time.Sleep(time.Second)
	wg1.Wait()
	switch mode {
	case 0:
		require.Equal(t, false, finished_by_watchers)
	case 1:
		require.Equal(t, true, finished_by_watchers)
	}
	if finished_by_watchers {
		tracer.LogEvent("344").SendMsg("FINISHED BY WATCHERS, TIMED OUT")
		bb, _ := json.MarshalIndent(example.ControlRegistry, "", "    ")
		fmt.Println(string(bb))
		require.Equal(t, true, example.ControlRegistry.MatchThePattern(map[string]int{
			"System01": 0,
			"System02": 0,
			"System03": 0,
			"aslast":   0,
			"itchot":   0,
			"ontras":   0,
			"uphili":   0,
		}))
		fmt.Println("OK, all components are DOWN")
		return
	}

	W01.Stop_Done()
	W02.Stop_Done()
	W03.Stop_Done()

	time.Sleep(time.Second)
	tracer.LogEvent("UP").Send()

	// make sure everything is started
	bb, _ := json.MarshalIndent(example.ControlRegistry, "", "    ")
	fmt.Println(string(bb))
	require.Equal(t, true, example.ControlRegistry.MatchThePattern(map[string]int{
		"System01": 12,
		"System02": 7,
		"System03": 5,
		"aslast":   1,
		"itchot":   1,
		"ontras":   1,
		"uphili":   1,
	}))
	var (
		_types int
		_total int
	)
	for _, v := range example.ControlRegistry.Components {
		_types++
		_total += v
	}
	fmt.Printf("OK, %v components of %v types are UP\n", _total, _types)

	time.Sleep(5 * time.Second)

	// we're going to stop then
	tracer.LogEvent("GO-DOWN").Send()

	// S_rastan.WindUp(tracer.NSTIC("wind-up").ITracer())

	A1.RecursiveWindUp(tracer.NSTIC("A1-recursive-wind-up").ITracer())

	// make sure everything is wind-up
	bb, _ = json.MarshalIndent(example.ControlRegistry, "", "    ")
	fmt.Println(string(bb))
	require.Equal(t, true, example.ControlRegistry.MatchThePattern(map[string]int{
		"System01": 0,
		"System02": 0,
		"System03": 0,
		"aslast":   0,
		"itchot":   0,
		"ontras":   0,
		"uphili":   0,
	}))
	fmt.Println("OK, all components are DOWN")

	tracer.LogEvent("END").Send()
}

// in real projects, this better be in separate file
// note that all these functions have idiomatic (but not exactly the same) signature
func Make_S_rastan(
	otdg *up.OTDG,
	final_func func(otdg *up.OTDG),
) (
	aslast *example.Aslast,
	uphili *example.Uphili,
) {
	otdg.Tracer.LogEvent("288").SendMsg("starting")
	controlflow.Try(func() (err any) {
		RandomSleep()

		// syntax example, variant 1
		aslast = up.Up(otdg, "aslast", func(name string) (*example.Aslast, upcommon.IRemovable) {
			s := example.NewAslast(name, otdg.Tracer.NSTIC(name).ITracer())
			return s, s
		})

		// syntax example, variant 2
		uphili = otdg.Up("uphili", func(name string) (any, upcommon.IRemovable) {
			s := example.NewUphili(name, otdg.Tracer.NSTIC(name).ITracer())
			return s, s
		}).(*example.Uphili)

		otdg.Tracer.LogEvent("389").SendMsg("done")
		SmallSleep()
		if final_func != nil {
			final_func(otdg)
		}
		return nil
	}).Catch(func(e *controlflow.Event) {
		otdg.Tracer.EventEvent("E320", e).Send()
	})
	return
}

// in real projects, this better be in separate file
// note that all these functions have idiomatic (but not exactly the same) signature
func Make_S_myinei(
	otdg *up.OTDG,
	final_func func(otdg *up.OTDG),
) (
	seigul *example.System03,
	sparip *example.System02,
) {
	otdg.Tracer.LogEvent("288").SendMsg("starting")
	controlflow.Try(func() (err any) {
		RandomSleep()

		sostie := up.Up(otdg, "sostie", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name, otdg.Tracer.NSTIC(name).ITracer())
			return s, s
		})

		razinn := up.Up(otdg, "razinn", func(name string) (*example.System02, upcommon.IRemovable) {
			s := example.NewSystem02(name, otdg.Tracer.NSTIC(name).ITracer(), sostie)
			return s, s
		})
		_ = razinn

		seigul = up.Up(otdg, "seigul", func(name string) (*example.System03, upcommon.IRemovable) {
			s := example.NewSystem03(name, otdg.Tracer.NSTIC(name).ITracer(), sostie)
			return s, s
		})

		kaffin := up.Up(otdg, "kaffin", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name, otdg.Tracer.NSTIC(name).ITracer(), sostie, seigul)
			return s, s
		})

		sparip = up.Up(otdg, "sparip", func(name string) (*example.System02, upcommon.IRemovable) {
			s := example.NewSystem02(name, otdg.Tracer.NSTIC(name).ITracer(), sostie, kaffin)
			return s, s
		})

		otdg.Tracer.LogEvent("389").SendMsg("done")
		SmallSleep()
		if final_func != nil {
			final_func(otdg)
		}
		return nil
	}).Catch(func(e *controlflow.Event) {
		otdg.Tracer.EventEvent("E375", e).Send()
	})
	return
}

// in real projects, this better be in separate file
// note that all these functions have idiomatic (but not exactly the same) signature
func Make_S_armati(
	otdg *up.OTDG,
	aslast *example.Aslast,
	uphili *example.Uphili,
	final_func func(otdg *up.OTDG),
) {
	otdg.Tracer.LogEvent("288").SendMsg("starting")
	controlflow.Try(func() (err any) {
		RandomSleep()

		fatoca := up.Up(otdg, "fatoca", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name, otdg.Tracer.NSTIC(name).ITracer())
			return s, s
		})

		paudat := up.Up(otdg, "paudat", func(name string) (*example.System02, upcommon.IRemovable) {
			s := example.NewSystem02(name, otdg.Tracer.NSTIC(name).ITracer())
			return s, s
		})

		teashi := up.Up(otdg, "teashi", func(name string) (*example.System03, upcommon.IRemovable) {
			s := example.NewSystem03(name, otdg.Tracer.NSTIC(name).ITracer(), aslast, uphili, fatoca, paudat)
			return s, s
		})
		_ = teashi

		otdg.Tracer.LogEvent("389").SendMsg("done")
		SmallSleep()
		if final_func != nil {
			final_func(otdg)
		}
		return nil
	}).Catch(func(e *controlflow.Event) {
		otdg.Tracer.EventEvent("E399", e).Send()
	})
}

// in real projects, this better be in separate file
// note that all these functions have idiomatic (but not exactly the same) signature
func Make_S_xylane(
	otdg *up.OTDG,
	aslast *example.Aslast,
	uphili *example.Uphili,
	final_func func(otdg *up.OTDG),
) {
	otdg.Tracer.LogEvent("288").SendMsg("starting")
	controlflow.Try(func() (err any) {
		RandomSleep()

		requac := up.Up(otdg, "requac", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name, otdg.Tracer.NSTIC(name).ITracer())
			return s, s
		})

		rebrei := up.Up(otdg, "rebrei", func(name string) (*example.System02, upcommon.IRemovable) {
			s := example.NewSystem02(name, otdg.Tracer.NSTIC(name).ITracer(), requac, aslast, uphili)
			return s, s
		})
		_ = rebrei

		otdg.Tracer.LogEvent("389").SendMsg("done")
		SmallSleep()
		if final_func != nil {
			final_func(otdg)
		}
		return nil
	}).Catch(func(e *controlflow.Event) {
		otdg.Tracer.EventEvent("E422", e).Send()
	})
}

// in real projects, this better be in separate file
// note that all these functions have idiomatic (but not exactly the same) signature
func Make_S_obdulu(
	otdg *up.OTDG,
	aslast *example.Aslast,
	uphili *example.Uphili,
	final_func func(otdg *up.OTDG),
) (
	itchot *example.Itchot,
) {
	otdg.Tracer.LogEvent("288").SendMsg("starting")
	controlflow.Try(func() (err any) {
		RandomSleep()

		unborr := up.Up(otdg, "unborr", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name, otdg.Tracer.NSTIC(name).ITracer(), aslast, uphili)
			return s, s
		})

		axseps := up.Up(otdg, "axseps", func(name string) (*example.System02, upcommon.IRemovable) {
			s := example.NewSystem02(name, otdg.Tracer.NSTIC(name).ITracer(), unborr)
			return s, s
		})

		asdeed := up.Up(otdg, "asdeed", func(name string) (*example.System03, upcommon.IRemovable) {
			s := example.NewSystem03(name, otdg.Tracer.NSTIC(name).ITracer(), axseps, unborr)
			return s, s
		})
		_ = asdeed

		itchot = up.Up(otdg, "itchot", func(name string) (*example.Itchot, upcommon.IRemovable) {
			s := example.NewItchot(name, otdg.Tracer.NSTIC(name).ITracer(), unborr, axseps)
			return s, s
		})

		otdg.Tracer.LogEvent("389").SendMsg("done")
		SmallSleep()
		if final_func != nil {
			final_func(otdg)
		}
		return nil
	}).Catch(func(e *controlflow.Event) {
		otdg.Tracer.EventEvent("E447", e).Send()
	})
	return
}

// in real projects, this better be in separate file
// note that all these functions have idiomatic (but not exactly the same) signature
func Make_S_dodism(
	otdg *up.OTDG,
	final_func func(otdg *up.OTDG),
) (
	ontras *example.Ontras,
) {
	otdg.Tracer.LogEvent("288").SendMsg("starting")
	controlflow.Try(func() (err any) {
		RandomSleep()

		_ = up.Up(otdg, "fopicu", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name, otdg.Tracer.NSTIC(name).ITracer())
			return s, s
		})

		_ = up.Up(otdg, "ejacti", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name, otdg.Tracer.NSTIC(name).ITracer())
			return s, s
		})

		_ = up.Up(otdg, "hoorgo", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name, otdg.Tracer.NSTIC(name).ITracer())
			return s, s
		})

		_ = up.Up(otdg, "volath", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name, otdg.Tracer.NSTIC(name).ITracer())
			return s, s
		})

		_ = up.Up(otdg, "vexati", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name, otdg.Tracer.NSTIC(name).ITracer())
			return s, s
		})

		ontras = up.Up(otdg, "ontras", func(name string) (*example.Ontras, upcommon.IRemovable) {
			s := example.NewOntras(name, otdg.Tracer.NSTIC(name).ITracer())
			return s, s
		})

		otdg.Tracer.LogEvent("389").SendMsg("done")
		SmallSleep()
		if final_func != nil {
			final_func(otdg)
		}
		return nil
	}).Catch(func(e *controlflow.Event) {
		otdg.Tracer.EventEvent("E471", e).Send()
	})
	return
}

// in real projects, this better be in separate file
// note that all these functions have idiomatic (but not exactly the same) signature
func Make_S_nutari(
	otdg *up.OTDG,
	final_func func(otdg *up.OTDG),
) (
	robles *example.System01,
) {
	otdg.Tracer.LogEvent("288").SendMsg("starting")
	controlflow.Try(func() (err any) {
		RandomSleep()

		robles = up.Up(otdg, "robles", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name, otdg.Tracer.NSTIC(name).ITracer())
			return s, s
		})

		otdg.Tracer.LogEvent("389").SendMsg("done")
		SmallSleep()
		if final_func != nil {
			final_func(otdg)
		}
		return nil
	}).Catch(func(e *controlflow.Event) {
		otdg.Tracer.EventEvent("E495", e).Send()
	})
	return
}

// in real projects, this better be in separate file
// note that all these functions have idiomatic (but not exactly the same) signature
func Make_S_uromer(
	otdg *up.OTDG,
	ontras *example.Ontras,
	itchot *example.Itchot,
	final_func func(otdg *up.OTDG),
) (
	petorn *example.System03,
	histil *example.System02,
) {
	otdg.Tracer.LogEvent("288").SendMsg("starting")
	controlflow.Try(func() (err any) {
		RandomSleep()

		petorn = up.Up(otdg, "petorn", func(name string) (*example.System03, upcommon.IRemovable) {
			s := example.NewSystem03(name, otdg.Tracer.NSTIC(name).ITracer())
			return s, s
		})

		histil = up.Up(otdg, "histil", func(name string) (*example.System02, upcommon.IRemovable) {
			s := example.NewSystem02(name, otdg.Tracer.NSTIC(name).ITracer(), petorn, itchot, ontras)
			return s, s
		})

		otdg.Tracer.LogEvent("389").SendMsg("done")
		SmallSleep()
		if final_func != nil {
			final_func(otdg)
		}
		return nil
	}).Catch(func(e *controlflow.Event) {
		otdg.Tracer.EventEvent("E522", e).Send()
	})
	return
}

// in real projects, this better be in separate file
// note that all these functions have idiomatic (but not exactly the same) signature
func Make_A1(
	otdg *up.OTDG,
	seigul *example.System03,
	sparip *example.System02,
	robles *example.System01,
	petorn *example.System03,
	histil *example.System02,
	final_func func(otdg *up.OTDG),
) {
	otdg.Tracer.LogEvent("288").SendMsg("starting")
	controlflow.Try(func() (err any) {
		RandomSleep()

		// syntax example, variant 1
		a1_nguloo := up.Up(otdg, "a1_nguloo", func(name string) (*example.System01, upcommon.IRemovable) {
			s := example.NewSystem01(name, otdg.Tracer.NSTIC(name).ITracer())
			return s, s
		})

		// syntax example, variant 2
		a1_fakint := otdg.Up("a1_fakint", func(name string) (any, upcommon.IRemovable) {
			s := example.NewSystem02(name, otdg.Tracer.NSTIC(name).ITracer())
			return s, s
		}).(*example.System02)

		up.Up(otdg, "a1_mestic", func(name string) (*example.System03, upcommon.IRemovable) {
			s := example.NewSystem03(name, otdg.Tracer.NSTIC(name).ITracer(), a1_nguloo, a1_fakint, seigul, sparip, robles, petorn, histil)
			return s, s
		})

		otdg.Tracer.LogEvent("389").SendMsg("done")
		SmallSleep()
		if final_func != nil {
			final_func(otdg)
		}
		return nil
	}).Catch(func(e *controlflow.Event) {
		otdg.Tracer.EventEvent("E566", e).Send()
	})
}

func RandomSleep() {
	time.Sleep(time.Duration(rand.IntN(1500))*time.Millisecond + 300*time.Millisecond)
}

func SmallSleep() {
	time.Sleep(300 * time.Millisecond)
}
