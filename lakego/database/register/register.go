package register

import(
    "lakego-admin/lakego/database/interfaces"
    driverRegister "lakego-admin/lakego/database/register/driver"
    databaseRegister "lakego-admin/lakego/database/register/Database"
)

/**
 * 注册驱动
 */
func RegisterDriver(name string, f func() interfaces.Driver) {
    driverRegister.New().With(name, f)
}

/**
 * 批量注册驱动
 */
func RegisterDrivers(drivers map[string]func() interfaces.Driver) {
    for name, f := range drivers {
        RegisterDriver(name, f)
    }
}

/**
 * 获取已注册驱动
 */
func GetDriver(name string, once ...bool) interfaces.Driver {
    return driverRegister.New().Get(name, once...)
}

/**
 * 注册数据库
 */
func RegisterDatabase(name string, f func() interfaces.Database) {
    databaseRegister.New().With(name, f)
}

/**
 * 批量注册驱动
 */
func RegisterDatabases(databases map[string]func() interfaces.Database) {
    for name, f := range databases {
        RegisterDatabase(name, f)
    }
}

/**
 * 获取已注册数据库
 */
func GetDatabase(name string, once ...bool) interfaces.Database {
    return databaseRegister.New().Get(name, once...)
}
