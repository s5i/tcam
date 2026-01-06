package enum

import "fmt"

type DatCategory uint8

func (c DatCategory) String() string {
	if str, ok := datCategoryMap[c]; ok {
		return str
	}
	return fmt.Sprintf("Unknown-%d", c)
}

const (
	DatCategoryItem DatCategory = iota
	DatCategoryCreature
	DatCategoryEffect
	DatCategoryMissile
	DatCategoryInvalid
	DatCategoryLast = DatCategoryInvalid
)

var datCategoryMap = map[DatCategory]string{
	DatCategoryItem:     "Item",
	DatCategoryCreature: "Creature",
	DatCategoryEffect:   "Effect",
	DatCategoryMissile:  "Missile",
	DatCategoryInvalid:  "Invalid",
}

type DatAttribute uint8

func (c DatAttribute) String() string {
	if str, ok := datAttributeMap[c]; ok {
		return fmt.Sprintf("%s-%d", str, c)
	}
	return fmt.Sprintf("Unknown-%d", c)
}

const (
	DatAttributeGround          DatAttribute = 0
	DatAttributeGroundBorder    DatAttribute = 1
	DatAttributeOnBottom        DatAttribute = 2
	DatAttributeOnTop           DatAttribute = 3
	DatAttributeContainer       DatAttribute = 4
	DatAttributeStackable       DatAttribute = 5
	DatAttributeForceUse        DatAttribute = 6
	DatAttributeMultiUse        DatAttribute = 7
	DatAttributeWritable        DatAttribute = 8
	DatAttributeWritableOnce    DatAttribute = 9
	DatAttributeFluidContainer  DatAttribute = 10
	DatAttributeSplash          DatAttribute = 11
	DatAttributeNotWalkable     DatAttribute = 12
	DatAttributeNotMoveable     DatAttribute = 13
	DatAttributeBlockProjectile DatAttribute = 14
	DatAttributeNotPathable     DatAttribute = 15
	DatAttributePickupable      DatAttribute = 16
	DatAttributeHangable        DatAttribute = 17
	DatAttributeHookSouth       DatAttribute = 18
	DatAttributeHookEast        DatAttribute = 19
	DatAttributeRotateable      DatAttribute = 20
	DatAttributeLight           DatAttribute = 21
	DatAttributeDontHide        DatAttribute = 22
	DatAttributeFloorChange     DatAttribute = 23
	DatAttributeDisplacement    DatAttribute = 24
	DatAttributeElevation       DatAttribute = 25
	DatAttributeLyingCorpse     DatAttribute = 26
	DatAttributeAnimateAlways   DatAttribute = 27
	DatAttributeMinimapColor    DatAttribute = 28
	DatAttributeLensHelp        DatAttribute = 29
	DatAttributeFullGround      DatAttribute = 30
	DatAttributeLook            DatAttribute = 31
	DatAttributeCloth           DatAttribute = 32
	DatAttributeMarket          DatAttribute = 33
	DatAttributeUsable          DatAttribute = 34
	DatAttributeWrapable        DatAttribute = 35
	DatAttributeUnwrapable      DatAttribute = 36
	DatAttributeTopEffect       DatAttribute = 37
	DatAttributeBones           DatAttribute = 38
	DatAttributeOpacity         DatAttribute = 100
	DatAttributeNotPreWalkable  DatAttribute = 101
	DatAttributeNoMoveAnimation DatAttribute = 253
	DatAttributeChargeable      DatAttribute = 254
	DatAttributeInvalid         DatAttribute = 255

	DatAttributeLast = DatAttributeInvalid
)

var datAttributeMap = map[DatAttribute]string{
	DatAttributeGround:          "Ground",
	DatAttributeGroundBorder:    "GroundBorder",
	DatAttributeOnBottom:        "OnBottom",
	DatAttributeOnTop:           "OnTop",
	DatAttributeContainer:       "Container",
	DatAttributeStackable:       "Stackable",
	DatAttributeForceUse:        "ForceUse",
	DatAttributeMultiUse:        "MultiUse",
	DatAttributeWritable:        "Writable",
	DatAttributeWritableOnce:    "WritableOnce",
	DatAttributeFluidContainer:  "FluidContainer",
	DatAttributeSplash:          "Splash",
	DatAttributeNotWalkable:     "NotWalkable",
	DatAttributeNotMoveable:     "NotMoveable",
	DatAttributeBlockProjectile: "BlockProjectile",
	DatAttributeNotPathable:     "NotPathable",
	DatAttributePickupable:      "Pickupable",
	DatAttributeHangable:        "Hangable",
	DatAttributeHookSouth:       "HookSouth",
	DatAttributeHookEast:        "HookEast",
	DatAttributeRotateable:      "Rotateable",
	DatAttributeLight:           "Light",
	DatAttributeDontHide:        "DontHide",
	DatAttributeDisplacement:    "Displacement",
	DatAttributeElevation:       "Elevation",
	DatAttributeLyingCorpse:     "LyingCorpse",
	DatAttributeAnimateAlways:   "AnimateAlways",
	DatAttributeMinimapColor:    "MinimapColor",
	DatAttributeLensHelp:        "LensHelp",
	DatAttributeFullGround:      "FullGround",
	DatAttributeLook:            "Look",
	DatAttributeCloth:           "Cloth",
	DatAttributeMarket:          "Market",
	DatAttributeUsable:          "Usable",
	DatAttributeWrapable:        "Wrapable",
	DatAttributeUnwrapable:      "Unwrapable",
	DatAttributeTopEffect:       "TopEffect",
	DatAttributeBones:           "Bones",
	DatAttributeOpacity:         "Opacity",
	DatAttributeNotPreWalkable:  "NotPreWalkable",
	DatAttributeFloorChange:     "FloorChange",
	DatAttributeNoMoveAnimation: "NoMoveAnimation",
	DatAttributeChargeable:      "Chargeable",
	DatAttributeInvalid:         "Invalid",
}
