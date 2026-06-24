package dat

import "fmt"

func readProperties(br *binReader, version formatVersion, id uint16) (Properties, error) {
	switch version {
	case formatV1:
		return readPropertiesV1(br, id)
	case formatV2:
		return readPropertiesV2(br, id)
	case formatV3:
		return readPropertiesV3(br, id)
	case formatV4:
		return readPropertiesV4(br, id)
	case formatV5:
		return readPropertiesV5(br, id)
	case formatV6:
		return readPropertiesV6(br, id)
	default:
		return Properties{}, fmt.Errorf("dat: unsupported format version %d", version)
	}
}

func readPropertiesV1(br *binReader, id uint16) (Properties, error) {
	var p Properties
	for {
		flag, err := br.u8()
		if err != nil {
			return p, err
		}
		if flag == 0xFF {
			return p, nil
		}
		switch flag {
		case 0x00:
			if _, err := br.u16(); err != nil {
				return p, err
			}
			p.Ground = true
		case 0x01:
			p.OnBottom = true
		case 0x02:
			p.OnTop = true
		case 0x03:
			p.Container = true
		case 0x04:
			p.Stackable = true
		case 0x05:
			p.MultiUse = true
		case 0x06:
			p.ForceUse = true
		case 0x07, 0x08:
			if _, err := br.u16(); err != nil {
				return p, err
			}
		case 0x09:
			p.FluidContainer = true
		case 0x0A:
			p.Fluid = true
		case 0x0B:
			p.Unpassable = true
		case 0x0C:
			p.Unmoveable = true
		case 0x0D:
			p.BlockMissile = true
		case 0x0E:
			p.BlockPathfind = true
		case 0x0F:
			p.Pickupable = true
		case 0x10:
			if err := br.skip(4); err != nil {
				return p, err
			}
		case 0x11:
		case 0x12:
		case 0x13:
			if _, err := br.u16(); err != nil {
				return p, err
			}
		case 0x14:
		case 0x16:
			if _, err := br.u16(); err != nil {
				return p, err
			}
		case 0x17:
		case 0x18:
		case 0x19:
		case 0x1A:
			if _, err := br.u16(); err != nil {
				return p, err
			}
		default:
			return p, unknownFlag(flag, id)
		}
	}
}

func readPropertiesV2(br *binReader, id uint16) (Properties, error) {
	var p Properties
	for {
		flag, err := br.u8()
		if err != nil {
			return p, err
		}
		if flag == 0xFF {
			return p, nil
		}
		switch flag {
		case 0x00:
			if _, err := br.u16(); err != nil {
				return p, err
			}
			p.Ground = true
		case 0x01:
			p.OnBottom = true
		case 0x02:
			p.OnTop = true
		case 0x03:
			p.Container = true
		case 0x04:
			p.Stackable = true
		case 0x05:
			p.MultiUse = true
		case 0x06:
			p.ForceUse = true
		case 0x07, 0x08:
			if _, err := br.u16(); err != nil {
				return p, err
			}
		case 0x09:
			p.FluidContainer = true
		case 0x0A:
			p.Fluid = true
		case 0x0B:
			p.Unpassable = true
		case 0x0C:
			p.Unmoveable = true
		case 0x0D:
			p.BlockMissile = true
		case 0x0E:
			p.BlockPathfind = true
		case 0x0F:
			p.Pickupable = true
		case 0x10:
			if err := br.skip(4); err != nil {
				return p, err
			}
		case 0x11, 0x12, 0x14, 0x18, 0x19, 0x1B, 0x1C:
		case 0x13:
			p.Hangable = true
		case 0x15, 0x16:
		case 0x1A:
			if _, err := br.u16(); err != nil {
				return p, err
			}
		default:
			return p, unknownFlag(flag, id)
		}
	}
}

func readPropertiesV3(br *binReader, id uint16) (Properties, error) {
	var p Properties
	for {
		flag, err := br.u8()
		if err != nil {
			return p, err
		}
		if flag == 0xFF {
			return p, nil
		}
		switch flag {
		case 0x00:
			if _, err := br.u16(); err != nil {
				return p, err
			}
			p.Ground = true
		case 0x01:
			p.GroundBorder = true
		case 0x02:
			p.OnBottom = true
		case 0x03:
			p.OnTop = true
		case 0x04:
			p.Container = true
		case 0x05:
			p.Stackable = true
		case 0x06:
			p.ForceUse = true
		case 0x07:
			p.MultiUse = true
		case 0x08, 0x09:
			if _, err := br.u16(); err != nil {
				return p, err
			}
		case 0x0A:
			p.FluidContainer = true
		case 0x0B:
			p.Fluid = true
		case 0x0C:
			p.Unpassable = true
		case 0x0D:
			p.Unmoveable = true
		case 0x0E:
			p.BlockMissile = true
		case 0x0F:
			p.BlockPathfind = true
		case 0x10:
			p.Pickupable = true
		case 0x11:
			p.Hangable = true
		case 0x12, 0x13, 0x14, 0x17, 0x1A, 0x1B, 0x1E:
		case 0x15:
			if err := br.skip(4); err != nil {
				return p, err
			}
		case 0x18:
			if err := br.skip(4); err != nil {
				return p, err
			}
		case 0x19:
			if _, err := br.u16(); err != nil {
				return p, err
			}
		case 0x1C:
			if _, err := br.u16(); err != nil {
				return p, err
			}
		case 0x1D:
			if _, err := br.u16(); err != nil {
				return p, err
			}
		default:
			return p, unknownFlag(flag, id)
		}
	}
}

