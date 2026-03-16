package gamedata

/*
// 1. Update signatures to 7.72; DAT = 439D5A33, SPR = 439852BE (little endian, remember to reverse).
// 2. Load in ObjectBuilder (https://github.com/ottools/ObjectBuilder).
// 3. Export as 8.54 v3.
// 4. Run:

	import "badc0de.net/pkg/go-tibia/dat"

	func main() {
		f, _ := os.Open("Tibiantis.dat")
		defer f.Close()
		ds, _ := dat.NewDataset(f)
		var fluidContainers, fluids, stackable []int
		for i := ds.MinItemID(); i <= ds.MaxItemID(); i++ {
			item := ds.Item(i)

			if item.FluidContainer {
				fluidContainers = append(fluidContainers, int(i))
			}
			if item.Splash {
				fluids = append(fluids, int(i))
			}
			if item.IsStackable {
				stackable = append(stackable, int(i))
			}
		}

		fmt.Fprintf(os.Stderr, "sFluidContainers = %#v\n", fluidContainers)
		fmt.Fprintf(os.Stderr, "sFluids = %#v\n", fluids)
		fmt.Fprintf(os.Stderr, "sStackable = %#v\n", stackable)
	}
*/
var (
	IsFluidContainer = map[int]bool{}
	IsFluid          = map[int]bool{}
	IsStackable       = map[int]bool{}

	sFluidContainers = []int{
		2524, 2873, 2874, 2875, 2876, 2877, 2879, 2880, 2881, 2882,
		2885, 2893, 2901, 2902, 2903, 2904, 3465, 3477, 3478, 3479,
		3480,
	}
	sFluids = []int{
		2886, 2887, 2888, 2889, 2890, 2891, 2895, 2896, 2897, 2898,
		2899, 2900,
	}
	sStackable = []int{
		1781, 2784, 2992, 3026, 3027, 3028, 3029, 3030, 3031, 3032,
		3033, 3034, 3035, 3040, 3042, 3043, 3114, 3145, 3146, 3207,
		3215, 3250, 3277, 3287, 3298, 3445, 3446, 3447, 3448, 3449,
		3450, 3492, 3533, 3534, 3547, 3548, 3560, 3577, 3578, 3579,
		3580, 3581, 3582, 3583, 3584, 3585, 3586, 3587, 3588, 3589,
		3590, 3591, 3595, 3596, 3597, 3598, 3599, 3600, 3601, 3602,
		3603, 3604, 3605, 3606, 3721, 3722, 3723, 3724, 3725, 3726,
		3727, 3728, 3729, 3730, 3731, 3732, 3734, 3735, 3736, 3737,
		3738, 3739, 3740, 3741, 5021,
	}
)

func init() {
	for _, x := range sFluidContainers {
		IsFluidContainer[x] = true
	}
	for _, x := range sFluids {
		IsFluid[x] = true
	}
	for _, x := range sStackable {
		IsStackable[x] = true
	}
}
