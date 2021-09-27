package jwt

import (
    "time"
    "errors"
    "io/ioutil"
    "github.com/dgrijalva/jwt-go"

    "lakego-admin/lakego/support/path"
    "lakego-admin/lakego/support/base64"
)

// JWT
func New() *JWT {
    claims := make(map[string]interface{})

    return &JWT{
        Secret: "123456",
        SigningMethod: "HS256",
        Claims: claims,
    }
}

// 验证方式列表
var signingMethodList = map[string]interface{} {
    "ES256": jwt.SigningMethodES256,
    "ES384": jwt.SigningMethodES384,
    "ES512": jwt.SigningMethodES512,

    "HS256": jwt.SigningMethodHS256,
    "HS384": jwt.SigningMethodHS384,
    "HS512": jwt.SigningMethodHS512,

    "RS256": jwt.SigningMethodRS256,
    "RS384": jwt.SigningMethodRS384,
    "RS512": jwt.SigningMethodRS512,

    "PS256": jwt.SigningMethodPS256,
    "PS384": jwt.SigningMethodPS384,
    "PS512": jwt.SigningMethodPS512,
}

type JWT struct {
    Claims map[string]interface{}

    SigningMethod string
    Secret string
    PrivateKey string
    PublicKey string
    PrivateKeyPassword string // 私钥密码
}

// Audience
func (jwter *JWT) WithAud(aud string) *JWT {
    jwter.Claims["aud"] = aud
    return jwter
}

// ExpiresAt
func (jwter *JWT) WithExp(exp int64) *JWT {
    jwter.Claims["exp"] = time.Now().Add(time.Second * time.Duration(exp)).Unix()
    return jwter
}

// Id
func (jwter *JWT) WithJti(jti string) *JWT {
    jwter.Claims["jti"] = jti
    return jwter
}

// Issuer
func (jwter *JWT) WithIss(iss string) *JWT {
    jwter.Claims["iss"] = iss
    return jwter
}

// NotBefore
func (jwter *JWT) WithNbf(nbf int64) *JWT {
    jwter.Claims["nbf"] = time.Now().Add(time.Second * time.Duration(nbf)).Unix()
    return jwter
}

// Subject
func (jwter *JWT) WithSub(sub string) *JWT {
    jwter.Claims["sub"] = sub
    return jwter
}

// 设置自定义载荷
func (jwter *JWT) WithClaim(key string, value interface{}) *JWT {
    jwter.Claims[key] = value
    return jwter
}

// 设置验证方式
func (jwter *JWT) WithSigningMethod(method string) *JWT {
    jwter.SigningMethod = method
    return jwter
}

// 设置秘钥
func (jwter *JWT) WithSecret(secret string) *JWT {
    jwter.Secret = secret
    return jwter
}

// 设置私钥
func (jwter *JWT) WithPrivateKey(privateKey string) *JWT {
    jwter.PrivateKey = privateKey
    return jwter
}

// 设置公钥
func (jwter *JWT) WithPublicKey(publicKey string) *JWT {
    jwter.PublicKey = publicKey
    return jwter
}

// 设置私钥密码
func (jwter *JWT) WithPrivateKeyPassword(password string) *JWT {
    jwter.PrivateKeyPassword = password
    return jwter
}

// 生成token
func (jwter *JWT) MakeToken() (token string, err error) {
    var methodType jwt.SigningMethod
    if method, ok := signingMethodList[jwter.SigningMethod]; ok {
        methodType = method.(jwt.SigningMethod)
    } else {
        methodType = jwt.SigningMethodHS256
    }

    // 载荷
    claims := make(jwt.MapClaims)
    if len(jwter.Claims) > 0 {
        for k, v := range jwter.Claims {
            claims[k] = v
        }
    }

    jwtToken := jwt.NewWithClaims(methodType, claims)

    var secret interface{}

    if jwter.SigningMethod == "RS256" || jwter.SigningMethod == "RS384" || jwter.SigningMethod == "RS512" {
        // 文件
        keyFile := jwter.FormatPath(jwter.PrivateKey)

        if jwter.PrivateKeyPassword != "" {
            if keyData, e := ioutil.ReadFile(keyFile); e == nil {
                // 密码
                password := base64.Decode(jwter.PrivateKeyPassword)

                secret, err = jwt.ParseRSAPrivateKeyFromPEMWithPassword(keyData, password)

                if err != nil {
                    token = ""
                    return
                }
            } else {
                token = ""
                err = errors.New("PrivateKey not exists")
                return
            }
        } else {
            if keyData, e := ioutil.ReadFile(keyFile); e == nil {
                secret, err = jwt.ParseRSAPrivateKeyFromPEM(keyData)

                if err != nil {
                    token = ""
                    return
                }
            } else {
                token = ""
                err = errors.New("PrivateKey not exists")
                return
            }
        }
    } else if jwter.SigningMethod == "PS256" || jwter.SigningMethod == "PS384" || jwter.SigningMethod == "PS512" {
        // 文件
        keyFile := jwter.FormatPath(jwter.PrivateKey)

        if keyData, e := ioutil.ReadFile(keyFile); e == nil {
            secret, err = jwt.ParseRSAPrivateKeyFromPEM(keyData)

            if err != nil {
                token = ""
                return
            }
        } else {
            token = ""
            err = errors.New("PrivateKey not exists")
            return
        }
    } else if jwter.SigningMethod == "HS256" || jwter.SigningMethod == "HS384" || jwter.SigningMethod == "HS512" {
        secret = jwter.Secret

        // 密码
        secret = base64.Decode(secret.(string))
        secret = []byte(secret.(string))
    } else if jwter.SigningMethod == "ES256" || jwter.SigningMethod == "ES384" || jwter.SigningMethod == "ES512" {
        // 文件
        keyFile := jwter.FormatPath(jwter.PrivateKey)

        if keyData, e := ioutil.ReadFile(keyFile); e == nil {
            secret, err = jwt.ParseECPrivateKeyFromPEM(keyData)

            if err != nil {
                token = ""
                return
            }
        } else {
            token = ""
            err = errors.New("PrivateKey not exists")
            return
        }
    }

    if secret == "" {
        token = ""
        err = errors.New("JWT encode error")
        return
    }

    token, err = jwtToken.SignedString(secret)
    return
}

