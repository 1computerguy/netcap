package maltego

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"testing"
)

// additional transforms
var transforms = []TransformCoreInfo{
	{"ToApplicationCategories", "netcap.IPAddr", "Retrieve categories of classified applications"},
	{"ToApplications", "netcap.IPAddr", "Show all applications used by the selected host"},
	{"ToApplicationsForCategory", "maltego.Service", "Retrieve applications seen for a given category"},
	{"ToCookiesForHTTPHost", "maltego.Website", "Retrieve cookies seen for the given host"},
	{"ToCookiesValues", "netcap.HTTPCookie", "Retrieve values for a given cookie identifier"},
	{"ToDHCP", "netcap.IPAddr", "Fetch DHCP options for host"},
	{"ToDNSQuestions", "netcap.IPAddr", "Show all captured DNS questions for the selected host"},
	{"ToDeviceContacts", "netcap.Device", "Get contacts for a device"},
	{"ToDeviceIPs", "netcap.Device", "Get all IPs that the device has been using"},
	{"ToDeviceProfiles", "netcap.PCAP", "Get profiles for devices from network packet captures"},
	{"ToDeviceProfilesWithDPI", "netcap.PCAP", "Retrieve device profiles with DPI enabled"},
	{"ToDstPorts", "netcap.IPAddr", "Retrieve all destination ports seen for the selected host"},
	{"ToFileType", "maltego.File", "Retrieve file type via unix file util"},
	{"ToFileTypes", "netcap.IPAddr", "Get all fille types for the selected IPAddr"},
	{"ToFiles", "netcap.IPAddr", "Get all files seen from the selected IP"},
	{"ToFilesForContentType", "netcap.ContentType", "Get all files for a given content type"},
	{"ToGeolocation", "netcap.IPAddr", "Retrieve the geolocation of an IP address"},
	{"ToHTTPContentTypes", "netcap.IPAddr", "Show all HTTP Content Types seen for the selected host"},
	{"ToHTTPCookies", "netcap.IPAddr", "Retrieve HTTP cookies"},
	{"ToHTTPHosts", "netcap.IPAddr", "Retrieve all hostnames seen via HTTP for the selected host"},
	{"ToHTTPHostsFiltered", "netcap.IPAddr", "Get a list of hosts filtered against a DNS whitelist"},
	{"ToHTTPParameters", "netcap.IPAddr", "Retrieve HTTP parameters"},
	{"ToHTTPServerNames", "netcap.IPAddr", "Retrieve the server names that have been contacted by the selected host"},
	{"ToHTTPStatusCodes", "netcap.IPAddr", "Show all HTTP status codes observed for the selected host"},
	{"ToHTTPURLs", "netcap.IPAddr", "Retrieve all URLs seen for the selected host"},
	{"ToHTTPUserAgents", "netcap.IPAddr", "Retrieve all HTTP user agents seen from the selected host"},
	{"ToIncomingFlowsFiltered", "netcap.IPAddr", "Show all incoming flows filtered against the configured whitelist"},
	{"ToMailAuthToken", "netcap.IPAddr", "Retrieve POP3 auth tokens"},
	{"ToMailFrom", "netcap.IPAddr", "Retrieve all email addresses from the 'From' field"},
	{"ToMailTo", "netcap.IPAddr", "Retrieve all email addresses from the 'To' field"},
	{"ToMailUserPassword", "maltego.Person", "Retrieve the password for a mail user"},
	{"ToMailUsers", "netcap.IPAddr", "Retrieve email users"},
	{"ToMails", "netcap.IPAddr", "Show mails fetched over POP3"},
	{"ToOutgoingFlowsFiltered", "netcap.IPAddr", "Show all outgoing flows filtered against the configured whitelist"},
	{"ToParameterValues", "netcap.HTTPParameter", "Retrieve all values seen for an HTTP parameter"},
	{"ToParametersForHTTPHost", "maltego.Website", "Retrieve HTTP params for a host"},
	{"ToSNIs", "netcap.IPAddr", "Retrieve the TLS Server Name Indicators seen for the selected host"},
	{"ToSrcPorts", "netcap.IPAddr", "Retrieve all source ports seen for the selected host"},
	{"ToURLsForHTTPHost", "maltego.Website", "Retrieve all urls for a given host"},
	{"ToAuditRecords", "netcap.PCAP", "Transform PCAP file into audit records"},
	{"ToCaptureProcess", "netcap.Interface", "Start network capture on the given interface"},
	{"ToDHCPClients", "netcap.DHCPv4AuditRecords", "Show all DHCP Clients"},
	{"ToDNSQuestions", "netcap.DNSAuditRecords", "Show all DNS questions"},
	{"ToFetchedMails", "netcap.POP3AuditRecords", "Show emails fetched over POP3"},
	{"ToFileTypes", "netcap.FileAuditRecords", "Show MIME types for extracted files"},
	{"ToHTTPHostnames", "netcap.HTTPAuditRecords", "Show all visited website hostnames"},
	{"ToIANAServices", "netcap.FlowAuditRecords", "Show all IANA services identified by the flows destination port"},
	{"ToLiveAuditRecords", "netcap.CaptureProcess", "Show current state of captured traffic"},
	{"ToLoginInformation", "netcap.CredentialsAuditRecords", "Show captured login credentials"},
	{"ToProducts", "netcap.SoftwareAuditRecords", "Show software products and version information"},
	{"ToSSHClients", "netcap.SSHAuditRecords", "Show detected SSH client"},
	{"ToSSHServers", "netcap.SSHAuditRecords", "Show all SSH server software"},
	{"ToSoftwareExploits", "netcap.ExploitAuditRecords", "Show potential exploits "},
	{"ToSoftwareVulnerabilities", "netcap.VulnerabilityAuditRecords", "Show all discovered vulnerable software"},
	{"ToTCPServices", "netcap.ServiceAuditRecords", "Show detected TCP services"},
	{"ToUDPServices", "netcap.ServiceAuditRecords", "Show detected UDP services"},
	{"ToDevices", "netcap.DeviceProfileAuditRecords", "Show all discovered device audit records"},

	{"OpenExploit", "netcap.Exploit", "Open the exploit source with the default system program for the filetype"},
	{"OpenFile", "maltego.File", "Opens a file with the default application"},
	{"OpenFolder", "maltego.File", "Open the folder in which the file resides"},
	{"OpenNetcapFolder", "netcap.PCAP", "Open the storage folder for the selected PCAP file"},
	{"OpenVulnerability", "netcap.Vulnerability", "Open link to the vulnerability on NVD homepage"},

	{"LookupDHCPFingerprint", "netcap.DHCPClient", "Resolve the clients DHCP fingerprint via the fingerbank API"},
	{"StopCaptureProcess", "netcap.CaptureProcess", "Stop the NETCAP capture process"},
}

