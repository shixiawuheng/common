package lanzouw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"shixia/common"
	"strconv"
)

func Parse(uri, pwd string, fs *File, data *string) error {
	Header.Set("User-Agent", UA())
	Header.Set("Accept", "*/*")
	Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	Header.Set("Connection", "keep-alive")
	Header.Set("Sec-Fetch-Dest", "empty")
	Header.Set("Sec-Fetch-Mode", "cors")
	Header.Set("Sec-Fetch-Site", "same-origin")
	Header.Set("X-Requested-With", "XMLHttpRequest")
	Header.Set("X-Lanzou-Web-Request", "xhr")
	Header.Set("X-Lanzou-Web-Version", "2.0")
	Header.Set("X-Lanzou-Web-Dev", "web")
	Client = &http.Client{
		Jar: newCookieJar(),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	var param string
	{
		body, err := get(uri)
		if err != nil {
			return err
		}
		sign := common.GetBetweenStr(string(body), "var skdklds = '", "';")
		param = "action=downprocess&sign=" + sign + "&p=" + pwd
		//match, err := regexp.Compile("[^\\/]data : .*,")
		//if err != nil {
		//	return err
		//}
		//res := match.FindAll(body, -1)
		//for i := 0; i < len(res); i++ {
		//	param = common.GetBetweenStr(string(res[i]), "'", "'") + pwd
		//}
	}
	Header.Set("Referer", uri)
	{
		body, err := post(fmt.Sprintf("https://%s.lanzouf.com/ajaxm.php", common.GetBetweenStr(uri, "//", ".")), []byte(param))
		if err != nil {
			return err
		}
		err = json.Unmarshal(body, fs)
		if err != nil {
			return err
		}
		if fs.Zt != 1 {
			return fmt.Errorf(fs.Inf)
		}
	}
	Test := func(df *File, data *string) error {
		req, err := http.NewRequest("GET", fmt.Sprintf("%v/file/%v", df.Dom, df.Url), nil)
		if err != nil {
			return err
		}
		req.Header = Header.Clone()
		req.Header.Set("Referer", fmt.Sprintf("%v/f/%v", df.Dom, df.Url))
		res, err := Client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode != 302 {
			return fmt.Errorf("下载异常! 错误码: " + strconv.Itoa(res.StatusCode))
		}
		*data = res.Header.Get("Location")
		return nil
	}
	return Test(fs, data)
}

func UA() string {
	ua := map[int]string{
		0: "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.87 Safari/537.36",
		1: "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534.57.2 (KHTML, like Gecko) Version/5.1.7 Safari/534.57.2",
		2: "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:46.0) Gecko/20100101 Firefox/46.0",
	}

	return ua[rand.Intn(len(ua))]

}

var Header = http.Header{}

var Client *http.Client

func get(uri string) ([]byte, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header = Header.Clone()
	resp, err := Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func post(uri string, param []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", uri, bytes.NewReader(param))
	if err != nil {
		return nil, err
	}
	req.Header = Header.Clone()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

type File struct {
	Zt  int         `json:"zt"`
	Dom string      `json:"dom"`
	Url interface{} `json:"url"`
	Inf string      `json:"inf"`
}

func Download(uri string, file string, arge ...interface{}) error {
	root := filepath.Dir(file)
	if err := os.MkdirAll(root, 0755); err != nil {
		fmt.Println("创建目录失败:", err)
		return err
	}
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = f.Truncate(0); err != nil {
		fmt.Println("清空文件内容失败:", err)
		return err
	}
	if _, err = f.Seek(0, 0); err != nil {
		fmt.Println("定位文件读写指针失败:", err)
		return err
	}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}
	res, err := Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("下载失败! 错误码: " + strconv.Itoa(res.StatusCode))
	}
	size := res.ContentLength
	// 创建进度条
	bar := pb.Full.Start64(size)
	bar.Set(pb.Bytes, true)
	if len(arge) > 1 {
		bar.SetWriter(arge[0].(io.Writer))
	}
	reader := bar.NewProxyReader(res.Body)
	defer bar.Finish()
	_, err = io.Copy(f, reader)
	if err != nil {
		return err
	}
	return nil
}
