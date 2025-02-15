package signer

import (
    "errors"
    "github.com/golang-jwt/jwt/v4"

    "github.com/deatil/lakego-jwt/jwt/sm2"
    "github.com/deatil/lakego-jwt/jwt/config"
    "github.com/deatil/lakego-jwt/jwt/interfaces"
)

// SignerGmSM2
func SignerGmSM2(conf config.SignerConfig) interfaces.Signer {
    return GmSM2{
        Config: conf,
        SigningMethod: sm2.SigningMethodGmSM2,
    }
}

/**
 * GmSM2
 *
 * @create 2023-2-5
 * @author deatil
 */
type GmSM2 struct {
    // 配置
    Config config.SignerConfig

    // 签名
    SigningMethod jwt.SigningMethod
}

// 获取签名
func (this GmSM2) GetSigner() jwt.SigningMethod {
    return this.SigningMethod
}

// 签名密钥
func (this GmSM2) GetSignSecrect() (secret any, err error) {
    // 私钥
    keyByte := this.Config.PrivateKey
    if len(keyByte) == 0 {
        err = errors.New("GmSM2 私钥内容不能为空")
        return
    }

    password := this.Config.PrivateKeyPassword

    if password != "" {
        secret, err = sm2.ParseSM2PrivateKeyFromPEMWithPassword(keyByte, password)
    } else {
        secret, err = sm2.ParseSM2PrivateKeyFromPEM(keyByte)
    }

    return
}

// 验证密钥
func (this GmSM2) GetVerifySecrect() (secret any, err error) {
    // 公钥
    keyByte := this.Config.PublicKey
    if len(keyByte) == 0 {
        err = errors.New("GmSM2 公钥内容不能为空")
        return nil, err
    }

    secret, err = sm2.ParseSM2PublicKeyFromPEM(keyByte)
    return
}
