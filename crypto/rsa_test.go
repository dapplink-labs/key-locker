package crypto

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

func TestNewRsa(t *testing.T) {
	data := strings.Repeat("H", 245) + "q"
	privateKey, publicKey := NewRsa("", "").CreatePkcs8Keys(2048)
	fmt.Printf("public key：%v \n private key: %v \n", publicKey, privateKey)
	rsaObj := NewRsa(publicKey, privateKey)
	sData, err := rsaObj.Encrypt([]byte(data))
	if err != nil {
		fmt.Println("encrypt fail:", err)
	}
	pData, err := rsaObj.Decrypt(sData)
	if err != nil {
		fmt.Println("decrypt fail：", err)
	}
	sign, _ := rsaObj.Sign([]byte(data), crypto.SHA256)
	verify := rsaObj.Verify([]byte(data), sign, crypto.SHA256)
	fmt.Printf(" encrypt：%v\n decrypt：%v\n sign：%v\n verify sign result：%v\n",
		hex.EncodeToString(sData),
		string(pData),
		hex.EncodeToString(sign),
		verify,
	)
}
