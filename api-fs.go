package wopan

type File struct {
	FamilyId     int    `json:"familyId"`
	Fid          string `json:"fid"`
	Creator      string `json:"creator"`
	Size         int64  `json:"size"`
	CreateTime   string `json:"createTime"`
	Name         string `json:"name"`
	ShootingTime string `json:"shootingTime"`
	Id           string `json:"id"`
	Type         int    `json:"type"`
	ThumbUrl     string `json:"thumbUrl"`
	FileType     string `json:"fileType"`
}

type QueryAllFilesData struct {
	Files []*File `json:"files"`
}

func (w *WoClient) QueryAllFiles(spaceType, parentDirectoryId string, pageNum, pageSize int, sortRule int, familyId string, opts ...RestyOption) (*QueryAllFilesData, error) {
	var resp QueryAllFilesData
	param := Json{
		"spaceType":         spaceType,
		"parentDirectoryId": parentDirectoryId,
		"pageNum":           pageNum,
		"pageSize":          pageSize,
		"sortRule":          sortRule,
		"clientId":          w.clientID,
	}
	if spaceType == SpaceTypeFamily && familyId != "" {
		param["familyId"] = familyId
	}
	if spaceType == SpaceTypePrivate {
		if w.psToken == "" {
			return nil, ErrInvalidPsToken
		}
		param["psToken"] = w.psToken
	}
	_, err := w.RequestWoHome(KeyQueryAllFiles, param, JsonSecret, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (w *WoClient) QueryAllFilesPersonal(parentDirectoryId string, pageNum, pageSize int, sortRule int, opts ...RestyOption) (*QueryAllFilesData, error) {
	return w.QueryAllFiles(SpaceTypePersonal, parentDirectoryId, pageNum, pageSize, sortRule, "", opts...)
}

func (w *WoClient) QueryAllFilesFamily(parentDirectoryId string, pageNum, pageSize int, sortRule int, familyId string, opts ...RestyOption) (*QueryAllFilesData, error) {
	return w.QueryAllFiles(SpaceTypeFamily, parentDirectoryId, pageNum, pageSize, sortRule, familyId, opts...)
}

// GetDownloadUrlEntry is one element of GetDownloadUrl/GetDownloadUrlV2's list.
// m3u8 and videoDuration are returned for video files and are the ONLY usable
// streaming form: the app prefers m3u8 over downloadUrl for every video.
type GetDownloadUrlEntry struct {
	Fid           string `json:"fid"`
	DownloadUrl   string `json:"downloadUrl"`
	M3u8          string `json:"m3u8"`
	VideoDuration string `json:"videoDuration"`
}

type GetDownloadUrlV2Data struct {
	Type int `json:"type"`
	// DownloadUrl is only populated when the request carried dirId, i.e. when
	// asking for the download URL of a whole directory.
	DownloadUrl string                `json:"downloadUrl"`
	List        []GetDownloadUrlEntry `json:"list"`
}

func (w *WoClient) downloadUrlV2Param() Json {
	param := Json{
		"type":     "1",
		"clientId": w.clientID,
	}
	if w.psToken != "" {
		param["psToken"] = w.psToken
	}
	return param
}

func (w *WoClient) GetDownloadUrlV2(fidList []string, opts ...RestyOption) (*GetDownloadUrlV2Data, error) {
	var resp GetDownloadUrlV2Data
	param := w.downloadUrlV2Param()
	param["fidList"] = fidList
	_, err := w.RequestWoHome(KeyGetDownloadUrlV2, param, JsonSecret, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDownloadUrlV2Dir asks for the download URL of a whole directory. The
// server packs it and answers with the top level downloadUrl field.
func (w *WoClient) GetDownloadUrlV2Dir(dirId string, opts ...RestyOption) (*GetDownloadUrlV2Data, error) {
	var resp GetDownloadUrlV2Data
	param := w.downloadUrlV2Param()
	param["dirId"] = dirId
	_, err := w.RequestWoHome(KeyGetDownloadUrlV2, param, JsonSecret, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type GetDownloadUrlData struct {
	Fid           string `json:"fid"`
	DownloadUrl   string `json:"downloadUrl"`
	M3u8          string `json:"m3u8"`
	VideoDuration string `json:"videoDuration"`
}

func (w *WoClient) GetDownloadUrl(spaceType string, fidList []string, opts ...RestyOption) ([]GetDownloadUrlData, error) {
	var resp []GetDownloadUrlData
	param := Json{
		"fidList":   fidList,
		"clientId":  w.clientID,
		"spaceType": spaceType,
	}
	_, err := w.RequestWoHome(KeyGetDownloadUrl, param, JsonSecret, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type CreateDirectoryData struct {
	Id string `json:"id"`
}

func (w *WoClient) CreateDirectory(spaceType, parentDirectoryId string, directoryName, familyId string, opts ...RestyOption) (*CreateDirectoryData, error) {
	var resp CreateDirectoryData
	param := Json{
		"spaceType":         spaceType,
		"parentDirectoryId": parentDirectoryId,
		"directoryName":     directoryName,
		"clientId":          w.clientID,
	}
	// The app only sends familyId when it is actually set
	// (CreateDirectoryRequest.bodyObject checks TextUtils.isEmpty first).
	if familyId != "" {
		param["familyId"] = familyId
	}
	if spaceType == SpaceTypePrivate {
		if w.psToken == "" {
			return nil, ErrInvalidPsToken
		}
		param["psToken"] = w.psToken
	}
	_, err := w.RequestWoHome(KeyCreateDirectory, param, JsonSecret, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// RenameFileOrDirectory
// _type: 1: file, 0: directory
func (w *WoClient) RenameFileOrDirectory(spaceType string, _type int, id string, name string, familyId string, opts ...RestyOption) error {
	fileType := "0"
	if _type != 0 {
		fileType = w.GetFileType(name)
	}
	param := Json{
		"spaceType": spaceType,
		"type":      _type,
		"fileType":  fileType,
		"id":        id,
		"name":      name,
		"clientId":  w.clientID,
	}
	if spaceType == SpaceTypeFamily && familyId != "" {
		param["familyId"] = familyId
	}
	if spaceType == SpaceTypePrivate {
		if w.psToken == "" {
			return ErrInvalidPsToken
		}
		param["psToken"] = w.psToken
	}
	_, err := w.RequestWoHome(KeyRenameFileOrDirectory, param, JsonSecret, nil, opts...)
	return err
}

func (w *WoClient) RenameFileOrDirectoryPersonal(_type int, id string, name string, opts ...RestyOption) error {
	return w.RenameFileOrDirectory(SpaceTypePersonal, _type, id, name, "", opts...)
}

func (w *WoClient) RenameFileOrDirectoryFamily(_type int, id string, name string, familyId string, opts ...RestyOption) error {
	return w.RenameFileOrDirectory(SpaceTypeFamily, _type, id, name, familyId, opts...)
}

func (w *WoClient) MoveFile(dirList, fileList []string, targetDirId string, sourceType, targetType string, fromFamilyId, targetFamilyId string, opts ...RestyOption) error {
	param := Json{
		"targetDirId": targetDirId,
		"sourceType":  sourceType,
		"targetType":  targetType,
		"dirList":     dirList,
		"fileList":    fileList,
		"secret":      false,
		"clientId":    w.clientID,
	}
	if sourceType == SpaceTypeFamily && fromFamilyId != "" {
		param["fromFamilyId"] = fromFamilyId
	}
	if targetType == SpaceTypeFamily && targetFamilyId != "" {
		param["familyId"] = targetFamilyId
	}
	if (sourceType == SpaceTypePrivate || targetType == SpaceTypePrivate) && w.psToken != "" {
		param["psToken"] = w.psToken
	}
	_, err := w.RequestWoHome(KeyMoveFile, param, JsonSecret, nil, opts...)
	return err
}

func (w *WoClient) CopyFile(dirList, fileList []string, targetDirId string, sourceType, targetType string, fromFamilyId, targetFamilyId string, opts ...RestyOption) error {
	param := Json{
		"targetDirId": targetDirId,
		"sourceType":  sourceType,
		"targetType":  targetType,
		"dirList":     dirList,
		"fileList":    fileList,
		"secret":      false,
		"clientId":    w.clientID,
	}
	if sourceType == SpaceTypeFamily && fromFamilyId != "" {
		param["fromFamilyId"] = fromFamilyId
	}
	if targetType == SpaceTypeFamily && targetFamilyId != "" {
		param["familyId"] = targetFamilyId
	}
	if (sourceType == SpaceTypePrivate || targetType == SpaceTypePrivate) && w.psToken != "" {
		param["psToken"] = w.psToken
	}
	_, err := w.RequestWoHome(KeyCopyFile, param, JsonSecret, nil, opts...)
	return err
}

// DeleteFile removes files/dirs.
//
// vipLevel is not decorative: DeleteFileRequest sends it on every delete, and
// the official client feeds it from the membership state it fetched at startup
// (CloudManager's vipLevel, populated from QueryCloudUsageInfo). Deleting with a
// hardcoded "0" makes the server treat a paying account as a free one.
func (w *WoClient) DeleteFile(spaceType, vipLevel string, dirList, fileList []string, opts ...RestyOption) error {
	if vipLevel == "" {
		vipLevel = "0"
	}
	param := Json{
		"spaceType": spaceType,
		"vipLevel":  vipLevel,
		"dirList":   dirList,
		"fileList":  fileList,
		"clientId":  w.clientID,
	}
	if spaceType == SpaceTypePrivate && w.psToken != "" {
		param["psToken"] = w.psToken
	}
	_, err := w.RequestWoHome(KeyDeleteFile, param, JsonSecret, nil, opts...)
	return err
}

func (w *WoClient) EmptyRecycleData(opts ...RestyOption) error {
	param := Json{
		"clientId": w.clientID,
	}
	_, err := w.RequestWoHome(KeyEmptyRecycleData, param, JsonSecret, nil, opts...)
	return err
}
