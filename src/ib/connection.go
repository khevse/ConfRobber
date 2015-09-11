package ib

/**
* Информация о подключении к информационной базе 1С
 */
type IBConnectionSettings struct {
	ibPath          string // Путь к информационной базе 1С
	ibUserName      string // Имя пользователя информационной базы 1С
	ibUserPwd       string // Пароль пользователя информационной базы 1С
	storagePath     string // Путь к каталогу хранилища
	storageUserName string // Имя пользователя к хранилищу конфигурации 1С
	storageUserPwd  string // Пароль пользователя к хранилищу конфигурации 1С
}

func CreateIBConnectionSettings(ibPath string,
	ibUserName string,
	ibUserPwd string,
	storagePath string,
	storageUserName string,
	storageUserPwd string) (cn *IBConnectionSettings) {

	return &IBConnectionSettings{ibPath: ibPath,
		ibUserName:      ibUserName,
		ibUserPwd:       ibUserPwd,
		storagePath:     storagePath,
		storageUserName: storageUserName,
		storageUserPwd:  storageUserPwd}
}

func (cn IBConnectionSettings) Empty() bool {
	return len(cn.ibPath) > 0
}

func (cn IBConnectionSettings) GetIbPath() string {
	return cn.ibPath
}

func (cn IBConnectionSettings) GetIbUserName() string {
	return cn.ibUserName
}

func (cn IBConnectionSettings) GetIbUserPwd() string {
	return cn.ibUserPwd
}

func (cn IBConnectionSettings) GetStoragePath() string {
	return cn.storagePath
}

func (cn IBConnectionSettings) GetStorageUserName() string {
	return cn.storageUserName
}

func (cn IBConnectionSettings) GetStorageUserPwd() string {
	return cn.storageUserPwd
}
