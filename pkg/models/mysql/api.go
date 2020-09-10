package mysql

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/ssrdive/mysequel"
	"github.com/ssrdive/straddle/pkg/sql/queries"
)

const otpChars = "1234567890"

type ApiModel struct {
	DB *sql.DB
}

func (m *ApiModel) VerifyHash(countryCode, number, hash string) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	var correct int
	err = tx.QueryRow(queries.CHECK_IF_HASH_CORRECT, countryCode, number, hash).Scan(&correct)
	if err != nil {
		return err
	}
	return nil
}

func (m *ApiModel) VerifyPin(countryCode, number, pin string) (string, error) {
	tx, err := m.DB.Begin()
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	var correct int
	err = tx.QueryRow(queries.CHECK_IF_PIN_CORRECT, countryCode, number, pin).Scan(&correct)
	if err != nil {
		return "", err
	}

	hashInput := fmt.Sprintf("%s%s%s", time.Now().Format("2006-01-02 15:04:05"), number, countryCode)
	h := sha1.New()
	h.Write([]byte(hashInput))
	sha1Hash := hex.EncodeToString(h.Sum(nil))

	_, err = mysequel.Update(mysequel.UpdateTable{
		Table: mysequel.Table{
			TableName: "user",
			Columns:   []string{"hash"},
			Vals:      []interface{}{sha1Hash},
			Tx:        tx,
		},
		WColumns: []string{"country_code", "number"},
		WVals:    []string{countryCode, number},
	})
	if err != nil {
		return "", err
	}

	return sha1Hash, nil
}

func (m *ApiModel) SignUp(countryCode, number, clockworkAPI, env string) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	var exists int
	err = tx.QueryRow(queries.CHECK_IF_NUMBER_EXISTS, countryCode, number).Scan(&exists)
	if err != nil {
		return err
	}

	otp, err := GenerateOTP(6)
	if err != nil {
		return err
	}

	if exists == 1 {
		_, err = mysequel.Update(mysequel.UpdateTable{
			Table: mysequel.Table{
				TableName: "user",
				Columns:   []string{"pin"},
				Vals:      []interface{}{otp},
				Tx:        tx,
			},
			WColumns: []string{"country_code", "number"},
			WVals:    []string{countryCode, number},
		})
	} else {
		_, err = mysequel.Insert(mysequel.Table{
			TableName: "user",
			Columns:   []string{"country_code", "number", "pin"},
			Vals:      []interface{}{countryCode, number, otp},
			Tx:        tx,
		})
	}
	if err != nil {
		return err
	}

	// resp, err := http.Get(fmt.Sprintf("http://www.textit.biz/sendmsg/?id=94768237192&pw=6200&to=%s%s&text=Your+Straddle+verification%20key+is+%s", countryCode, number, otp))
	if env == "prod" {
		resp, err := http.Get(fmt.Sprintf("https://api.clockworksms.com/http/send.aspx?key=%s&to=%s%s&content=Your+Straddle+verification+key+is+%s", clockworkAPI, countryCode, number, otp))
		if err != nil {
			return err
		}

		defer resp.Body.Close()
	}

	return nil
}

func GenerateOTP(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}
