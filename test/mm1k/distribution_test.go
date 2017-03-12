package mm1k_test

import (
	"fmt"
	"mm1k"
	"testing"
)

func TestExpDistribution(t *testing.T) {
	e := mm1k.NewExpDistribution(0.5, 42)
	e.Get()
}

func ExampleExpDistribution() {
	e := mm1k.NewExpDistribution(0.5, 42)
	for i := 0; i < 10; i++ {
		fmt.Println(e.Get())
	}
	// Output:
	// 0.9914768298047957
	// 0.26109436751792503
	// 0.3064669044181516
	// 0.6768927252188972
	// 0.23192847758808366
	// 2.111317879104961
	// 1.7180304580523118
	// 0.29726736216152605
	// 2.7956095955449896
	// 2.8525479779829563
}

func ExampleExpDistribution2() {
	e := mm1k.NewExpDistribution(0.1, 42)
	for i := 0; i < 10; i++ {
		fmt.Println(e.Get())
	}
	// Output:
	// 4.957384149023978
	// 1.305471837589625
	// 1.532334522090758
	// 3.384463626094486
	// 1.1596423879404183
	// 10.556589395524805
	// 8.590152290261559
	// 1.4863368108076302
	// 13.978047977724946
	// 14.262739889914782

}
