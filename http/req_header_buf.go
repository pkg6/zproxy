package http

import (
	"bytes"
	"fmt"
	"net"
	"net/url"
	"strings"
)

type ReqHeadBuf struct {
	Conn net.Conn
	//header内容
	HeadBuf []byte
	//请求方法
	Method string
	//请求域名
	hostOrURL string
	//处理之后地址
	Host string
	//请求域名或者在header中host
	URL string
}

func NewHeadBuf(inConn net.Conn, bufSize int) (req ReqHeadBuf, err error) {
	req.Conn = inConn
	buf := make([]byte, bufSize)
	len, err := req.Conn.Read(buf[:])
	if err != nil {
		return
	}
	req.HeadBuf = buf[:len]
	index := bytes.IndexByte(req.HeadBuf, '\n')
	_, err = fmt.Sscanf(string(req.HeadBuf[:index]), "%s%s", &req.Method, &req.hostOrURL)
	if err != nil {
		return
	}
	req.Method = strings.ToUpper(req.Method)
	if req.IsHTTPS() {
		req.HTTPS()
	} else {
		if err = req.HTTP(); err != nil {
			return
		}
	}
	return req, nil
}

func (req *ReqHeadBuf) HTTP() (err error) {
	if !strings.HasPrefix(req.hostOrURL, "/") {
		req.URL = req.hostOrURL
	}
	if req.URL == "" {
		_host, err := req.GetHeader("host")
		if err != nil {
			return err
		}
		req.URL = fmt.Sprintf("http://%s%s", _host, req.hostOrURL)
	}
	if err == nil {
		u, _ := url.Parse(req.URL)
		req.Host = u.Host
		req.addPortIfNot()
	}
	return
}
func (req *ReqHeadBuf) HTTPS() {
	req.Host = req.hostOrURL
	req.addPortIfNot()
	return
}
func (req *ReqHeadBuf) addPortIfNot() {
	port := "80"
	if req.IsHTTPS() {
		port = "443"
	}
	if (!strings.HasPrefix(req.Host, "[") && strings.Index(req.Host, ":") == -1) || (strings.HasPrefix(req.Host, "[") && strings.HasSuffix(req.Host, "]")) {
		req.Host = req.Host + ":" + port
	}
	return
}
func (req *ReqHeadBuf) IsHTTPS() bool {
	return req.Method == "CONNECT"
}

func (req *ReqHeadBuf) GetHeader(key string) (val string, err error) {
	key = strings.ToUpper(key)
	lines := strings.Split(string(req.HeadBuf), "\r\n")
	for _, line := range lines {
		line := strings.SplitN(strings.Trim(line, "\r\n "), ":", 2)
		if len(line) == 2 {
			k := strings.ToUpper(strings.Trim(line[0], " "))
			v := strings.Trim(line[1], " ")
			if key == k {
				val = v
				return
			}
		}
	}
	err = fmt.Errorf("can not find  header")
	return
}

func (req *ReqHeadBuf) Close() {
	if req.Conn != nil {
		_ = req.Conn.Close()
	}
}
