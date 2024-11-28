// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package tbclient

// Errors
type vndError struct {
	Logref  string `json:"logref"`
	Message string `json:"message"`
}

// Collections
type collection[T any] struct {
	Embedded T `json:"_embedded"`
}

type embeddedData[R any] interface {
	getData() []R
}

// Topology
type Topology struct {
	Uid         string `json:"uid,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Datacenter  string `json:"datacenter"`
	Notes       string `json:"notes"`
	Status      string `json:"status,omitempty"`
}

type topologyCollection struct {
	Data []Topology `json:"topologies"`
}

func (t topologyCollection) getData() []Topology {
	return t.Data
}

// InventoryNetwork
type InventoryNetwork struct {
	Id     string `json:"id,omitempty"`
	Type   string `json:"type,omitempty"`
	Subnet string `json:"subnet,omitempty"`
}

type inventoryNetworkCollection struct {
	Data []InventoryNetwork `json:"inventoryNetworks"`
}

func (i inventoryNetworkCollection) getData() []InventoryNetwork {
	return i.Data
}

// Network
type Network struct {
	Uid              string            `json:"uid,omitempty"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	InventoryNetwork *InventoryNetwork `json:"inventoryNetwork"`
	Topology         *Topology         `json:"topology"`
}

type networkCollection struct {
	Data []Network `json:"networks"`
}

func (n networkCollection) getData() []Network {
	return n.Data
}

// Nic Type
type NicType struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type nicTypeCollection struct {
	Data []NicType `json:"vmNetworkInterfaceTypes"`
}

func (n nicTypeCollection) getData() []NicType {
	return n.Data
}

// OS Family
type OsFamily struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type osFamilyCollection struct {
	Data []OsFamily `json:"osFamilies"`
}

func (o osFamilyCollection) getData() []OsFamily {
	return o.Data
}

// Inventory VM

type InventoryVmNic struct {
	InventoryNetworkId string `json:"inventoryNetworkId,omitempty"`
	Name               string `json:"name,omitempty"`
	IpAddress          string `json:"ipAddress,omitempty"`
	MacAddress         string `json:"macAddress,omitempty"`
	Type               string `json:"type,omitempty"`
	RdpEnabled         bool   `json:"rdpEnabled"`
	SshEnabled         bool   `json:"sshEnabled"`
}

type InventoryVmRemoteAccess struct {
	RdpAutoLogin bool `json:"rdpAutoLogin"`
	RdpEnabled   bool `json:"rdpEnabled"`
	SshEnabled   bool `json:"sshEnabled"`
}

type InventoryVm struct {
	Id                  string                   `json:"id,omitempty"`
	Datacenter          string                   `json:"datacenter,omitempty"`
	Name                string                   `json:"name,omitempty"`
	Description         string                   `json:"description,omitempty"`
	OriginalName        string                   `json:"originalName,omitempty"`
	OriginalDescription string                   `json:"originalDescription,omitempty"`
	CpuQty              uint64                   `json:"cpuQty,omitempty"`
	MemoryMb            uint64                   `json:"memoryMb,omitempty"`
	NetworkInterfaces   []InventoryVmNic         `json:"networkInterfaces,omitempty"`
	RemoteAccess        *InventoryVmRemoteAccess `json:"remoteAccess,omitempty"`
}

type inventoryVmCollection struct {
	Data []InventoryVm `json:"inventoryVms"`
}

func (v inventoryVmCollection) getData() []InventoryVm {
	return v.Data
}

// Virtual Machine

type VmNicRdp struct {
	Enabled   bool `json:"enabled"`
	AutoLogin bool `json:"autoLogin"`
}

type VmNicSsh struct {
	Enabled bool `json:"enabled"`
}

type VmNic struct {
	Uid        string    `json:"uid,omitempty"`
	Name       string    `json:"name,omitempty"`
	MacAddress string    `json:"macAddress,omitempty"`
	IpAddress  string    `json:"ipAddress,omitempty"`
	Type       string    `json:"type,omitempty"`
	InUse      bool      `json:"inUse"`
	AssignDhcp bool      `json:"assignDhcp"`
	Rdp        *VmNicRdp `json:"rdp"`
	Ssh        *VmNicSsh `json:"ssh"`
	Network    *Network  `json:"network"`
}

