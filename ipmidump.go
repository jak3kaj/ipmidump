package ipmidump

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"

	"github.com/u-root/u-root/pkg/ipmi"
)

var t bool = true

var event = map[int]string{
	0x10: "IPMI command",
	0x08: "power fault",
	0x04: "power interlock",
	0x02: "power overload",
	0x01: "AC failed",
	0x00: "none",
}

var policy = map[int]string{
	0x0: "always-off",
	0x1: "previous",
	0x2: "always-on",
	0x3: "unknown",
}

var statusBool = map[byte]*bool{
	0x80: &t,
	0x00: new(bool),
}

var lanFields = []string {
	"setInProgress",             /* bit  0 */
	"",                          /* bit  1 */
	"",                          /* bit  2 */
	"IPAddress",                 /* bit  3 */
	"IPAddressSrc",              /* bit  4 */
	"MACAddress",                /* bit  5 */
	"Netmask",                   /* bit  6 */
	"",                          /* bit  7 */
	"RMCPPort",                  /* bit  8 */
	"RMCPPortAlt",               /* bit  9 */
	"ARPSettings",               /* bit 10 */
	"ARPInterval",               /* bit 11 */
	"DefaultGWIP",               /* bit 12 */
	"DefaultGWMAC",              /* bit 13 */
	"BackupGWIP",                /* bit 14 */
	"BackupGWMAC",               /* bit 15 */
	"SNMPCommunityString",       /* bit 16 */
	"",                          /* bit 17 */
	"",                          /* bit 18 */
	"",                          /* bit 19 */
	"VLANID",                    /* bit 20 */
	"VLANPriority",              /* bit 21 */
	"",                          /* bit 22 */
	"",                          /* bit 23 */
	"",                          /* bit 24 */
	"",                          /* bit 25 */
	"",                          /* bit 26 */
	"",                          /* bit 27 */
	"",                          /* bit 28 */
	"DHCPServerIP",              /* bit 29 */
	"DHCPServerMAC",             /* bit 30 */
	"DHCPEnable",                /* bit 31 */
}

var lanFieldStrings = map[string][]string {
	"setInProgress": {
		"Set Complete",          /* bit 0 */
		"Set In Progress",       /* bit 1 */
		"Commit Write",          /* bit 2 */
		"Reserved",              /* bit 3 */
	},
	"IPAddressSrc": {
		"Unspecified",           /* bit 0 */
		"Static Address",        /* bit 1 */
		"DHCP Address",          /* bit 2 */
		"BIOS Assigned Address", /* bit 3 */
	},
}

var byteCount = map[string]int{
	"setInProgress":         3,
	"IPAddress":             6,
	"IPAddressSrc":          3,
	"MACAddress":            8,
	"Netmask":               6,
	"RMCPPort":              4,
	"RMCPPortAlt":           4,
	"ARPSettings":           3,
	"ARPInterval":           3,
	"DefaultGWIP":           6,
	"DefaultGWMAC":          8,
	"BackupGWIP":            6,
	"BackupGWMAC":           8,
	"SNMPCommunityString":  20,
	"VLANID":                4,
	"VLANPriority":          3,
	"DHCPServerIP":          6,
	"DHCPServerMAC":         8,
	"DHCPEnable":            3,
}


var adtlDevSupport = []string{
	"Sensor Device",             /* bit 0 */
	"SDR Repository Device",     /* bit 1 */
	"SEL Device",                /* bit 2 */
	"FRU Inventory Device",      /* bit 3 */
	"IPMB Event Receiver",       /* bit 4 */
	"IPMB Event Generator",      /* bit 5 */
	"Bridge",                    /* bit 6 */
	"Chassis Device",            /* bit 7 */
}

func itob(i int) *bool {
	if i != 0 {
		return &t
	} else {
		return new(bool)
	}
}

