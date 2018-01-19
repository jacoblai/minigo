package engine

type ConfigYml struct {
	Debug    int    `json:"debug"`    //是否显示sql日志
	Port     int    `json:"port"`     //服务端口
	DbType   string `json:"dbtype"`   //数据库类型 mssql mysql oci
	DbConStr string `json:"dbconstr"` //数据库用户
}