type VmAdvancedSettings struct {
	NameInHypervisor      string `json:"nameInHypervisor,omitempty"`
	BiosUuid              string `json:"biosUuid,omitempty"`
	NotStarted            bool   `json:"notStarted"`
	AllDisksNonPersistent bool   `json:"allDisksNonPersistent"`
}

type VmRemoteAccessDisplayCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type VmRemoteAccessInternalUrl struct {
	Location    string `json:"location,omitempty"`
	Description string `json:"description,omitempty"`
}

type VmRemoteAccess struct {
	Username           string                            `json:"username,omitempty"`
	Password           string                            `json:"password,omitempty"`
	VmConsoleEnabled   bool                              `json:"vmConsoleEnabled"`
	DisplayCredentials *VmRemoteAccessDisplayCredentials `json:"displayCredentials,omitempty"`
	InternalUrls       []VmRemoteAccessInternalUrl       `json:"internalUrls"`
}

type VmGuestAutomation struct {
	Command   string `json:"command,omitempty"`
	DelaySecs uint32 `json:"delaySecs"`
}

type VmDhcpConfig struct {
	DefaultGatewayIp string `json:"defaultGatewayIp"`
}

type Vm struct {
	Uid                  string              `json:"uid,omitempty"`
	Name                 string              `json:"name,omitempty"`
	Description          string              `json:"description,omitempty"`
	MemoryMb             uint64              `json:"memoryMb,omitempty"`
	CpuQty               uint64              `json:"cpuQty,omitempty"`
	NestedHypervisor     *bool               `json:"nestedHypervisor"`
	InventoryVmId        string              `json:"inventoryVmId,omitempty"`
	TopologyInvariantUid string              `json:"topologyInvariantUid,omitempty"`
	OsFamily             string              `json:"osFamily,omitempty"`
	RemoteAccess         *VmRemoteAccess     `json:"remoteAccess"`
	VmNetworkInterfaces  []VmNic             `json:"vmNetworkInterfaces"`
	AdvancedSettings     *VmAdvancedSettings `json:"advancedSettings"`
	GuestAutomation      *VmGuestAutomation  `json:"guestAutomation"`
	ShutdownAutomation   *VmGuestAutomation  `json:"shutdownAutomation"`
	DhcpConfig           *VmDhcpConfig       `json:"dhcpConfig"`
	Topology             *Topology           `json:"topology"`
}

type vmCollection struct {
	Data []Vm `json:"vms"`
}

func (v vmCollection) getData() []Vm {
	return v.Data
}

// Hardware Scripts

type InventoryHwScript struct {
	Uid  string `json:"uid"`
	Name string `json:"name,omitempty"`
}

type inventoryHwScriptCollection struct {
	Data []InventoryHwScript `json:"inventoryHardwareScripts"`
}

func (hws inventoryHwScriptCollection) getData() []InventoryHwScript {
	return hws.Data
}

// Hardware

type InventoryHwNic struct {
	Id string `json:"id"`
}

type InventoryHw struct {
	Id                       string           `json:"id,omitempty"`
	Name                     string           `json:"name,omitempty"`
	Description              string           `json:"description,omitempty"`
	PowerControlAvailable    bool             `json:"powerControlAvailable"`
	HardwareConsoleAvailable bool             `json:"hardwareConsoleAvailable"`
	NetworkInterfaces        []InventoryHwNic `json:"networkInterfaces"`
}

type inventoryHwCollection struct {
	Data []InventoryHw `json:"inventoryHardwareItems"`
}

func (hw inventoryHwCollection) getData() []InventoryHw {
	return hw.Data
}

type HwNic struct {
	Uid              string         `json:"uid,omitempty"`
	Network          Network        `json:"network"`
	NetworkInterface InventoryHwNic `json:"networkInterface"`
}

