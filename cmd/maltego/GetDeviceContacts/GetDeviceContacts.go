package main

import (
	"flag"
	"fmt"
	"github.com/dreadl0ck/netcap"
	maltego "github.com/dreadl0ck/netcap/cmd/maltego/maltego"
	"github.com/dreadl0ck/netcap/resolvers"
	"github.com/dreadl0ck/netcap/types"
	"github.com/dustin/go-humanize"
	"github.com/gogo/protobuf/proto"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	//"strconv"
	"strings"
)

var (
	flagVersion = flag.Bool("version", false, "print version and exit")
)

func main() {

	lt := maltego.ParseLocalArguments(os.Args)
	profilesFile := lt.Values["path"]
	mac := lt.Values["mac"]

	// print version and exit
	if *flagVersion {
		fmt.Println(netcap.Version)
		os.Exit(0)
	}

	stdout := os.Stdout
	os.Stdout = os.Stderr
	netcap.PrintBuildInfo()

	f, err := os.Open(profilesFile)
	if err != nil {
		log.Fatal(err)
	}

	// check if its an audit record file
	if !strings.HasSuffix(f.Name(), ".ncap.gz") && !strings.HasSuffix(f.Name(), ".ncap") {
		log.Fatal("input file must be an audit record file")
	}

	os.Stdout = stdout

	r, err := netcap.Open(profilesFile, netcap.DefaultBufferSize)
	if err != nil {
		panic(err)
	}

	// read netcap header
	header := r.ReadHeader()
	if header.Type != types.Type_NC_DeviceProfile {
		panic("file does not contain DeviceProfile records: " + header.Type.String())
	}

	var (
		profile = new(types.DeviceProfile)
		pm  proto.Message
		ok  bool
		trx = maltego.MaltegoTransform{}
	)
	pm = profile

	if _, ok = pm.(types.AuditRecord); !ok {
		panic("type does not implement types.AuditRecord interface")
	}

	for {
		err := r.Next(profile)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		} else if err != nil {
			panic(err)
		}

		if profile.MacAddr == mac {

			for _, ip := range profile.Contacts {

				var ent *maltego.MaltegoEntityObj
				if resolvers.IsPrivateIP(net.ParseIP(ip.Addr)) {
					ent = trx.AddEntity("netcap.InternalContact", ip.Addr)
					ent.SetType("netcap.InternalContact")
				} else {
					ent = trx.AddEntity("netcap.ExternalContact", ip.Addr)
					ent.SetType("netcap.ExternalContact")
				}
				ent.SetValue(ip.Addr + "\n" + ip.Geolocation)

				di := "<h3>Heading</h3><p>Timestamp: " + profile.Timestamp + "</p>"
				ent.AddDisplayInformation(di, "Other")

				ent.AddProperty("addr", "Address", "strict", ip.Addr)
				ent.AddProperty("path", "Path", "strict", profilesFile)
				ent.AddProperty("numPackets", "Num Packets", "strict", strconv.FormatInt(profile.NumPackets, 10))

				ent.SetLinkLabel(strconv.FormatInt(ip.NumPackets, 10) + " pkts\n" + humanize.Bytes(ip.Bytes))
				ent.SetLinkColor("#000000")
				ent.SetLinkThickness(getThickness(ip.NumPackets))

				ip.SrcPorts = nil
				ip.DstPorts = nil
				note := strings.ReplaceAll(proto.MarshalTextString(ip), "\"", "'")
				note = strings.ReplaceAll(note, "<", "")
				note = strings.ReplaceAll(note, ">", "")
				ent.SetNote(note)
			}
		}
	}

	trx.AddUIMessage("completed!","Inform")
	fmt.Println(trx.ReturnOutput())
}

func getThickness(val int64) int {
	switch {
	case val < 10:
		return 1
	case val < 100:
		return 2
	case val < 1000:
		return 3
	case val < 10000:
		return 4
	case val < 100000:
		return 5
	default:
		return 1
	}
}