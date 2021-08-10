package register

import(
    "lakego-admin/lakego/cache/interfaces"
    driverRegister "lakego-admin/lakego/cache/register/driver"
    cacheRegister "lakego-admin/lakego/cache/register/cache"
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
 * 注册缓存
 */
func RegisterCache(name string, f func() interfaces.Cache) {
    cacheRegister.New().With(name, f)
}

/**
 * 批量注册缓存
 */
func RegisterCaches(caches map[string]func() interfaces.Cache) {
    for name, f := range caches {
        RegisterCache(name, f)
    }
}

/**
 * 获取已注册缓存
 */
func GetCache(name string, once ...bool) interfaces.Cache {
    return cacheRegister.New().Get(name, once...)
}
