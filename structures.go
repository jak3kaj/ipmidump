package ipmidump

type IPMIDump struct {
	CurrentPowerState       *CurrentPowerState
	ChassisState            *ChassisState
	DeviceID                *DeviceID
	FrontPanelButton        *FrontPanelButton
	LAN                     *LAN
}

type CurrentPowerState struct {
	PowerRestorePolicy      string   `json:"power_restore_policy,omitempty"`
	PowerStatus             *bool    `json:"power_status,omitempty"`
	PowerControlFault       *bool    `json:"power_control_fault,omitempty"`
	PowerFault              *bool    `json:"power_fault,omitempty"`
	PowerInterlock          *bool    `json:"power_interlock,omitempty"`
	PowerOverload           *bool    `json:"power_overload,omitempty"`
	LastPowerEvent          string   `json:"last_power_event,omitempty"`
}

type ChassisState struct {
	ChassisIntrusion        *bool    `json:"chassis_intrusion,omitempty"`
	FrontPanelLockout       *bool    `json:"front_panel_lockout,omitempty"`
	DriveFault              *bool    `json:"drive_fault,omitempty"`
	FanFault                *bool    `json:"fan_fault,omitempty"`
}

type DeviceID struct {
	ID                      int      `json:"device_id"`
	Rev                     int      `json:"device_revision"`
	FwRev                   string   `json:"firmware_revision,omitempty"`
	IpmiVersion             string   `json:"ipmi_version,omitempty"`
	ManufacturerID          int      `json:"manufacturer_id"`
	ProductID               int      `json:"product_id"`
	DeviceAvailable         *bool    `json:"device_available"`
	ProvidesSDRs            *bool    `json:"provides_device_sdrs"`
	AdditionalDeviceSupport []string `json:"additional_device_support,omitempty"`
	AuxFwRev                int      `json:"aux_firmware_revision_info,omitempty"`
}

type FrontPanelButton struct {
	PoweroffButton          *bool    `json:"poweroff_button,omitempty"`
	ResetButton             *bool    `json:"reset_button,omitempty"`
	DiagnosticButton        *bool    `json:"diagnostic_button,omitempty"`
	StandbyButton           *bool    `json:"standby_button,omitempty"`
	PoweroffButtonDisable   *bool    `json:"poweroff_button_disable,omitempty"`
	ResetButtonDisable      *bool    `json:"reset_button_disable,omitempty"`
	DiagnosticButtonDisable *bool    `json:"diagnostic_button_disable,omitempty"`
	StandbyButtonDisable    *bool    `json:"standby_button_disable,omitempty"`
}

type LAN struct {
	setInProgress           string   `json:"set_in_progress,omitempty"`
	IPAddress               string   `json:"ip"`
	IPAddressSrc            string   `json:"ip_address_source"`
	MACAddress              string   `json:"mac"`
	Netmask                 string   `json:"netmask"`
	RMCPPort                int      `json:"rmcp_port,omitempty"`
	RMCPPortAlt             int      `json:"rmcp_port_alt,omitempty"`
	GratuitousARP           *bool    `json:"gratuitous_arp,omitempty"`
	ARPResponses            *bool    `json:"arp_responses,omitempty"`
	ARPInterval             int      `json:"arp_interval,omitempty"`
	DefaultGWIP             string   `json:"default_gw_ip,omitempty"`
	DefaultGWMAC            string   `json:"default_gw_mac,omitempty"`
	BackupGWIP              string   `json:"backup_default_gw_ip,omitempty"`
	BackupGWMAC             string   `json:"backup_default_gw_mac,omitempty"`
	SNMPCommunityString     string   `json:"snmp_community_string,omitempty"`
	VLANID                  int      `json:"vlan_id,omitempty"`
	VLANPriority            int      `json:"vlan_priority,omitempty"`
	DHCPServerIP            string   `json:"dhcp_server_ip,omitempty"`
	DHCPServerMAC           string   `json:"dhcp_server_mac,omitempty"`
	DHCPEnable              *bool    `json:"dhcp_enable,omitempty"`
}
