package kernel

import(
    "sync"

    "github.com/deatil/lakego-doak/lakego/provider/interfaces"
)

var instance *Register
var once sync.Once

/**
 * 单例模式
 */
func NewRegister() *Register {
    once.Do(func() {
        instance = &Register{
            providers: make([]Provider, 0),
        }
    })

    return instance
}

// 添加服务提供者
func AddProvider(f func() any) {
    NewRegister().WithProvider(f)
}

// 获取全部服务提供者
func GetAllProvider() []Provider {
    return NewRegister().GetRegisteredProviders()
}

type (
    // 服务提供者接口
    IServiceProvider = interfaces.ServiceProvider

    // 服务提供者函数
    Provider = func() IServiceProvider
)

/**
 * 注册器
 *
 * @create 2021-12-19
 * @author deatil
 */
type Register struct {
    // 锁定
    mu sync.RWMutex

    // 服务提供者
    providers []Provider
}

// 注册
func (this *Register) WithProvider(f func() any) *Register {
    this.mu.Lock()
    defer this.mu.Unlock()

    addProvider := f()

    // 判断是否为服务提供者
    switch addProvider.(type) {
        case IServiceProvider:
            this.providers = append(this.providers, func() IServiceProvider {
                return addProvider.(IServiceProvider)
            })
    }

    return this
}

/**
 * 获取注册的全部服务提供者
 */
func (this *Register) GetRegisteredProviders() []Provider {
    this.mu.RLock()
    defer this.mu.RUnlock()

    return this.providers
}
