package pgsql

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

func GenerateDbConnectAdminAuthToken(parsedUrl *url.URL) (string, error) {
	// Fetch credentials
	sess, err := session.NewSession()
	if err != nil {
		return "", err
	}

	creds, err := sess.Config.Credentials.Get()
	if err != nil {
		return "", err
	}
	staticCredentials := credentials.NewStaticCredentials(
		creds.AccessKeyID,
		creds.SecretAccessKey,
		creds.SessionToken,
	)

	region := strings.Split(parsedUrl.Hostname(), ".")[2]

	// The scheme is arbitrary and is only needed because validation of the URL requires one.
	endpoint := "https://" + parsedUrl.Hostname()
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}
	values := req.URL.Query()
	values.Set("Action", "DbConnectAdmin")
	req.URL.RawQuery = values.Encode()

	signer := v4.Signer{
		Credentials: staticCredentials,
	}
	_, err = signer.Presign(req, nil, "dsql", region, 60*time.Minute, time.Now())
	if err != nil {
		return "", err
	}

	url := req.URL.String()[len("https://"):]

	return url, nil
}
