package h2proxy

import (
	"crypto/tls"
	"fmt"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

var (
	addr = "localhost:3010"
	user = UserInfo{
		Username: "aaa",
		Passwd:   "bbb",
	}
)

func testServerStart(t *testing.T) {
	//InitLogger()
	ca, err := tls.X509KeyPair([]byte(crt), []byte(key))
	if err != nil {
		t.Fatal(err)
	}

	server := &http.Server{
		Addr: addr,
		TLSConfig: &tls.Config{
			Certificates:       []tls.Certificate{ca},
			InsecureSkipVerify: true,
			NextProtos:         []string{"h2", "h2c", "h2i"},
		},
		Handler: http.HandlerFunc(handle(&ServerConfig{
			NeedAuth: true,
			User:     &user,
		})),
	}

	// require cert.
	// generate cert for test:
	// openssl req -new -x509 -days 365 -key test1.key -out test1.crt
	t.Fatal(server.ListenAndServeTLS("", ""))
}

func TestServer(t *testing.T) {
	InitLogger()
	go testServerStart(t)

	tr := NewTransportWithProxy(addr)

	fmt.Println(Log)
	remoteAddr := "http://www.baidu.com/s?ie=utf-8&f=8&rsv_bp=1&rsv_idx=1&tn=baidu&wd=test&rsv_pq=b490c49a0000626b&rsv_t=132ffj2JcJlsvnHjGuDY6aR7woxPXQeCGImDWkR73XJBOuQrytnW9Racfew&rqlang=cn&rsv_enter=1&rsv_sug3=4&rsv_sug1=4&rsv_sug7=100&rsv_sug2=0&inputT=764&rsv_sug4=764"
	Log.Info(remoteAddr)

	req, err := http.NewRequest(
		http.MethodGet,
		remoteAddr,
		nil,
	)

	//req.Header = from.Header
	if err != nil {
		Log.Error("", zap.Error(err))
	}
	req.Header.Set("User-Agent", "curl/7.54.0")
	req.Header.Set("Accept", "*/*")
	SetAuthInHeader(&user, req)

	tr.DisableCompression = true
	resp, err := tr.RoundTrip(req)
	if err != nil {
		Log.Error(err)
		return
	}
	defer closeConn(resp.Body)

	if resp.StatusCode != 200 {
		Log.Info(resp.StatusCode)
		io.Copy(os.Stdout, resp.Body)
		Log.Info("Connect Proxy Server Error")
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b)[:10])
}