type Hw struct {
	Uid                    string             `json:"uid,omitempty"`
	Name                   string             `json:"name,omitempty"`
	PowerControlEnabled    *bool              `json:"powerControlEnabled"`
	HardwareConsoleEnabled *bool              `json:"hardwareConsoleEnabled"`
	StartupScript          *InventoryHwScript `json:"inventoryStartupScript"`
	CustomScript           *InventoryHwScript `json:"inventoryCustomScript"`
	ShutdownScript         *InventoryHwScript `json:"inventoryShutdownScript"`
	TemplateConfigScript   *InventoryHwScript `json:"inventoryTemplateConfigScript"`
	NetworkInterfaces      []HwNic            `json:"hardwareNetworkInterfaces"`
	InventoryHardwareItem  *InventoryHw       `json:"inventoryHardwareItem"`
	Topology               *Topology          `json:"topology"`
}

type hwCollection struct {
	Data []Hw `json:"hardwareItems"`
}

func (hw hwCollection) getData() []Hw {
	return hw.Data
}

// License

type InventoryLicense struct {
	Id       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Quantity int    `json:"quantity,omitempty"`
}

type inventoryLicenseCollection struct {
	Data []InventoryLicense `json:"inventoryLicenses"`
}

func (i inventoryLicenseCollection) getData() []InventoryLicense {
	return i.Data
}

type License struct {
	Uid              string            `json:"uid,omitempty"`
	Quantity         int               `json:"quantity,omitempty"`
	InventoryLicense *InventoryLicense `json:"inventoryLicense"`
	Topology         *Topology         `json:"topology"`
}

type licenseCollection struct {
	Data []License `json:"licenses"`
}

func (l licenseCollection) getData() []License {
	return l.Data
}

// Start/Stop Order

type VmStartPosition struct {
	Position     int `json:"position"`
	DelaySeconds int `json:"delaySeconds"`
	Vm           *Vm `json:"vm"`
}

type VmStartOrder struct {
	Uid       string            `json:"uid,omitempty"`
	Ordered   bool              `json:"ordered"`
	Positions []VmStartPosition `json:"positions"`
	Topology  *Topology         `json:"topology"`
}

type VmStopPosition struct {
	Position int `json:"position"`
	Vm       *Vm `json:"vm"`
}

type VmStopOrder struct {
	Uid       string           `json:"uid,omitempty"`
	Ordered   bool             `json:"ordered"`
	Positions []VmStopPosition `json:"positions"`
	Topology  *Topology        `json:"topology"`
}

type HwStartPosition struct {
	Position     int `json:"position"`
	DelaySeconds int `json:"delaySeconds"`
	Hw           *Hw `json:"hardwareItem"`
}

type HwStartOrder struct {
	Uid       string            `json:"uid,omitempty"`
	Ordered   bool              `json:"ordered"`
	Positions []HwStartPosition `json:"positions"`
	Topology  *Topology         `json:"topology"`
}

// Remote Access

type RemoteAccess struct {
	Uid                string    `json:"uid,omitempty"`
	AnyconnectEnabled  bool      `json:"anyconnectEnabled"`
	EndpointKitEnabled bool      `json:"endpointKitEnabled"`
	Topology           *Topology `json:"topology"`
}

// Scenario

type ScenarioOption struct {
	Uid          string `json:"uid,omitempty"`
	InternalName string `json:"internalName,omitempty"`
	DisplayName  string `json:"displayName,omitempty"`
}

type Scenario struct {
	Uid      string           `json:"uid,omitempty"`
	Question string           `json:"question,omitempty"`
	Enabled  bool             `json:"enabled"`
	Options  []ScenarioOption `json:"scenarioOptions"`
	Topology *Topology        `json:"topology"`
}

// IP NAT Rule

type IpNatTarget struct {
	IpAddress string `json:"ipAddress"`
	Name      string `json:"name"`
}

type IpNatRule struct {
	Uid      string      `json:"uid,omitempty"`
	EastWest bool        `json:"eastWest"`
	Scope    *string     `json:"scope"`
	Target   IpNatTarget `json:"target"`
	Topology *Topology   `json:"topology"`
}

// VM NAT Rule

type VmNatTarget struct {
	IpAddress string `json:"ipAddress,omitempty"`
	Name      string `json:"name,omitempty"`
	VmNic     *VmNic `json:"targetItem"`
}

type VmNatRule struct {
	Uid      string      `json:"uid,omitempty"`
	EastWest bool        `json:"eastWest"`
	Scope    *string     `json:"scope"`
	Target   VmNatTarget `json:"target"`
	Topology *Topology   `json:"topology"`
}

