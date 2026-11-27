package wopan

type FCloudProductOrdListQryCtxData struct {
	FcloudProductOrds []struct {
		ActiviteCode       string `json:"activiteCode"`
		ApplyTime          string `json:"applyTime"`
		EffectState        string `json:"effectState"`
		ExpireTime         string `json:"expireTime"`
		IsAutoSub          string `json:"isAutoSub"`
		OrderState         string `json:"orderState"`
		PackageProductCode string `json:"packageProductCode"`
		VipDesc            string `json:"vipDesc"`
		VipLevel           string `json:"vipLevel"`
	} `json:"fcloudProductOrds"`
	MaxVipLevel string `json:"maxVipLevel"`
	IsShowInlet string `json:"isShowInlet"`
}

func (w *WoClient) FCloudProductOrdListQry(opts ...RestyOption) (*FCloudProductOrdListQryCtxData, error) {
	var resp FCloudProductOrdListQryCtxData
	_, err := w.RequestWoHome(KeyFCloudProductOrdListQry, nil, Json{
		"qryType":  "1",
		"clientId": w.clientID,
	}, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type QueryCloudUsageInfoData struct {
	Code      string `json:"code"`
	UsageInfo struct {
		TotalSize     string `json:"totalSize"`
		UsedSize      int64  `json:"usedSize"`
		ImageSize     int64  `json:"imageSize"`
		VideoSize     int64  `json:"videoSize"`
		AudioSize     int64  `json:"audioSize"`
		TextSize      int64  `json:"textSize"`
		OtherSize     int64  `json:"otherSize"`
		ByteUsedSize  int64  `json:"byteUsedSize"`
		ByteTotalSize string `json:"byteTotalSize"`
	} `json:"usageInfo"`
	VipLevel   string `json:"vipLevel"`
	ExpireTime string `json:"expireTime"`
	ApplyTime  string `json:"applyTime"`
	PayType    string `json:"payType"`
	Source     string `json:"source"`
	OrderState string `json:"orderState"`
	Status     string `json:"status"`
}

// QueryCloudUsageInfo returns the quota and the membership level. phoneNum is
// what the Android client sends; when Phone is empty it is taken from
// AppQueryUser.
func (w *WoClient) QueryCloudUsageInfo(opts ...RestyOption) (*QueryCloudUsageInfoData, error) {
	err := w.InitPhone()
	if err != nil {
		return nil, err
	}
	var resp QueryCloudUsageInfoData
	_, err = w.RequestWoHome(KeyQueryCloudUsageInfo, Json{
		"phoneNum": w.Phone,
		"clientId": w.clientID,
	}, JsonSecret, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// FCloudProductPackage is not required

type ClassifyRuleData struct {
	FileTypes map[string]struct {
		SubType string `json:"subType"`
		Ability string `json:"ability"`
		Type    string `json:"type"`
	} `json:"fileTypes"`
	FileIcons struct {
		App struct {
			Zip   string `json:"zip"`
			Image string `json:"image"`
			Other string `json:"other"`
			Xmind string `json:"xmind"`
			Gif   string `json:"gif"`
			Csv   string `json:"csv"`
			Video string `json:"video"`
			Excel string `json:"excel"`
			Pdf   string `json:"pdf"`
			Ppt   string `json:"ppt"`
			Doc   string `json:"doc"`
			Audio string `json:"audio"`
			Text  string `json:"text"`
			Word  string `json:"word"`
		} `json:"app"`
		H5 struct {
		} `json:"h5"`
	} `json:"fileIcons"`
}

func (w *WoClient) ClassifyRule(opts ...RestyOption) (*ClassifyRuleData, error) {
	var resp ClassifyRuleData
	_, err := w.RequestWoHome(KeyClassifyRule, Json{}, Json{
		"key": true,
	}, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type GetZoneInfoData struct {
	Url string `json:"url"`
}

func (w *WoClient) GetZoneInfo(opts ...RestyOption) (*GetZoneInfoData, error) {
	var resp GetZoneInfoData
	_, err := w.RequestWoHome(KeyGetZoneInfo, Json{
		"appId": w.appID,
	}, Json{
		"key": true,
	}, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type FamilyUserCurrentEncodeData struct {
	Count           string `json:"count"`
	DefaultHomeId   int    `json:"defaultHomeId"`
	DefaultHomeName string `json:"defaultHomeName"`
	GroupHeadUrl    string `json:"groupHeadUrl"`
	GroupName       string `json:"groupName"`
	Id              int    `json:"id"`
	MemberRole      string `json:"memberRole"`
	OwnerId         string `json:"ownerId"`
	UnreadFlag      string `json:"unreadFlag"`
}

func (w *WoClient) FamilyUserCurrentEncode(opts ...RestyOption) (*FamilyUserCurrentEncodeData, error) {
	var resp FamilyUserCurrentEncodeData
	_, err := w.RequestWoHome(KeyFamilyUserCurrentEncode, Json{
		"clientId": w.clientID,
	}, JsonSecret, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// FamilyGroup is the subset of FamilyListBean's group entries the driver cares
// about: the id is what QueryAllFiles expects as familyId.
type FamilyGroup struct {
	Id           string `json:"id"`
	GroupId      string `json:"groupId"`
	GroupName    string `json:"groupName"`
	GroupHeadUrl string `json:"groupHeadUrl"`
	IsDefault    string `json:"isDefault"`
	OwnerId      string `json:"ownerId"`
	Status       string `json:"status"`
}

type FamilyListData struct {
	FamilyCount        string        `json:"familyCount"`
	MaxableFamilyCount string        `json:"maxableFamilyCount"`
	ICreateGroups      []FamilyGroup `json:"iCreateGroups"`
	IJoinGroups        []FamilyGroup `json:"iJoinGroups"`
	ICreateGroupsCount int           `json:"iCreateGroupsCount"`
	IJoinGroupsCount   int           `json:"iJoinGroupsCount"`
	UnreadFlag         string        `json:"unreadFlag"`
}

// FamilyList lists the family groups the account owns or joined
// (FamilyListRequest.java:19-25). It is the only way to discover a familyId:
// login itself never returns one.
func (w *WoClient) FamilyList(opts ...RestyOption) (*FamilyListData, error) {
	var resp FamilyListData
	_, err := w.RequestWoHome(KeyFamilyList, Json{
		"clientId": w.clientID,
	}, JsonSecret, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type PrivateSpaceLoginData struct {
	PsToken string `json:"psToken"`
	IsPass  string `json:"isPass"`
	Desc    string `json:"desc"`
}

func (w *WoClient) PrivateSpaceLogin(pwd string, opts ...RestyOption) error {
	var resp PrivateSpaceLoginData
	_, err := w.RequestWoHome(KeyPrivateSpaceLogin, Json{
		"pwd":      pwd,
		"clientId": w.clientID,
	}, JsonSecret, &resp, opts...)
	if err != nil {
		return err
	}
	w.SetPsToken(resp.PsToken)
	return nil
}

// DisplayMembershipCardData carries the per-file upload limit of the account
// (uploadFileSize, in MB), which the client uses to refuse an oversized upload
// before sending a single byte.
type DisplayMembershipCardData struct {
	VipLevel           string `json:"vipLevel"`
	EffectiveEndDate   string `json:"effectiveEndDate"`
	IsAutomaticRenewal string `json:"isAutomaticRenewal"`
	IsNeedUpgrade      string `json:"isNeedUpgrade"`
	UploadFileSize     string `json:"uploadFileSize"`
}

// DisplayMembershipCard is not a dispatcher endpoint: it lives at
// /wohome/v1/displayMembershipCard but shares the same signed+encrypted
// envelope.
func (w *WoClient) DisplayMembershipCard(opts ...RestyOption) (*DisplayMembershipCardData, error) {
	var resp DisplayMembershipCardData
	_, err := w.RequestPath("/wohome/v1/displayMembershipCard", ChannelWoHome,
		"displayMembershipCard", Json{"clientId": w.clientID}, JsonSecret, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
