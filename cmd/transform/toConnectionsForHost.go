package transform

import (
	"github.com/dreadl0ck/netcap/maltego"
	"github.com/dreadl0ck/netcap/resolvers"
	"github.com/dreadl0ck/netcap/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"strconv"
	"strings"
)

func toConnectionsForHost() {

	var (
		resolverLog = zap.New(zapcore.NewNopCore())
	)

	defer func() {
		err := resolverLog.Sync()
		if err != nil {
			log.Println(err)
		}
	}()

	resolvers.SetLogger(resolverLog)

	stdOut := os.Stdout
	os.Stdout = os.Stderr
	resolvers.InitServiceDB()
	os.Stdout = stdOut

	maltego.ConnectionTransform(nil, func(lt maltego.LocalTransform, trx *maltego.Transform, conn *types.Connection, min, max uint64, path string, mac string, ip string, sizes *[]int) {

		// TODO: make showing empty stream configurable, or add a dedicated transform?
		if conn.AppPayloadSize == 0 {
			return
		}

		if conn.SrcIP == ip || conn.DstIP == ip {
			i, err := strconv.Atoi(conn.DstPort)
			if err != nil {
				return
			}
			service := resolvers.LookupServiceByPort(i, strings.ToLower(conn.TransportProto))
			addConn(trx, conn, path, min, max, service)
		}
	})
}