// Inbound Proxy Rule

type TrafficVmNicTarget struct {
	Uid       string `json:"uid,omitempty"`
	IpAddress string `json:"ipAddress,omitempty"`
	Vm        *Vm    `json:"vm,omitempty"`
}

type InboundProxyHyperlink struct {
	Show bool   `json:"show"`
	Text string `json:"text"`
}

type InboundProxyRule struct {
	Uid         string                 `json:"uid,omitempty"`
	TcpPort     int                    `json:"tcpPort"`
	Ssl         bool                   `json:"ssl"`
	UrlPath     string                 `json:"urlPath"`
	VmNicTarget *TrafficVmNicTarget    `json:"vmNicTarget"`
	Hyperlink   *InboundProxyHyperlink `json:"hyperlink"`
	Topology    *Topology              `json:"topology"`
}

type inboundProxyRuleCollection struct {
	Data []InboundProxyRule `json:"inboundProxyRules"`
}

func (r inboundProxyRuleCollection) getData() []InboundProxyRule {
	return r.Data
}

// DNS Asset

type InventoryDnsAsset struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type inventoryDnsAssetCollection struct {
	Data []InventoryDnsAsset `json:"inventoryDnsAssets"`
}

func (a inventoryDnsAssetCollection) getData() []InventoryDnsAsset {
	return a.Data
}

type InventorySrvProtocol struct {
	Id       string `json:"id"`
	Protocol string `json:"protocol"`
}

type inventorySrvProtocolCollection struct {
	Data []InventorySrvProtocol `json:"srvProtocols"`
}

func (p inventorySrvProtocolCollection) getData() []InventorySrvProtocol {
	return p.Data
}

type ExternalDnsNatRule struct {
	Uid string `json:"uid"`
}

type ExternalDnsSrvRecord struct {
	Uid      string `json:"uid,omitempty"`
	Service  string `json:"service"`
	Protocol string `json:"protocol"`
	Port     int    `json:"port"`
}

type ExternalDnsRecord struct {
	Uid               string                 `json:"uid,omitempty"`
	Hostname          string                 `json:"hostname,omitempty"`
	ARecord           string                 `json:"aRecord,omitempty"`
	InventoryDnsAsset *InventoryDnsAsset     `json:"inventoryDnsAsset"`
	NatRule           *ExternalDnsNatRule    `json:"natRule"`
	SrvRecords        []ExternalDnsSrvRecord `json:"srvRecords"`
	Topology          *Topology              `json:"topology"`
}

type externalDnsRecordCollection struct {
	Data []ExternalDnsRecord `json:"externalDnsRecords"`
}

func (r externalDnsRecordCollection) getData() []ExternalDnsRecord {
	return r.Data
}

// Mail Server

type MailServer struct {
	Uid               string              `json:"uid,omitempty"`
	InventoryDnsAsset *InventoryDnsAsset  `json:"inventoryDnsAsset"`
	VmNicTarget       *TrafficVmNicTarget `json:"vmNicTarget"`
	Topology          *Topology           `json:"topology"`
}

type mailServerCollection struct {
	Data []MailServer `json:"mailServers"`
}

func (m mailServerCollection) getData() []MailServer {
	return m.Data
}

// Documentation

type Documentation struct {
	Uid              string `json:"uid,omitempty"`
	DocumentationUrl string `json:"documentationUrl"`
}

// Telephony

type InventoryTelephonyItem struct {
	Id          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type inventoryTelephonyItemCollection struct {
	Data []InventoryTelephonyItem `json:"inventoryTelephonyItems"`
}

func (t inventoryTelephonyItemCollection) getData() []InventoryTelephonyItem {
	return t.Data
}

type TelephonyItem struct {
	Uid                    string                  `json:"uid,omitempty"`
	Name                   string                  `json:"name,omitempty"`
	InventoryTelephonyItem *InventoryTelephonyItem `json:"inventoryTelephonyItem"`
	Topology               *Topology               `json:"topology"`
}

type telephonyItemCollection struct {
	Data []TelephonyItem `json:"telephonyItems"`
}

func (t telephonyItemCollection) getData() []TelephonyItem {
	return t.Data
}
