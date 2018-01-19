package engine

type ConfigYml struct {
	Debug  int    `json:"debug" comment:"是否显示sql日志"` //是否显示sql日志
	Port   int    `json:"port" comment:"服务端口"`       //服务端口
	DbType string `json:"dbtype" comment:"数据库类型"`    //数据库类型 mssql mysql oci
	DbUser string `json:"dbuser" comment:"数据库用户"`    //数据库用户
	DbPwd  string `json:"dbpwd" comment:"数据库密码"`     //数据库密码
	DbIp   string `json:"dbip" comment:"数据库IP"`      //数据库IP
	DbPort int64  `json:"dbport" comment:"数据库端口"`    //数据库端口
	DbName string `json:"dbname" comment:"数据库名称"`    //数据库名称
}
