package globals

import "fmt"

// Basic ProtocolID
var protocolNames = map[uint16]string{
	0x7B: "AA-User",
	0x19: "Account-Management",
	0x73: "Datastore",
	0x14: "Friends",
	0x70: "Ranking",
	0x0B: "Secure-Connection",
	0x0A: "Ticket-Granting",
	0x03: "Nat-Traversal",
	0x15: "Matchmaking",
	0x32: "Matchmaking-Ext",
	0x6D: "Matchmaking-Extension",
	0x6E: "Store-Manager",
}

var methodNames = map[uint16]map[uint32]string{
	// Ticket-Granting
	0x0A: {
		1: "Login",
		2: "LoginEx",
		3: "RequestTicket",
		4: "GetPID",
		5: "GetName",
		6: "LoginWithParam",
	},

	// Secure-Connection
	0x0B: {
		1: "Register",
		2: "RequestConnectionData",
		3: "RequestURLs",
		4: "RegisterEx",
		5: "TestConnectivity",
		6: "UpdateURLs",
		7: "ReplaceURL",
		8: "SendReport",
	},

	// Nat-Traversal
	0x03: {
		1: "RequestProbeInitiationExt",
		2: "ReportNATTraversalResult",
		3: "ReportNATProperties",
		4: "GetRelaySignatureKey",
		5: "ReportNATTraversalResultDetail",
	},

	// AA-User
	0x7B: {
		1: "RegisterApplication",
		2: "UnregisterApplication",
		3: "SetApplicationInfo",
		4: "GetApplicationInfo",
	},

	// Account-Management
	0x19: {
		1:  "CreateAccount",
		2:  "DeleteAccount",
		3:  "DisableAccount",
		4:  "ChangePassword",
		5:  "TestCapability",
		6:  "GetName",
		7:  "GetAccountData",
		8:  "GetPrivateData",
		9:  "GetPublicData",
		10: "GetMultiplePublicData",
		11: "UpdateAccountName",
		12: "UpdateAccountEmail",
		13: "UpdateCustomData",
	},

	// Datastore
	0x73: {
		1:  "PrepareGetObjectV1",
		2:  "PreparePostObjectV1",
		3:  "CompletePostObjectV1",
		4:  "DeleteObject",
		5:  "DeleteObjects",
		6:  "ChangeMeta",
		7:  "ChangeMetas",
		8:  "GetMeta",
		9:  "GetMetas",
		10: "PrepareUpdateObject",
		11: "CompleteUpdateObject",
		12: "SearchObject",
		13: "GetNotificationUrl",
		14: "GetNewArrivedNotifications",
		15: "RateObject",
		16: "GetRating",
		17: "GetRatings",
		18: "ResetRating",
		19: "ResetRatings",
		20: "GetSpecificMetaV1",
		21: "PostMetaBinary",
		22: "TouchObject",
		23: "GetRatingWithLog",
		24: "PreparePostObject",
		25: "PrepareGetObject",
		26: "CompletePostObject",
		27: "GetNewArrivedNotificationsV1",
		28: "GetNewArrivedNotifications",
		29: "GetPersistenceInfo",
		30: "GetPersistenceInfos",
		31: "PerpetuateObject",
		32: "UnperpetuateObject",
		33: "PrepareGetObjectOrMetaBinary",
		34: "GetPasswordInfo",
		35: "GetPasswordInfos",
		36: "GetMetasMultipleParam",
		37: "CompletePostObjects",
	},

	// Ranking
	0x70: {
		1:  "UploadScore",
		2:  "DeleteScore",
		3:  "DeleteAllScores",
		4:  "UploadCommonData",
		5:  "DeleteCommonData",
		6:  "GetCommonData",
		7:  "ChangeAttributes",
		8:  "ChangeAllAttributes",
		9:  "GetRanking",
		10: "GetApproxOrder",
		11: "GetStats",
		12: "GetRankingByPrincipalID",
		13: "GetCategorySetting",
	},

	// Matchmaking
	0x15: {
		1:  "RegisterGathering",
		2:  "UnregisterGathering",
		3:  "UnregisterGatherings",
		4:  "UpdateGathering",
		5:  "Invite",
		6:  "AcceptInvitation",
		7:  "DeclineInvitation",
		8:  "CancelInvitation",
		9:  "GetInvitationsSent",
		10: "GetInvitationsReceived",
		11: "Participate",
		12: "CancelParticipation",
		13: "GetParticipants",
		14: "AddParticipants",
		15: "GetDetailedParticipants",
		16: "GetParticipantsURLs",
		17: "FindBySingleID",
		18: "FindByOwner",
		19: "FindByParticipants",
		20: "FindInvitationsByGatheringID",
		21: "FindByDescriptionLike",
		22: "FindByDescriptionRegex",
		23: "FindByID",
		24: "FindBySQLQuery",
		25: "LaunchSession",
		26: "UpdateSessionURL",
		27: "GetSessionURL",
		28: "GetState",
		29: "SetState",
		30: "ReportStats",
		31: "GetStats",
		32: "DeleteGathering",
		33: "GetPendingDeletions",
		34: "DeleteFromDeletions",
		35: "MigrateGatheringOwnership",
	},

	// Matchmaking-Ext
	0x32: {
		1: "EndParticipation",
		2: "GetParticipants",
		3: "GetDetailedParticipants",
		4: "GetParticipantsURLs",
		5: "GetGatheringRelations",
		6: "DeleteFromDeletions",
	},

	// Matchmaking-Extension
	0x6D: {
		1:  "CloseParticipation",
		2:  "OpenParticipation",
		3:  "AutoMatchmake_Postpone",
		4:  "BrowseMatchmakeSession",
		5:  "BrowseMatchmakeSessionWithHostUrls",
		6:  "CreateMatchmakeSession",
		7:  "JoinMatchmakeSession",
		8:  "ModifyCurrentGameAttribute",
		9:  "UpdateNotificationData",
		10: "GetFriendNotificationData",
		11: "UpdateApplicationBuffer",
		12: "UpdateMatchmakeSessionAttribute",
		13: "GetlstFriendNotificationData",
		14: "UpdateMatchmakeSession",
		15: "AutoMatchmakeWithSearchCriteria_Postpone",
		16: "GetPlayingSession",
	},

	// Store-Manager / Utility
	0x6E: {
		1: "AcquireNexUniqueID",
		2: "AcquireNexUniqueIDWithPassword",
		3: "AssociateNexUniqueIDWithMyPrincipalID",
		4: "AssociateNexUniqueIDsWithMyPrincipalID",
		5: "GetAssociatedNexUniqueIDWithMyPrincipalID",
		6: "GetAssociatedNexUniqueIDsWithMyPrincipalID",
		7: "GetIntegerSettings",
		8: "GetStringSettings",
	},
}

// GetProtocolNameByID returns the protocol name for a given protocol ID
func GetProtocolNameByID(protocolID uint16) string {
	if name, ok := protocolNames[protocolID]; ok {
		return name
	}
	return "Unknown"
}

// GetMethodNameByID returns the method name for a given protocol ID and method ID
func GetMethodNameByID(protocolID uint16, methodID uint32) string {
	if methods, ok := methodNames[protocolID]; ok {
		if name, ok := methods[methodID]; ok {
			return name
		}
	}
	return "Unknown"
}

// FormatCall formats the protocol and method IDs into a human-readable string
func FormatCall(protocolID uint16, methodID uint32) string {
	return fmt.Sprintf("%s (0x%02X) -> %s (%d)",
		GetProtocolNameByID(protocolID), protocolID,
		GetMethodNameByID(protocolID, methodID), methodID)
}