// generate all transforms and pack as archive
func TestGenerateAllTransforms(t *testing.T) {

	genTransformArchive()

	// generate additional entities
	for _, t := range transforms {
		genTransform(t.ID, t.Description, t.InputEntity)
	}

	genServerListing()
	packTransformArchive()

	copyFile("transforms.mtz", filepath.Join(os.Getenv("HOME"), "transforms.mtz"))
}

func TestToTransformDisplayName(t *testing.T) {
	res := toTransformDisplayName("ToTCPServices")
	if res != "To TCP Services [NETCAP]" {
		t.Fatal("unexpected result", res)
	}
}

func genServerListing() {

	srv := Server{
		Name:        "Local",
		Enabled:     true,
		Description: "Local transforms hosted on this machine",
		URL:         "http://localhost",
		LastSync:    "2020-06-23 20:47:24.433 CEST", // TODO
		Protocol: struct {
			Text    string `xml:",chardata"`
			Version string `xml:"version,attr"`
		}{
			Version: "0.0",
		},
		Authentication: struct {
			Text string `xml:",chardata"`
			Type string `xml:"type,attr"`
		}{
			Type: "none",
		},
		Seeds: "",
	}

	for _, t := range transforms {
		srv.Transforms.Transform = append(srv.Transforms.Transform, struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
		}{
			Name: netcapPrefix + t.ID,
		})
	}

	data, err := xml.MarshalIndent(srv, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(filepath.Join("transforms", "Servers", "Local.tas"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func TestGenerateTransformServerListing(t *testing.T) {

	// File: Servers/Local.tas
	expected := `<MaltegoServer name="Local" enabled="true" description="Local transforms hosted on this machine" url="http://localhost">
 <LastSync>2020-06-23 20:47:24.433 CEST</LastSync>
 <Protocol version="0.0"/>
 <Authentication type="none"/>
 <Transforms>
  <Transform name="netcap.ToAuditRecords"/>
 </Transforms>
 <Seeds/>
</MaltegoServer>`

	srv := Server{
		Name:        "Local",
		Enabled:     true,
		Description: "Local transforms hosted on this machine",
		URL:         "http://localhost",
		LastSync:    "2020-06-23 20:47:24.433 CEST", // TODO
		Protocol: struct {
			Text    string `xml:",chardata"`
			Version string `xml:"version,attr"`
		}{
			Version: "0.0",
		},
		Authentication: struct {
			Text string `xml:",chardata"`
			Type string `xml:"type,attr"`
		}{
			Type: "none",
		},
		Transforms: struct {
			Text      string `xml:",chardata"`
			Transform []struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
			} `xml:"Transform"`
		}{
			Text:      "",
			Transform: []struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
			}{
				{
					Name: "netcap.ToAuditRecords",
				},
			},
		},
		Seeds: "",
	}

	data, err := xml.MarshalIndent(&srv, "", " ")
	if err != nil {
		t.Fatal(err)
	}

	compareGeneratedXML(data, expected, t)
}

func TestGenerateTransformSettings(t *testing.T) {

	// File: TransformRepositories/Local/netcap.ToAuditRecords.transformsettings
	expected := `<TransformSettings enabled="true" disclaimerAccepted="false" showHelp="true" runWithAll="true" favorite="false">
 <Properties>
  <Property name="transform.local.command" type="string" popup="false">/usr/local/bin/net</Property>
  <Property name="transform.local.parameters" type="string" popup="false">transform ToAuditRecords</Property>
  <Property name="transform.local.working-directory" type="string" popup="false">/usr/local/</Property>
  <Property name="transform.local.debug" type="boolean" popup="false">true</Property>
 </Properties>
</TransformSettings>`

	tr := TransformSettings{
		Enabled:            true,
		DisclaimerAccepted: false,
		ShowHelp:           true,
		RunWithAll:         true,
		Favorite:           false,
		Property: TransformSettingProperties{
			Items: []TransformSettingProperty{
				{
					Name:  "transform.local.command",
					Type:  "string",
					Popup: false,
					Text:  "/usr/local/bin/net",
				},
				{
					Name:  "transform.local.parameters",
					Type:  "string",
					Popup: false,
					Text:  "transform ToAuditRecords",
				},
				{
					Name:  "transform.local.working-directory",
					Type:  "string",
					Popup: false,
					Text:  "/usr/local/",
				},
				{
					Name:  "transform.local.debug",
					Type:  "boolean",
					Popup: false,
					Text:  "true",
				},
			},
		}}

	data, err := xml.MarshalIndent(&tr, "", " ")
	if err != nil {
		t.Fatal(err)
	}

	compareGeneratedXML(data, expected, t)
}
func TestGenerateTransform(t *testing.T) {

	// File: TransformRepositories/Local/netcap.ToAuditRecords.transform
	expected := `<MaltegoTransform name="netcap.ToAuditRecords" displayName="To Audit Records [NETCAP]" abstract="false" template="false" visibility="public" description="Transform PCAP file into audit records" author="Philipp Mieden" requireDisplayInfo="false">
 <TransformAdapter>com.paterva.maltego.transform.protocol.v2api.LocalTransformAdapterV2</TransformAdapter>
 <Properties>
  <Fields>
   <Property name="transform.local.command" type="string" nullable="false" hidden="false" readonly="false" description="The command to execute for this transform" popup="false" abstract="false" visibility="public" auth="false" displayName="Command line">
    <SampleValue></SampleValue>
   </Property>
   <Property name="transform.local.parameters" type="string" nullable="true" hidden="false" readonly="false" description="The parameters to pass to the transform command" popup="false" abstract="false" visibility="public" auth="false" displayName="Command parameters">
    <SampleValue></SampleValue>
   </Property>
   <Property name="transform.local.working-directory" type="string" nullable="true" hidden="false" readonly="false" description="The working directory used when invoking the executable" popup="false" abstract="false" visibility="public" auth="false" displayName="Working directory">
    <DefaultValue>/</DefaultValue>
    <SampleValue></SampleValue>
   </Property>
   <Property name="transform.local.debug" type="boolean" nullable="true" hidden="false" readonly="false" description="When this is set, the transform&apos;s text output will be printed to the output window" popup="false" abstract="false" visibility="public" auth="false" displayName="Show debug info">
    <SampleValue>false</SampleValue>
   </Property>
  </Fields>
 </Properties>
 <InputConstraints>
  <Entity type="netcap.PCAP" min="1" max="1"/>
 </InputConstraints>
 <OutputEntities/>
 <defaultSets/>
 <StealthLevel>0</StealthLevel>
</MaltegoTransform>`

	tr := newTransform("ToAuditRecords", "Transform PCAP file into audit records", "netcap.PCAP")

	data, err := xml.MarshalIndent(&tr, "", " ")
	if err != nil {
		t.Fatal(err)
	}

	compareGeneratedXML(data, expected, t)
}