var key = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCoTcpQVpC4J5OH
GCEHYoPiUr9VVe7d4WoOBM9pSVMq9MZ6c+8KVGPO+e08iWDQeUO+6indtHLkuXVl
U4X23Jvq/lveNvIDoagm7iFELMyYsCgd2Vlo8mrlbE3Mx/Ii1/oUcatxCJpOwIjv
LsfXrKQrGR0mhr92QbiMwKYYhgN3m5FKE01IUV28LR+hb7qhxf19HSTQWsbS+Eg1
Vlej45Q+TpFgo1+AA4cBaNUwAz1rt+vOl+fXErO4929UGwdafei2liCDpbXDiQu3
2rqkLNH3rQ7Y4qa5ihbzKw+Z8tICTM96vFZ8wd2xYioHpjmQ4V714S726VI4raok
erZcsT1nAgMBAAECggEAMGvJBBoTotfICvr3La+7L9cMsxl5Ep7yqzvZDHqLEfXA
UKSgJkGnQXoINf10PCZXRksKZn2u/H0a+F7yUNahiRdLCQCX2lGdFi42pe2Zo2gy
8nuAzL/J312seHkiAcJPcahOYcPO2U9tVhsIZdWGmdus1lO4K+a2mNAoOS/9OWCx
9/PfJn47cX26JF0+k566JynhK4Dv4HdSlL6x/dgfS6LiSw5XwLdw/ovLiGf1vbNG
lFTMSlLhG47rS1wUt5BzM4GlsQQ7MN4Vs5db33XNGw5M6reBOwDHGzhe94711sEw
fFtBOPXyw194ThrwzMk0muulMsBZZJXhTu7bHbhLYQKBgQDZcuiYY1gv40ukdEnV
sJtbCNpU6nOR0pV6AexFmAzlTm9X2xPG0S3UOtMCA++cE9ijrSktxwRq/4FH/diI
A3h2qdgHbugxxeysfmWS6AlVUhV3zfWz91k1XopToImUdpWrxforfZSXfv2rOfL/
6jRt2pavw/TEWFIM6W8bNhH4cQKBgQDGJGWMmwotwYMOVugT4XCddYzmHnAFNDUy
3K3FPua61YCeeEKWrDYyNmIz8s0KNJIsJzKgD38v4rj7rqEaXs/3u4L6zuZTfb/o
3KftknzmOdLXevWxWYzyhmtmbLXkeA4ZTGgn+s15/WkoHkOz2bnjplioCL/mQ/iH
StHm4JI/VwKBgE9PINSLz1tP/IPTwiZFTrRqSy+Tf2ldNBWW4/USGwn7jJKvncvy
+VMhzVo700XK20X/XziKEOtxm1aFmFcrZOFq2xcC9X9J4COdyjBFnznWQWw723Sz
L39OpwcPU36prbdD8xWvrOWAdMbh0OZUJqE2i6U5xGlkiTCaZ2K2WuGRAoGBAIjD
zqRCz7/NdlyLeB1g2o6U+PBNyhyNcLruv7MKO9ByVhkMAUpnC/GUwCwDR6vnpY18
cOEyUSQIZo6ydtjw4LOqZjogXbL7dV+SDwdYuYVgHDxHzxbfLP6p8a/9EX/lrjWg
G7Sc1P+C/vaGDU0y17BevYsenvadrAoWhtPJ5qh5AoGAF5zIoLwqnKH2o0PEvGCk
2T7ZO4r1iUxF/6k3zSTexHK1oQLHO4fw39/B9C10B9zt4ASi9y7QfkPwPg7pJC/i
vS5eFawJM/bAsd3hTR3GFuPLlWJER+xRpLdowghBzrEpOpJ5OE8Zt/d2VrilF6Jg
nfRc0dbQZ8ny3t1noIQ/6qs=
-----END PRIVATE KEY-----`

var crt = `-----BEGIN CERTIFICATE-----
MIIDMjCCAhoCCQCAndLnW7yXOjANBgkqhkiG9w0BAQsFADBbMQswCQYDVQQGEwIx
MTEKMAgGA1UECAwBMjEKMAgGA1UEBwwBMTEKMAgGA1UECgwBMTEKMAgGA1UECwwB
MTEKMAgGA1UEAwwBMTEQMA4GCSqGSIb3DQEJARYBMTAeFw0xOTAzMjQxNTM2NTNa
Fw0yMDAzMjMxNTM2NTNaMFsxCzAJBgNVBAYTAjExMQowCAYDVQQIDAEyMQowCAYD
VQQHDAExMQowCAYDVQQKDAExMQowCAYDVQQLDAExMQowCAYDVQQDDAExMRAwDgYJ
KoZIhvcNAQkBFgExMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqE3K
UFaQuCeThxghB2KD4lK/VVXu3eFqDgTPaUlTKvTGenPvClRjzvntPIlg0HlDvuop
3bRy5Ll1ZVOF9tyb6v5b3jbyA6GoJu4hRCzMmLAoHdlZaPJq5WxNzMfyItf6FHGr
cQiaTsCI7y7H16ykKxkdJoa/dkG4jMCmGIYDd5uRShNNSFFdvC0foW+6ocX9fR0k
0FrG0vhINVZXo+OUPk6RYKNfgAOHAWjVMAM9a7frzpfn1xKzuPdvVBsHWn3otpYg
g6W1w4kLt9q6pCzR960O2OKmuYoW8ysPmfLSAkzPerxWfMHdsWIqB6Y5kOFe9eEu
9ulSOK2qJHq2XLE9ZwIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQBkwreMyqwrAxwo
DXr5y/7MYIsCEwTuuMuN28NB0kqEGtTc5rUVl2A89FAwL9LG/qknLC0MYVYqolMB
ZUzPcDhb5zZOJy91wSLO0QkZ3Ou8BE70k1jPqgCX5FlbmgLDpk9Esm8P99OvkCpJ
+8G1l+24JYwnskkNDu8mQFL8ZZEG3rXNbgE2fCXc0p9GtmNEcMFiCYe6WgwvGgg6
TQl3GmL13n1MooSzyvHZOfOfVHufZe1JDZyApsfUxCE+DNpeDmZhP/k24rlL+xxy
UlmSMAR8lmZoc4voVh2/EnaQiBd7+46kEGLEqz/qB06HbNrs9mqMYxO6UbdE0qbH
sgGLrMCt
-----END CERTIFICATE-----`
