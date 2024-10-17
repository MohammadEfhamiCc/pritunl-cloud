package imds

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-cloud/agent/constants"
	"github.com/pritunl/pritunl-cloud/errortypes"
	"github.com/pritunl/pritunl-cloud/utils"
)

var (
	client = &http.Client{
		Timeout: 10 * time.Second,
	}
)

type Imds struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	Secret  string `json:"secret"`
}

func (m *Imds) Get(query string) (val string, err error) {
	u := url.URL{}
	u.Scheme = "http"
	u.Host = fmt.Sprintf("%s:%d", m.Address, m.Port)
	u.Path = "/query" + query

	req, e := http.NewRequest("GET", u.String(), nil)
	if e != nil {
		err = &errortypes.RequestError{
			errors.Wrap(e, "agent: Failed to create imds request"),
		}
		return
	}

	req.Header.Set("User-Agent", "pritunl-imds")
	req.Header.Set("Auth-Token", m.Secret)

	resp, e := client.Do(req)
	if e != nil {
		err = &errortypes.RequestError{
			errors.Wrap(e, "agent: Imds request failed"),
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body := ""
		data, _ := ioutil.ReadAll(resp.Body)
		if data != nil {
			body = string(data)
		}

		errData := &errortypes.ErrorData{}
		err = json.Unmarshal(data, errData)
		if err != nil || errData.Error == "" {
			errData = nil
		}

		if errData != nil && errData.Message != "" {
			body = errData.Message
		}

		err = &errortypes.RequestError{
			errors.Newf(
				"agent: Imds server error %d - %s",
				resp.StatusCode, body),
		}
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "agent: Imds failed to read body"),
		}
		return
	}

	val = string(data)

	return
}

func (m *Imds) Init() (err error) {
	confData, err := utils.Read(constants.ImdsConfPath)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(confData), m)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "agent: Failed to unmarshal imds conf"),
		}
		return
	}

	return
}
