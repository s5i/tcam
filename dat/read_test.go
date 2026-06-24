package dat_test

import (
	"os"
	"testing"

	"github.com/s5i/tcam/dat"
)

var (
	tibiantisFluidContainers = []int{
		2524, 2873, 2874, 2875, 2876, 2877, 2879, 2880, 2881, 2882,
		2885, 2893, 2901, 2902, 2903, 2904, 3465, 3477, 3478, 3479,
		3480,
	}
	tibiantisFluids = []int{
		2886, 2887, 2888, 2889, 2890, 2891, 2895, 2896, 2897, 2898,
		2899, 2900,
	}
	tibiantisStackable = []int{
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

	tibiaRelicFluidContainers = []int{
		2524, 2873, 2874, 2875, 2876, 2877, 2879, 2880, 2881, 2882,
		2885, 2893, 2901, 2902, 2903, 2904, 3465, 3477, 3478, 3479,
		3480,
	}
	tibiaRelicFluids = []int{
		2886, 2887, 2888, 2889, 2890, 2891, 2895, 2896, 2897, 2898,
		2899, 2900,
	}
	tibiaRelicStackable = []int{
		1781, 2992, 3026, 3027, 3028, 3029, 3030, 3031, 3032, 3033,
		3034, 3035, 3040, 3042, 3043, 3044, 3114, 3207, 3215, 3239,
		3250, 3277, 3287, 3298, 3347, 3446, 3447, 3448, 3449, 3450,
		3492, 3533, 3534, 3547, 3548, 3577, 3578, 3579, 3580, 3581,
		3582, 3583, 3584, 3585, 3586, 3587, 3588, 3589, 3590, 3591,
		3595, 3596, 3597, 3598, 3599, 3600, 3601, 3602, 3603, 3604,
		3605, 3606, 3723, 3724, 3725, 3726, 3727, 3728, 3729, 3730,
		3731, 3732, 3734, 3735, 3736, 3737, 3738, 3739, 3740, 3741,
		4827, 4838, 5021,
	}
)

func TestRead_Tibiantis(t *testing.T) {
	file := readFile(t, "testdata/Tibiantis.dat")

	if file.Signature != 0x6970EFAD {
		t.Fatalf("Signature = 0x%X, want 0x6970EFAD", file.Signature)
	}
	if file.ItemCount != 5089 {
		t.Fatalf("ItemCount = %d, want 5089", file.ItemCount)
	}
	if file.OutfitCount != 254 {
		t.Fatalf("OutfitCount = %d, want 254", file.OutfitCount)
	}
	if file.EffectCount != 26 {
		t.Fatalf("EffectCount = %d, want 26", file.EffectCount)
	}
	if file.MissileCount != 16 {
		t.Fatalf("MissileCount = %d, want 16", file.MissileCount)
	}

	assertIDs(t, file.ItemCount, "IsFluidContainer", file.IsFluidContainer, tibiantisFluidContainers)
	assertIDs(t, file.ItemCount, "IsFluid", file.IsFluid, tibiantisFluids)
	assertIDs(t, file.ItemCount, "IsStackable", file.IsStackable, tibiantisStackable)

	if file.IsFluidContainer(3031) {
		t.Fatal("gold coin (3031) should not be a fluid container")
	}
	if !file.IsStackable(3031) {
		t.Fatal("gold coin (3031) should be stackable")
	}
	if file.IsFluid(3031) {
		t.Fatal("gold coin (3031) should not be a fluid")
	}
}

func TestRead_TibiaRelic(t *testing.T) {
	file := readFile(t, "testdata/TibiaRelic.dat")

	if file.Signature != 0x439D5A33 {
		t.Fatalf("Signature = 0x%X, want 0x439D5A33", file.Signature)
	}
	if file.ItemCount != 5161 {
		t.Fatalf("ItemCount = %d, want 5161", file.ItemCount)
	}
	if file.OutfitCount != 254 {
		t.Fatalf("OutfitCount = %d, want 254", file.OutfitCount)
	}
	if file.EffectCount != 26 {
		t.Fatalf("EffectCount = %d, want 26", file.EffectCount)
	}
	if file.MissileCount != 16 {
		t.Fatalf("MissileCount = %d, want 16", file.MissileCount)
	}

	assertIDs(t, file.ItemCount, "IsFluidContainer", file.IsFluidContainer, tibiaRelicFluidContainers)
	assertIDs(t, file.ItemCount, "IsFluid", file.IsFluid, tibiaRelicFluids)
	assertIDs(t, file.ItemCount, "IsStackable", file.IsStackable, tibiaRelicStackable)

	if file.IsFluidContainer(3031) {
		t.Fatal("gold coin (3031) should not be a fluid container")
	}
	if !file.IsStackable(3031) {
		t.Fatal("gold coin (3031) should be stackable")
	}
	if file.IsFluid(3031) {
		t.Fatal("gold coin (3031) should not be a fluid")
	}
}

func readFile(t *testing.T, path string) *dat.File {
	t.Helper()

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("Open(%q) error: %v", path, err)
	}
	defer f.Close()

	file, err := dat.Read(f)
	if err != nil {
		t.Fatalf("Read(%q) error: %v", path, err)
	}
	return file
}

func assertIDs(t *testing.T, itemCount uint16, name string, check func(int) bool, want []int) {
	t.Helper()
	for _, id := range want {
		if !check(id) {
			t.Errorf("%s(%d) = false, want true", name, id)
		}
	}
	for id := 100; id <= int(itemCount); id++ {
		wantTrue := false
		for _, w := range want {
			if w == id {
				wantTrue = true
				break
			}
		}
		if check(id) != wantTrue {
			t.Errorf("%s(%d) = %v, want %v", name, id, !wantTrue, wantTrue)
		}
	}
}
