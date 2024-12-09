package model

func init() {
	RegisterModel(&SystemInfo{})
}

const (
	// LoginTypeDex is the dex login type
	LoginTypeDex string = "dex"
	// LoginTypeLocal is the local login type
	LoginTypeLocal string = "local"
)

// SystemInfo systemInfo model
type SystemInfo struct {
	BaseModel
	SignedKey        string `json:"signedKey"`
	InstallID        string `json:"installID"`
	EnableCollection bool   `json:"enableCollection"`
	LoginType        string `json:"loginType"`
}

// ProjectRef set the project name and roles
type ProjectRef struct {
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}

// TableName return custom table name
func (u *SystemInfo) TableName() string {
	return tableNamePrefix + "system_info"
}

// ShortTableName is the compressed version of table name for storage
func (u *SystemInfo) ShortTableName() string {
	return "sysi"
}

// PrimaryKey return custom primary key
func (u *SystemInfo) PrimaryKey() string {
	return u.InstallID
}

// Index return custom index
func (u *SystemInfo) Index() map[string]interface{} {
	index := make(map[string]interface{})
	if u.InstallID != "" {
		index["installID"] = u.InstallID
	}
	return index
}
