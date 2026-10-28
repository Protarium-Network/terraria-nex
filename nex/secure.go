package nex

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	nexgo "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	commonmatchmaking "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making"
	commonmatchmakingext "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making-ext"
	commonmatchmakeextension "github.com/PretendoNetwork/nex-protocols-common-go/v2/matchmake-extension"
	commonnattraversal "github.com/PretendoNetwork/nex-protocols-common-go/v2/nat-traversal"
	commonsecure "github.com/PretendoNetwork/nex-protocols-common-go/v2/secure-connection"
	commonutility "github.com/PretendoNetwork/nex-protocols-common-go/v2/utility"
	matchmaking "github.com/PretendoNetwork/nex-protocols-go/v2/match-making"
	matchmakingext "github.com/PretendoNetwork/nex-protocols-go/v2/match-making-ext"
	matchmakeextension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"
	nattraversal "github.com/PretendoNetwork/nex-protocols-go/v2/nat-traversal"
	secure "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"
	utility "github.com/PretendoNetwork/nex-protocols-go/v2/utility"
	_ "github.com/lib/pq"

	"github.com/Protarium-Network/terraria-nex/globals"
)

var SecureServer *nexgo.PRUDPServer
var SecureEndpoint *nexgo.PRUDPEndPoint

func StartSecureServer() {
	dsn := os.Getenv("PN_TRA_POSTGRES_URI")
	if dsn == "" {
		dsn = "postgres://admin:password123@127.0.0.1/terraria?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	SecureServer = nexgo.NewPRUDPServer()

	// Wrong value = 106-0502 / endless CONNECT retransmit. Most 3.4.13 titles
	// on this stack need false; flip with PN_WO3_LEGACY_SIGNATURE=1.
	// LegacyConnectionSignature=True if the NEX version is lower than 3.2.0
	SecureServer.PRUDPV1Settings.LegacyConnectionSignature = globals.LegacyConnectionSignature

	SecureServer.LibraryVersions.SetDefault(nexgo.NewLibraryVersion(globals.NEXMajor, globals.NEXMinor, globals.NEXPatch))
	SecureServer.AccessKey = globals.AccessKey

	// UseStructureHeader=True if the NEX version is lower than 3.3.0
	SecureServer.ByteStreamSettings.UseStructureHeader = true

	SecureEndpoint = nexgo.NewPRUDPEndPoint(1)
	SecureEndpoint.IsSecureEndPoint = true
	SecureEndpoint.ServerAccount = globals.SecureServerAccount
	SecureEndpoint.AccountDetailsByPID = globals.AccountDetailsByPID
	SecureEndpoint.AccountDetailsByUsername = globals.AccountDetailsByUsername

	SecureServer.BindPRUDPEndPoint(SecureEndpoint)
	globals.SecureEndpoint = SecureEndpoint

	SecureEndpoint.OnData(func(packet nexgo.PacketInterface) {
		request := packet.RMCMessage()
		if request == nil {
			fmt.Printf("[TRA Secure] type=%T\n", packet)
			return
		}
		fmt.Printf("[TRA Secure] PID=%d proto=0x%02X method=0x%02X callID=%d paramsLen=%d\n",
			uint64(packet.Sender().PID()), request.ProtocolID, request.MethodID, request.CallID, len(request.Parameters))
	})

	SecureEndpoint.OnConnectionEnded(func(connection *nexgo.PRUDPConnection) {
		fmt.Printf("[TRA Secure] PID=%d disconnected\n", uint64(connection.PID()))
	})

	SecureEndpoint.OnError(func(err *nexgo.Error) {
		globals.Logger.Errorf("[TRA Secure] code=%d msg=%s", err.ResultCode, err.Message)
	})

	matchmakingManager := common_globals.NewMatchmakingManager(SecureEndpoint, db)
	matchmakingManager.GetUserFriendPIDs = func(pid uint32) []uint32 { return []uint32{} }

	globals.UtilityManager = common_globals.NewUtilityManager(SecureEndpoint, db)

	secureProtocol := secure.NewProtocol()
	SecureEndpoint.RegisterServiceProtocol(secureProtocol)
	commonSecureProtocol := commonsecure.NewCommonProtocol(secureProtocol)
	commonSecureProtocol.EnableInsecureRegister()
	commonSecureProtocol.CreateReportDBRecord = func(_ types.PID, _ types.UInt32, _ types.QBuffer) error {
		return nil
	}

	natTraversalProtocol := nattraversal.NewProtocol()
	SecureEndpoint.RegisterServiceProtocol(natTraversalProtocol)
	commonnattraversal.NewCommonProtocol(natTraversalProtocol)

	matchMakingProtocol := matchmaking.NewProtocol()
	SecureEndpoint.RegisterServiceProtocol(matchMakingProtocol)
	commonmatchmaking.NewCommonProtocol(matchMakingProtocol).SetManager(matchmakingManager)

	matchMakingExtProtocol := matchmakingext.NewProtocol()
	SecureEndpoint.RegisterServiceProtocol(matchMakingExtProtocol)
	commonmatchmakingext.NewCommonProtocol(matchMakingExtProtocol).SetManager(matchmakingManager)

	matchmakeExtensionProtocol := matchmakeextension.NewProtocol()
	SecureEndpoint.RegisterServiceProtocol(matchmakeExtensionProtocol)
	commonmatchmakeextension.NewCommonProtocol(matchmakeExtensionProtocol).SetManager(matchmakingManager)

	utilityProtocol := utility.NewProtocol()
	SecureEndpoint.RegisterServiceProtocol(utilityProtocol)
	commonutility.NewCommonProtocol(utilityProtocol).SetManager(globals.UtilityManager)

	port, err := strconv.Atoi(os.Getenv("PN_TRA_SECURE_PORT"))
	if err != nil || port < 1 || port > 65535 {
		port = 27103
	}

	globals.Logger.Successf("[TRA] Secure server listening on UDP %d", port)
	SecureServer.Listen(port)
}