func readPropertiesV4(br *binReader, id uint16) (Properties, error) {
	var p Properties
	for {
		flag, err := br.u8()
		if err != nil {
			return p, err
		}
		if flag == 0xFF {
			return p, nil
		}
		switch flag {
		case 0x00:
			if _, err := br.u16(); err != nil {
				return p, err
			}
			p.Ground = true
		case 0x01:
			p.GroundBorder = true
		case 0x02:
			p.OnBottom = true
		case 0x03:
			p.OnTop = true
		case 0x04:
			p.Container = true
		case 0x05:
			p.Stackable = true
		case 0x06:
			p.ForceUse = true
		case 0x07:
			p.MultiUse = true
		case 0x08:
			// has charges
		case 0x09, 0x0A:
			if _, err := br.u16(); err != nil {
				return p, err
			}
		case 0x0B:
			p.FluidContainer = true
		case 0x0C:
			p.Fluid = true
		case 0x0D:
			p.Unpassable = true
		case 0x0E:
			p.Unmoveable = true
		case 0x0F:
			p.BlockMissile = true
		case 0x10:
			p.BlockPathfind = true
		case 0x11:
			p.Pickupable = true
		case 0x12:
			p.Hangable = true
		case 0x13, 0x14, 0x17, 0x19, 0x1A, 0x1B, 0x1C, 0x1E, 0x1F:
		case 0x15:
			if err := br.skip(4); err != nil {
				return p, err
			}
		case 0x18:
			if err := br.skip(4); err != nil {
				return p, err
			}
		case 0x16:
			if _, err := br.u16(); err != nil {
				return p, err
			}
		case 0x1D:
			if _, err := br.u16(); err != nil {
				return p, err
			}
		default:
			return p, unknownFlag(flag, id)
		}
	}
}

func readPropertiesV5(br *binReader, id uint16) (Properties, error) {
	var p Properties
	for {
		flag, err := br.u8()
		if err != nil {
			return p, err
		}
		if flag == 0xFF {
			return p, nil
		}
		switch flag {
		case 0x00:
			if _, err := br.u16(); err != nil {
				return p, err
			}
			p.Ground = true
		case 0x01:
			p.GroundBorder = true
		case 0x02:
			p.OnBottom = true
		case 0x03:
			p.OnTop = true
		case 0x04:
			p.Container = true
		case 0x05:
			p.Stackable = true
		case 0x06:
			p.ForceUse = true
		case 0x07:
			p.MultiUse = true
		case 0x08, 0x09:
			if _, err := br.u16(); err != nil {
				return p, err
			}
		case 0x0A:
			p.FluidContainer = true
		case 0x0B:
			p.Fluid = true
		case 0x0C:
			p.Unpassable = true
		case 0x0D:
			p.Unmoveable = true
		case 0x0E:
			p.BlockMissile = true
		case 0x0F:
			p.BlockPathfind = true
		case 0x10:
			p.Pickupable = true
		case 0x11:
			p.Hangable = true
		case 0x12, 0x13, 0x14, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F:
		case 0x15:
			if err := br.skip(4); err != nil {
				return p, err
			}
		case 0x16:
			if err := br.skip(4); err != nil {
				return p, err
			}
		case 0x20:
			if _, err := br.u16(); err != nil {
				return p, err
			}
		case 0x21:
			if err := br.skip(12); err != nil {
				return p, err
			}
			nameLen, err := br.u16()
			if err != nil {
				return p, err
			}
			if err := br.skip(int(nameLen)*2 + 4); err != nil {
				return p, err
			}
		default:
			return p, unknownFlag(flag, id)
		}
	}
}

func readPropertiesV6(br *binReader, id uint16) (Properties, error) {
	var p Properties
	for {
		flag, err := br.u8()
		if err != nil {
			return p, err
		}
		if flag == 0xFF {
			return p, nil
		}
		switch flag {
		case 0x00:
			if _, err := br.u16(); err != nil {
				return p, err
			}
			p.Ground = true
		case 0x01:
			p.GroundBorder = true
		case 0x02:
			p.OnBottom = true
		case 0x03:
			p.OnTop = true
		case 0x04:
			p.Container = true
		case 0x05:
			p.Stackable = true
		case 0x06:
			p.ForceUse = true
		case 0x07:
			p.MultiUse = true
		case 0x08, 0x09:
			if _, err := br.u16(); err != nil {
				return p, err
			}
		case 0x0A:
			p.FluidContainer = true
		case 0x0B:
			p.Fluid = true
		case 0x0C:
			p.Unpassable = true
		case 0x0D:
			p.Unmoveable = true
		case 0x0E:
			p.BlockMissile = true
		case 0x0F:
			p.BlockPathfind = true
		case 0x10:
		case 0x11:
			p.Pickupable = true
		case 0x12:
			p.Hangable = true
		case 0x13, 0x14, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F:
		case 0x15:
			if err := br.skip(4); err != nil {
				return p, err
			}
		case 0x16:
			if err := br.skip(4); err != nil {
				return p, err
			}
		case 0x20:
			if _, err := br.u16(); err != nil {
				return p, err
			}
		case 0x21:
			if err := br.skip(12); err != nil {
				return p, err
			}
			nameLen, err := br.u16()
			if err != nil {
				return p, err
			}
			if err := br.skip(int(nameLen)*2 + 4); err != nil {
				return p, err
			}
		case 0x22:
			if _, err := br.u16(); err != nil {
				return p, err
			}
		case 0x23:
			p.Usable = true
		default:
			return p, unknownFlag(flag, id)
		}
	}
}

func unknownFlag(flag byte, id uint16) error {
	return fmt.Errorf("dat: unknown property flag 0x%02X for item %d", flag, id)
}