func Dump() (*IPMIDump, error) {
	// Open the default ipmi device using the u-root ipmi library
	ipmi, err := ipmi.Open(0)
	if err != nil {
		return &IPMIDump{}, fmt.Errorf("Failed to open ipmi device: %v\n", err)
	}
	defer ipmi.Close()

	cps := CurrentPowerState{}
	cs  := ChassisState{}
	dev := DeviceID{}
	fpb := FrontPanelButton{}
	lan := LAN{}

	// Get Chassis Status data from the BMC
	if status, err := ipmi.GetChassisStatus(); err == nil {
		// Current power status
		cpsdata               := int(status.CurrentPowerState)
		cps.PowerRestorePolicy = policy[(cpsdata>>5)&0x03]
		cps.PowerControlFault  = itob(cpsdata&0x10)
		cps.PowerFault         = itob(cpsdata&0x08)
		cps.PowerInterlock     = itob(cpsdata&0x04)
		cps.PowerOverload      = itob(cpsdata&0x02)
		cps.PowerStatus        = itob(cpsdata&0x01)

		// Last power event
		lpedata               := int(status.LastPowerEvent)
		cps.LastPowerEvent     = event[lpedata&0x1F]

		// Misc. chassis state
		mcsdata               := int(status.MiscChassisState)
		cs.FanFault            = itob(mcsdata&0x08)
		cs.DriveFault          = itob(mcsdata&0x04)
		cs.FrontPanelLockout   = itob(mcsdata&0x02)
		cs.ChassisIntrusion    = itob(mcsdata&0x01)

		// Front panel button (optional)
		fpbdata               := int(status.FrontPanelButton)
		// Check if any Front Panel Buttons are being pushed
		if *itob(fpbdata) {
			fpb.PoweroffButton           = itob(fpbdata&0x01)
			fpb.ResetButton              = itob(fpbdata&0x02)
			fpb.DiagnosticButton         = itob(fpbdata&0x04)
			fpb.StandbyButton            = itob(fpbdata&0x08)

			fpb.PoweroffButtonDisable    = itob(fpbdata&0x10)
			fpb.ResetButtonDisable       = itob(fpbdata&0x20)
			fpb.DiagnosticButtonDisable  = itob(fpbdata&0x40)
			fpb.StandbyButtonDisable     = itob(fpbdata&0x80)
		}
	}

	// Loop through each element in the lanFields Slice (order matters)
	for i, name := range lanFields {
		// Skip any blank entries in lanFields, but increment counter
		if name == "" {
			continue
		}
		// Increment lanConfigParameter byte 
		lanConfigParameter := byte(i)
		// Get Lan Config data from the BMC, in order, from all valid fields
		if buf, err := ipmi.GetLanConfig(1, lanConfigParameter); err == nil {
			buflen := len(buf)
			// Validate field has expected byte count
			if buflen == byteCount[name] {
				switch buflen {
				case 6:
					switch name {
					case "IPAddress":
						lan.IPAddress     = net.IP(buf[2:]).String()
					case "Netmask":
						lan.Netmask       = net.IP(buf[2:]).String()
					case "DefaultGWIP":
						lan.DefaultGWIP   = net.IP(buf[2:]).String()
					case "BackupGWIP":
						lan.BackupGWIP    = net.IP(buf[2:]).String()
					case "DHCPServerIP":
						lan.DHCPServerIP  = net.IP(buf[2:]).String()
					}
				case 8:
					switch name {
					case "MACAddress":
						lan.MACAddress    = net.HardwareAddr(buf[2:]).String()
					case "DefaultGWMAC":
						lan.DefaultGWMAC  = net.HardwareAddr(buf[2:]).String()
					case "BackupGWMAC":
						lan.BackupGWMAC   = net.HardwareAddr(buf[2:]).String()
					case "DHCPServerMAC":
						lan.DHCPServerMAC = net.HardwareAddr(buf[2:]).String()
					}
				case 4:
					switch name {
					case "RMCPPort":
						lan.RMCPPort      = int(binary.LittleEndian.Uint16(buf[2:]))
					case "RMCPPortAlt":
						lan.RMCPPortAlt   = int(binary.LittleEndian.Uint16(buf[2:]))
					case "VLANID":
						lan.VLANID        = int(binary.LittleEndian.Uint16(buf[2:]))
					}
				case 3:
					switch name {
				    case "setInProgress":
						if int(buf[2]) < len(lanFieldStrings[name]) {
							lan.setInProgress = lanFieldStrings[name][buf[2]]
						}
				    case "IPAddressSrc":
						if int(buf[2]) < len(lanFieldStrings[name]) {
							lan.IPAddressSrc  = lanFieldStrings[name][buf[2]]
						}
					// Two settings are configured in this field
					case "ARPSettings":
						lan.GratuitousARP = itob(int(buf[2]&1))
						lan.ARPResponses  = itob(int(buf[2]&2))
					case "DHCPEnable":
						lan.DHCPEnable    = itob(int(buf[2]&1))
					case "ARPInterval":
						lan.ARPInterval   = int(buf[2])
					case "VLANPriority":
						lan.VLANPriority  = int(buf[2])
					}
				case 20:
					switch name {
					case "SNMPCommunityString":
						// Deal with Null terminated C Strings properly
						if nullIndex := strings.Index(string(buf[2:]), "\x00"); nullIndex >= 0 {
							lan.SNMPCommunityString = string(buf[2:][:nullIndex])
						}
					}
				}
			}
		}
	}

	// Get Device ID data from BMC
	if info, err := ipmi.GetDeviceID(); err == nil {
		var ver uint8
		ver  = uint8(info.IpmiVersion)

		var mid uint32
		mid  = uint32(info.ManufacturerID[2]) << 16
		mid |= uint32(info.ManufacturerID[1]) << 8
		mid |= uint32(info.ManufacturerID[0])

		var pid uint16
		pid  = uint16(info.ProductID[1]) << 8
		pid |= uint16(info.ProductID[0])

		var supportedDevices []string
		for i := 0; i < 8; i++ {
			if *itob(int(info.AdtlDeviceSupport & (1 << i))) {
				supportedDevices = append(supportedDevices, adtlDevSupport[i])
			}
		}

		dev.ID                      = int(info.DeviceID)
		dev.Rev                     = int(info.DeviceRevision&0x0F)
		dev.FwRev                   = fmt.Sprintf("%d.%02x" , info.FwRev1&0x3F, info.FwRev2)
		dev.IpmiVersion             = fmt.Sprintf("%x.%x", ver&0x0F, (ver&0xF0)>>4)
		dev.ManufacturerID          = int(mid)
		dev.ProductID               = int(pid)
		dev.DeviceAvailable         = statusBool[(^info.FwRev1&0x80)]
		dev.ProvidesSDRs            = statusBool[(info.DeviceRevision&0x80)]
		dev.AdditionalDeviceSupport = supportedDevices
		dev.AuxFwRev                = int(binary.LittleEndian.Uint32(info.AuxFwRev[0:]))
	}

	return &IPMIDump{
		&cps,
		&cs,
		&dev,
		&fpb,
		&lan,
	}, nil

}
