package mysql

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ssrdive/basara/pkg/sql/queries"
	"github.com/ssrdive/mysequel"
)

const otpChars = "1234567890"

type ApiModel struct {
	DB *sql.DB
}

func (m *ApiModel) SignUp(countryCode, number, clockworkAPI string) error {
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
	resp, err := http.Get(fmt.Sprintf("https://api.clockworksms.com/http/send.aspx?key=%s&to=%s%s&content=Your+Straddle+verification+key+is+%s", clockworkAPI, countryCode, number, otp))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(body)

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
