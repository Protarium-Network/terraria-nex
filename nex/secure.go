package nex

import (
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

	"github.com/Protarium-Network/terraria-nex/database"
	"github.com/Protarium-Network/terraria-nex/globals"
)

func StartSecureServer() {
	globals.SecureServer = nexgo.NewPRUDPServer()

	// Wrong value = 106-0502 / endless CONNECT retransmit.
	// LegacyConnectionSignature are only used if the NEX version is lower than 3.2.0.
	// Terraria is NEX 3.8.3, so it must be false.
	globals.SecureServer.PRUDPV1Settings.LegacyConnectionSignature = globals.LegacyConnectionSignature

	globals.SecureServer.LibraryVersions.SetDefault(nexgo.NewLibraryVersion(globals.NEXMajor, globals.NEXMinor, globals.NEXPatch))
	globals.SecureServer.AccessKey = globals.AccessKey

	// StructureHeaders are used by NEX 3.5.0 and newer (Terraria is 3.8.3).
	globals.SecureServer.ByteStreamSettings.UseStructureHeader = true

	globals.SecureEndpoint = nexgo.NewPRUDPEndPoint(1)
	globals.SecureEndpoint.IsSecureEndPoint = true
	globals.SecureEndpoint.ServerAccount = globals.SecureServerAccount
	globals.SecureEndpoint.AccountDetailsByPID = globals.AccountDetailsByPID
	globals.SecureEndpoint.AccountDetailsByUsername = globals.AccountDetailsByUsername

	globals.SecureServer.BindPRUDPEndPoint(globals.SecureEndpoint)

	globals.SecureEndpoint.OnData(func(packet nexgo.PacketInterface) {
		request := packet.RMCMessage()
		if request == nil {
			fmt.Printf("[TRA Secure] Type=%T\n", packet)
			return
		}
		fmt.Printf("[TRA Secure] PID=%d ProtocolID=%s (0x%02X) MethodID=0x%02X CallID=%d ParamsLen=%d\n",
			uint64(packet.Sender().PID()), globals.GetProtocolNameByID(request.ProtocolID), request.ProtocolID, request.MethodID, request.CallID, len(request.Parameters))
	})

	globals.SecureEndpoint.OnConnectionEnded(func(connection *nexgo.PRUDPConnection) {
		fmt.Printf("[TRA Secure] PID=%d disconnected\n", uint64(connection.PID()))
	})

	globals.SecureEndpoint.OnError(func(err *nexgo.Error) {
		globals.Logger.Errorf("[TRA Secure] ResultCode=%d Message=%s", err.ResultCode, err.Message)
	})

	globals.MatchmakingManager = common_globals.NewMatchmakingManager(globals.SecureEndpoint, database.Postgres)
	globals.MatchmakingManager.GetUserFriendPIDs = globals.GetUserFriendPIDs

	globals.UtilityManager = common_globals.NewUtilityManager(globals.SecureEndpoint, database.Postgres)

	secureProtocol := secure.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(secureProtocol)
	commonSecureProtocol := commonsecure.NewCommonProtocol(secureProtocol)
	commonSecureProtocol.EnableInsecureRegister()
	commonSecureProtocol.CreateReportDBRecord = func(PID types.PID, ReportID types.UInt32, ReportData types.QBuffer) error {
		globals.Logger.Warningf("[TRA Secure] Player report from PID %d (report ID %d, %d bytes) discarded: this server does not store reports",
			PID, ReportID, len(ReportData))
		return nil
	}

	natTraversalProtocol := nattraversal.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(natTraversalProtocol)
	commonnattraversal.NewCommonProtocol(natTraversalProtocol)

	matchMakingProtocol := matchmaking.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(matchMakingProtocol)
	commonmatchmaking.NewCommonProtocol(matchMakingProtocol).SetManager(globals.MatchmakingManager)

	matchMakingExtProtocol := matchmakingext.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(matchMakingExtProtocol)
	commonmatchmakingext.NewCommonProtocol(matchMakingExtProtocol).SetManager(globals.MatchmakingManager)

	matchmakeExtensionProtocol := matchmakeextension.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(matchmakeExtensionProtocol)
	commonmatchmakeextension.NewCommonProtocol(matchmakeExtensionProtocol).SetManager(globals.MatchmakingManager)

	utilityProtocol := utility.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(utilityProtocol)
	commonutility.NewCommonProtocol(utilityProtocol).SetManager(globals.UtilityManager)

	port, err := strconv.Atoi(os.Getenv("PN_TRA_SECURE_PORT"))
	if err != nil || port < 1 || port > 65535 {
		port = 27103
	}

	globals.Logger.Successf("[TRA] Secure server listening on UDP %d", port)
	globals.SecureServer.Listen(port)
}
