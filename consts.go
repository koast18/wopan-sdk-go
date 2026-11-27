package wopan

// Client credential pairs.
//
// The dispatcher envelope is bound to a client identity in two independent
// ways, so a token can only be refreshed by the very client that minted it:
//
//   - the clientId travels inside the (encrypted) body of every request, and
//   - the AES key of the api-user channel is the client's secret.
//
// Two pairs are in use in the wild:
//
//	1001000035 / iELf0UL07o6I8eRK   Android app (com.chinaunicom.bol.cloudapp)
//	1001000021 / XFmi9GS2hzk98jGX   PC/Web surface, this SDK's historical default
//
// The previous revision hardcoded the PC pair, which made every token harvested
// from the Android app refresh with the wrong secret. Both are configurable now
// (see WoClient.SetClientCredentials); DefaultClientID/DefaultClientSecret keep
// the old values so existing callers behave exactly as before.
const (
	ClientIDAndroid     = "1001000035"
	ClientSecretAndroid = "iELf0UL07o6I8eRK"

	ClientIDPCWeb     = "1001000021"
	ClientSecretPCWeb = "XFmi9GS2hzk98jGX"

	DefaultClientID     = ClientIDPCWeb
	DefaultClientSecret = ClientSecretPCWeb
	DefaultAppID        = "10000001"
	DefaultBaseURL      = "https://panservice.mail.wo.cn"
	// DefaultZoneURL old https://gxupload.pan.wo.cn:8443
	DefaultZoneURL  = "https://tjupload.pan.wo.cn"
	DefaultUA       = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.37"
	DefaultPartSize = int64(8 * 1024 * 1024)

	// DefaultIV is the fixed AES-CBC initialisation vector of every dispatcher
	// envelope (base/util/a.java:30).
	DefaultIV = "wNSOYIB1k1DjY5lA"
)

const (
	ChannelAPIUser = "api-user"
	ChannelWoHome  = "wohome"
	ChannelWoCloud = "wocloud"
)

// spaceType values understood by QueryAllFiles (QueryAllFilesRequest.java:20-32).
const (
	SpaceTypePersonal    = "0"
	SpaceTypeFamily      = "1"
	SpaceTypeWechatBack  = "3"  // 微信备份
	SpaceTypePrivate     = "4"  // 私密空间, needs psToken
	SpaceTypeAlbumBackup = "5"  // 相册备份
	SpaceTypeWoCamera    = "6"  // 沃相机
	SpaceTypeWoRing      = "7"  // 沃铃
	SpaceTypeWoVideo     = "8"  // 沃视频
	SpaceTypeWoMail      = "9"  // 沃邮箱
	SpaceTypeWoReport    = "10" // 沃报告
	SpaceTypeBill        = "12" // 账单
	SpaceTypeOCR         = "13" // OCR
	SpaceTypeMediaSquare = "15" // 媒体广场
)

// api-user methods
const (
	KeyPcWebLogin        = "PcWebLogin"
	KeyPcLoginVerifyCode = "PcLoginVerifyCode"
	KeyAppQueryUser      = "AppQueryUser"
	KeyAppRefreshToken   = "AppRefreshToken"
	KeyAppLogout         = "AppLogout"

	// Android-app login keys (see analysis/02_auth-account.md §3.1/§3.2).
	KeyAppLoginAccount    = "AppLoginAccount"
	KeyAppLoginByMobile   = "AppLoginByMobile"
	KeyAppLoginDoubleAuth = "AppLoginDoubleAuth"
	// KeyAppSendMessageCode is not a dispatcher key: the SMS endpoint is a
	// plain POST to /api-user/sendMessageCodeBase with {func, clientId, param}.
	KeyAppSendMessageCode = "sendMessageCodeBase"
)

// wohome methods
const (
	KeyFCloudProductOrdListQry = "FCloudProductOrdListQry"
	KeyQueryCloudUsageInfo     = "QueryCloudUsageInfo"
	KeyFCloudProductPackage    = "FCloudProductPackage"
	KeyClassifyRule            = "ClassifyRule"
	KeyGetZoneInfo             = "GetZoneInfo"
	KeyQuerySysConfig          = "QuerySysConfig"
	KeyFamilyUserCurrentEncode = "FamilyUserCurrentEncode"
	KeyFamilyList              = "FamilyList"
	KeyQueryAllFiles           = "QueryAllFiles"
	KeyGetSearchDirectory      = "GetSearchDirectory"
	KeyGetDownloadUrlV2        = "GetDownloadUrlV2"
	KeyGetDownloadUrl          = "GetDownloadUrl"
	KeyCreateDirectory         = "CreateDirectory"
	KeyRenameFileOrDirectory   = "RenameFileOrDirectory"
	KeyMoveFile                = "MoveFile"
	KeyCopyFile                = "CopyFile"
	KeyDeleteFile              = "DeleteFile"
	KeyEmptyRecycleData        = "EmptyRecycleData"
	KeyUpload2C                = "upload2C"
	KeyPrivateSpaceLogin       = "PrivateSpaceLogin"
)

// Sort rules accepted by QueryAllFiles.
//
// These mirror QueryAllFilesRequest.SORTRULE_* of the Android client:
//
//	"0" time desc, the client default ("recent files")
//	"1" name ascending
//	"2" name descending
//	"3" time ascending
//
// The server has no size ordering at all, so the former SortSize* values were
// never understood by it (3 meant "time ascending", 4/5/6 did not exist).
const (
	SortTimeDesc = iota // 0
	SortNameAsc         // 1
	SortNameDesc        // 2
	SortTimeAsc         // 3
)