// 解析 token
func (jwter *JWT) ParseToken(strToken string) (*jwt.Token, error) {
    var err error
    var secret interface{}

    if jwter.SigningMethod == "RS256" || jwter.SigningMethod == "RS384" || jwter.SigningMethod == "RS512" {
        // 文件
        keyFile := jwter.FormatPath(jwter.PublicKey)

        if keyData, e := ioutil.ReadFile(keyFile); e == nil {
            secret, err = jwt.ParseRSAPublicKeyFromPEM(keyData)

            if err != nil {
                return nil, err
            }
        } else {
            err = errors.New("PublicKey not exists")
            return nil, err
        }
    } else if jwter.SigningMethod == "PS256" || jwter.SigningMethod == "PS384" || jwter.SigningMethod == "PS512" {
        // 文件
        keyFile := jwter.FormatPath(jwter.PublicKey)

        if keyData, e := ioutil.ReadFile(keyFile); e == nil {
            secret, err = jwt.ParseRSAPublicKeyFromPEM(keyData)

            if err != nil {
                return nil, err
            }
        } else {
            err = errors.New("PublicKey not exists")
            return nil, err
        }
    } else if jwter.SigningMethod == "HS256" || jwter.SigningMethod == "HS384" || jwter.SigningMethod == "HS512" {
        secret = jwter.Secret

        // 密码
        secret = base64.Decode(secret.(string))
        secret = []byte(secret.(string))
    } else if jwter.SigningMethod == "ES256" || jwter.SigningMethod == "ES384" || jwter.SigningMethod == "ES512" {
        // 文件
        keyFile := jwter.FormatPath(jwter.PublicKey)

        if keyData, e := ioutil.ReadFile(keyFile); e == nil {
            secret, err = jwt.ParseECPublicKeyFromPEM(keyData)

            if err != nil {
                return nil, err
            }
        } else {
            err = errors.New("PublicKey not exists")
            return nil, err
        }
    }

    if secret == "" {
        return nil, errors.New("JWT encode error")
    }

    token, err := jwt.Parse(strToken, func(token *jwt.Token) (interface{}, error) {
        return secret, nil
    })
    if err != nil {
        return nil, err
    }

    return token, nil
}

// 从 token 获取解析后的数据
func (jwter *JWT) GetClaimsFromToken(token *jwt.Token) (jwt.MapClaims, error) {
    var ok bool
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, errors.New("Token claims get error")
    }

    return claims, nil
}

// token 过期检测
func (jwter *JWT) Validate(token *jwt.Token) (bool, error) {
    if err := token.Claims.Valid(); err != nil {
        return false, err
    }

    return true, nil
}

// 验证 token 是否有效
func (jwter *JWT) Verify(token *jwt.Token) (bool, error) {
    if token.Claims.(jwt.MapClaims).VerifyAudience(jwter.Claims["aud"].(string), false) == false {
        return false, errors.New("Audience is error")
    }

    if token.Claims.(jwt.MapClaims).VerifyIssuer(jwter.Claims["iss"].(string), false) == false {
        return false, errors.New("Issuer is error")
    }

    return true, nil
}

// 格式化文件路径
func (jwter *JWT) FormatPath(file string) string {
    filename := path.FormatPath(file)

    return filename
}
