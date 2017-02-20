// types.go

// This file contains the various types used by the API

package atlas

// APIError is for errors returned by the RIPE API.
type APIError struct {
    Error struct {
        Status int
        Code   int
        Detail string
        Title  string
    }
}

// Key is holding the API key parameters
type Key struct {
	UUID      string `json:"uuid"`
	ValidFrom string `json:"valid_from"`
	ValidTo   string `json:"valid_to"`
	Enabled   bool
	IsActive  bool    `json:"is_active"`
	CreatedAt string  `json:"created_at"`
	Label     string  `json:"label"`
	Grants    []Grant `json:"grants"`
	Type      string  `json:"type"`
}

// Grant is the permission(s) associated with a key
type Grant struct {
	Permission string `json:"permission"`
	Target     struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"target"`
}

// Probe is holding probe's data
type Probe struct {
	AddressV4      string `json:"address_v4"`
	AddressV6      string `json:"address_v6"`
	AsnV4          int    `json:"asn_v4"`
	AsnV6          int    `json:"asn_v6"`
	CountryCode    string `json:"country_code"`
	Description    string `json:"description"`
	FirstConnected int    `json:"first_connected"`
	Geometry       struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	ID            int    `json:"id"`
	IsAnchor      bool   `json:"is_anchor"`
	IsPublic      bool   `json:"is_public"`
	LastConnected int    `json:"last_connected"`
	PrefixV4      string `json:"prefix_v4"`
	PrefixV6      string `json:"prefix_v6"`
	Status        struct {
		Since string `json:"since"`
		ID    int    `json:"id"`
		Name  string `json:"name"`
	} `json:"status"`
	StatusSince int `json:"status_since"`
	Tags        []struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"tags"`
	Type string `json:"type"`
}

// Measurement is what we are working with
type Measurement struct {
	Af                    int                    `json:"af"`
	CreationTime          int                    `json:"creation_time"`
	Description           string                 `json:"description"`
	DestinationOptionSize interface{}            `json:"destination_option_size"`
	DontFragment          interface{}            `json:"dont_fragment"`
	DuplicateTimeout      interface{}            `json:"duplicate_timeout"`
	FirstHop              int                    `json:"first_hop"`
	Group                 string                 `json:"group"`
	GroupID               int                    `json:"group_id"`
	HopByHopOptionSize    interface{}            `json:"hop_by_hop_option_size"`
	ID                    int                    `json:"id"`
	InWifiGroup           bool                   `json:"in_wifi_group"`
	Interval              int                    `json:"interval"`
	IsAllScheduled        bool                   `json:"is_all_scheduled"`
	IsOneoff              bool                   `json:"is_oneoff"`
	IsPublic              bool                   `json:"is_public"`
	MaxHops               int                    `json:"max_hops"`
	PacketInterval        interface{}            `json:"packet_interval"`
	Packets               int                    `json:"packets"`
	Paris                 int                    `json:"paris"`
	ParticipantCount      int                    `json:"participant_count"`
	ParticipationRequests []ParticipationRequest `json:"participation_requests"`
	Port                  interface{}            `json:"port"`
	ProbesRequested       int                    `json:"probes_requested"`
	ProbesScheduled       int                    `json:"probes_scheduled"`
	Protocol              string                 `json:"protocol"`
	ResolveOnProbe        bool                   `json:"resolve_on_probe"`
	ResolvedIPs           []string               `json:"resolved_ips"`
	ResponseTimeout       int                    `json:"response_timeout"`
	Result                string                 `json:"result"`
	Size                  int                    `json:"size"`
	Spread                interface{}            `json:"spread"`
	StartTime             int                    `json:"start_time"`
	Status                struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"status"`
	StopTime  int    `json:"stop_time"`
	Target    string `json:"target"`
	TargetASN int    `json:"target_asn"`
	TargetIP  string `json:"target_ip"`
	Type      string `json:"type"`
}

// ParticipationRequest allow you to add or remove probes from a measurement that
// was already created
type ParticipationRequest struct {
	Action        string `json:"action"`
	CreatedAt     int    `json:"created_at"`
	ID            int    `json:"id"`
	Self          string `json:"self"`
	Measurement   string `json:"measurement"`
	MeasurementID int    `json:"measurement_id"`
	Requested     int    `json:"requested"`
	Type          string `json:"type"`
	Value         string `json:"value"`
	Logs          string `json:"logs"`
}

var (
	// ProbeTypes should be obvious
	ProbeTypes = []string{"area", "country", "prefix", "asn", "probes", "msm"}
	// AreaTypes should also be obvious
	AreaTypes = []string{"WW", "West", "North-Central", "South-Central", "North-East", "South-East"}
)

// MeasurementRequest contains the different measurement to create/view
type MeasurementRequest struct {
	// see below for definition
	Definitions []Definition `json:"definitions"`

	// requested set of probes
	Probes ProbeSet `json:"probes"`
	//
	BillTo       int  `json:"bill_to,omitempty"`
	IsOneoff     bool `json:"is_oneoff,omitempty"`
	SkipDNSCheck bool `json:"skip_dns_check,omitempty"`
	Times        int  `json:"times,omitempty"`
	StartTime    int  `json:"start_time,omitempty"`
	StopTime     int  `json:"stop_time,omitempty"`
}

// ProbeSet is a set of probes obviously
type ProbeSet []struct {
	Requested int               `json:"requested"` // number of probes
	Type      string            `json:"type"`      // area, country, prefix, asn, probes, msm
	Value     string            `json:"value"`     // can be numeric or string
	Tags      map[string]string `json:"tags,omitempty"`
}

// Definition is used to create measurements
type Definition struct {
	// Required fields
	Description string `json:"description"`
	Type        string `json:"type"`
	AF          int    `json:"af"`

	// Required for all but "dns"
	Target string `json:"target,omitempty"`

	GroupID        int    `json:"group_id,omitempty"`
	Group          string `json:"group,omitempty"`
	InWifiGroup    bool   `json:"in_wifi_group,omitempty"`
	Spread         int    `json:"spread,omitempty"`
	Packets        int    `json:"packets,omitempty"`
	PacketInterval int    `json:"packet_interval,omitempty"`

	// Common parameters
	ExtraWait      int  `json:"extra_wait,omitempty"`
	IsOneoff       bool `json:"is_oneoff,omitempty"`
	IsPublic       bool `json:"is_public,omitempty"`
	ResolveOnProbe bool `json:"resolve_on_probe,omitempty"`

	// Default depends on type
	Interval int `json:"interval,omitempty"`

	// dns parameters
	Protocol         string `json:"protocol"`
	QueryClass       string `json:"query_class,omitempty"`
	QueryType        string `json:"query_type,omitempty"`
	QueryArgument    string `json:"query_argument,omitempty"`
	Retry            int    `json:"retry"`
	SetCDBit         bool   `json:"set_cd_bit"`
	SetDOBit         bool   `json:"set_do_bit"`
	SetNSIDBit       bool   `json:"set_nsid_bit"`
	SetRDBit         bool   `json:"set_rd_bit"`
	UDPPayloadSize   int    `json:"udp_payload_size"`
	UseProbeResolver bool   `json:"use_probe_resolver"`

	// ping parameters
	//   none (see target)

	// traceroute parameters
	//   none (see target)

	// ntp parameters
	//   none (see target)

	// http parameters
	//   none (see target)

	// sslcert parameters
	//   none (see target)

	// wifi parameters
	AnonymousIdentity string `json:"anonymous_identity,omitempty"`
	Cert              string `json:"cert,omitempty"`
	EAP               string `json:"eap,omitempty"`
}
