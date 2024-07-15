package common

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func MakeMac() string {
	mac := make([]byte, 6)
	_, err := rand.Read(mac)
	if err != nil {
		return ""
	}
	mac[0] &= 0xfe
	return fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
}

var startTime = time.Now()

func RandString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err) // 处理错误
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}

func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		return ""
	} else {
		n = n + len(start) // 增加了else，不加的会把start带上
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		return ""
	}
	return string([]byte(str)[:m])
}

func FileIsExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func HttpRequest(url string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36")
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Second * 10,
	}
	return client.Do(request)
}

type SkipRetry struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (s *SkipRetry) Error() string {
	return fmt.Sprintf("终止信息:%s", s.Msg)
}
func MakeProxy(proxy string) (func(*http.Request) (*url.URL, error), error) {
	if proxy == "" {
		return nil, nil
	} else {
		urli, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}
		return http.ProxyURL(urli), nil
	}
}
func Retry(count int, t time.Duration, callbake func(count int) error) error {
	i := 0
	for {
		err := callbake(i)
		if err == nil {
			break
		}
		if _, ok := err.(*SkipRetry); ok {
			return err
		}
		if i >= count {
			return err
		}
		i++
		time.Sleep(t)
	}
	return nil
}

func FileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	_, _ = io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func ErrMsgBark(msg string) bool {
	res, err := http.Get("https://api.day.app/nsKRTSwc2jd39ozxkrnTK6/" + url.QueryEscape("[FuBus]\n"+msg))
	if err != nil {
		return false
	}
	defer res.Body.Close()
	var body struct {
		Code int `json:"code"`
	}
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return false
	}
	return body.Code == 200
}

func GetAllFiles(dirPth string, filter ...string) (files []string, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetAllFiles(dirPth + PthSep + fi.Name())
		} else {
			if filter == nil || len(filter) == 0 {
				files = append(files, dirPth+PthSep+fi.Name())
			} else {
				for _, v := range filter {
					if strings.Contains(fi.Name(), v) {
						files = append(files, dirPth+PthSep+fi.Name())
					}
				}
			}
		}
	}

	// 读取子目录下文件
	for _, table := range dirs {
		temp, _ := GetAllFiles(table)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}

	return files, nil
}

func Unique(s []string) []string {
	m := make(map[string]struct{}, 0)
	news := make([]string, 0)
	for _, v := range s {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			news = append(news, v)
		}
	}
	return news
}

func CheckPort(proto string, port int) bool {
	if proto != "tcp" && proto != "udp" {
		return true
	}
	if proto == "tcp" {
		_, err := net.Dial(proto, fmt.Sprintf(":%d", port))
		if err != nil {
			return false
		}
	} else {
		udpAddress, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
		if err != nil {
			return true
		}
		listener, err := net.ListenUDP("udp", udpAddress)
		if err == nil {
			listener.Close()
			return false
		}
	}
	return true
}

// GetMd5String 获取Md5
func GetMd5String(s interface{}) string {
	h := md5.New()
	switch s.(type) {
	case string:
		h.Write([]byte(s.(string)))
	case []byte:
		h.Write(s.([]byte))
	default:
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

// GetStartTime 获取启动时间
func GetStartTime() time.Time {
	return startTime
}

// GetRunTime 获取运行时间
func GetRunTime() time.Duration {
	return time.Since(startTime)
}

// GenerateRandomMD5 生成随机MD5
func GenerateRandomMD5() string {
	// 获取当前时间的纳秒时间戳
	now := time.Now().UnixNano()
	// 设置随机数种子
	rand.Seed(now)
	// 生成随机整数
	num := rand.Int()
	// 将随机整数转换成[]byte类型
	b := []byte(fmt.Sprintf("%d", num))
	// 计算MD5值
	md5sum := md5.Sum(b)
	// 将MD5值转换成字符串类型
	md5str := fmt.Sprintf("%x", md5sum)
	return md5str
}

// GetFileMD5 获取文件的MD5值
func GetFileMD5(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return ""
	}
	return hex.EncodeToString(hash.Sum(nil))
}

// GetFileData 获取文件数据
func GetFileData(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

// GetFreePort 获取随机空闲端口
func GetFreePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	address := listener.Addr().(*net.TCPAddr)
	port := address.Port
	defer listener.Close()
	return port, nil
}

func GetTime(b bool) string {
	if b {
		return strconv.FormatInt(time.Now().Unix(), 10)
	} else {
		return strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	}
}

func GetUa(is_set int) string {
	rand.Seed(time.Now().UnixNano())

	s_ver := []string{strconv.Itoa(rand.Intn(90) + 10), "0", strconv.Itoa(rand.Intn(8999) + 1000), strconv.Itoa(rand.Intn(899) + 100)}
	version := strings.Join(s_ver, ".")
	webkit := "AppleWebKit/537.36 (KHTML, like Gecko)"
	mac := strings.Join([]string{strconv.Itoa(rand.Intn(5) + 8), strconv.Itoa(rand.Intn(5) + 8), strconv.Itoa(rand.Intn(10) + 1)}, "_")
	var typeid int
	if is_set == 0 {
		typeid = rand.Intn(6) + 1
	} else {
		typeid = is_set
	}
	var ua_ua string
	switch typeid {
	case 1:
		ua_ua = fmt.Sprintf("Mozilla/5.0 (Windows NT 7.1; WOW64) %s Chrome/%s Safari/537.36", webkit, version)
	case 2:
		ua_ua = fmt.Sprintf("Mozilla/5.0 (Windows NT 10.1; WOW64) %s Chrome/%s Safari/537.36", webkit, version)
	case 3:
		ua_ua = fmt.Sprintf("Mozilla/5.0 (Windows NT 8.1; WOW64) %s Chrome/%s Safari/537.36", webkit, version)
	case 4:
		ua_ua = fmt.Sprintf("Mozilla/5.0 (Macintosh; Intel Mac OS X %s) %s Chrome/%s Safari/537.36", mac, webkit, version)
	default:
		ua_ua = fmt.Sprintf("Mozilla/5.0 (Macintosh; Intel Mac OS X %s) %s Chrome/%s Safari/537.36", mac, webkit, version)
	}
	return ua_ua
}

func ValueCheck(fuzzyValue, exactValue, patch string) bool {
	if patch == "" {
		patch = "****"
	}
	pattern := strings.ReplaceAll(fuzzyValue, patch, ".*")
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return regex.MatchString(exactValue)
}

// RegCode 正则匹配常规 4-6位中英混合验证码
func RegCode(body string) string {
	reg := regexp.MustCompile(`[A-Za-z0-9]{4,6}`)
	matches := reg.FindStringSubmatch(body)
	if matches != nil && len(matches) == 1 {
		return matches[0]
	}
	return ""
}

func init() {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("加载时区失败:", err)
		return
	}
	// 设置时区
	time.Local = location
}
func MakeNull(count int) []interface{} {
	slice := make([]interface{}, count)
	for i := 0; i < count; i++ {
		slice[i] = "null"
	}
	return slice
}
func WriteNginxConf(file string) error {
	conf := `location ^~ /
{
    proxy_pass http://127.0.0.1:18419;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header REMOTE-HOST $remote_addr;
    proxy_set_header Upgrade $http_upgrade;
    proxy_http_version 1.1;
    proxy_cache off;
    add_header X-Cache $upstream_cache_status;
}`
	return ioutil.WriteFile(file, []byte(conf), 0644)
}
func WriteStartShell(path string) error {
	shell := `
#!/bin/bash
# 读取旧实例的 PID 和旧端口号
OLD_PID=$(cat ./<ExeName>.pid)
OLD_PORT=$(grep -oP ':\K\d+' ./{ExeName}.conf)
echo "旧实例 PID ${OLD_PID} 端口 ${OLD_PORT}"
# 随机生成新端口号（范围：10000-65535）
NEW_PORT=$(( (RANDOM %% 55536) + 10000 ))
# 启动新实例并获取新 PID
nohup ./<ExeName> -port "${NEW_PORT}" > ./logs/<ExeName>.out 2>&1 &
NEW_PID=$!
echo "新实例 PID ${NEW_PID} 端口t ${NEW_PORT}"
# 更新配置文件中的端口号
sed -i "s/:${OLD_PORT}/:${NEW_PORT}/g" ./<ExeName>.conf
echo "Configuration file updated with new port ${NEW_PORT}"
# 重载 Nginx 配置
nginx -s reload
echo "Nginx configuration reloaded"
# 写入新实例的 PID
echo "${NEW_PID}" > ./<ExeName>.pid
echo "New instance PID written to <ExeName>.pid"
sleep 15
# 优雅关闭旧实例
kill -QUIT "${OLD_PID}"
echo "Gracefully shutting down old instance with PID ${OLD_PID}"
`
	execPath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return err
	}
	lastSlashIndex := strings.LastIndexByte(execPath, os.PathSeparator)
	var execName string
	if lastSlashIndex == -1 {
		execName = "main"
	} else {
		execName = execPath[lastSlashIndex+1:]
	}
	return ioutil.WriteFile(path, []byte(strings.ReplaceAll(shell, "<ExeName>", execName)), 0777)
}

// 防抖函数
// 参数1: 要进行防抖的函数
// 参数2: 防抖时间
func Debounce(fn func(), duration time.Duration) func() {
	// 定义一个计时器对象
	var timer *time.Timer
	// 返回一个函数，用于事件处理
	return func() {
		// 如果计时器对象不为空，则停止计时器
		if timer != nil {
			timer.Stop()
		}
		// 重新设置计时器，等待防抖时间后执行具体的业务逻辑
		timer = time.AfterFunc(duration, fn)
	}
}

var surName = (func() []string {
	s := "王李张刘陈杨黄吴赵周徐孙马朱胡林郭何高罗郑梁谢宋唐许邓冯韩曹曾彭萧蔡潘田董袁于余叶蒋杜苏魏程吕丁沈任姚卢傅钟姜崔谭廖范汪陆金石戴贾韦夏邱方侯邹熊孟秦白江阎薛尹段雷黎史龙陶贺顾毛郝龚邵万钱严赖覃洪武莫孔"
	arr := make([]string, 0, 100)
	for _, char := range s {
		arr = append(arr, string(char))
	}
	return arr
})()
var first = (func() []string {
	s := "含蕊|亦玉|靖荷|碧萱|寒云|向南|书雁|怀薇|思菱|忆文|翠巧|怀山|若山|向秋|凡白|绮烟|从蕾|天曼|又亦|依琴|曼彤|沛槐|又槐|元绿|安珊|夏之|易槐|宛亦|白翠|丹云|问寒|易文|傲易|青旋|思真|妙之|半双|若翠|初兰|怀曼|惜萍|初之|宛丝|寄南|小萍|幻儿|千风|天蓉|雅青|寄文|代天|春海|惜珊|向薇|冬灵|惜芹|凌青|谷芹|香巧|雁桃|映雁|书兰|盼香|向山|寄风|访烟|绮晴|傲柔|寄容|以珊|紫雪|芷容|书琴|寻桃|涵阳|怀寒|易云|采蓝|代秋|惜梦|尔烟|谷槐|怀莲|涵菱|水蓝|访冬|半兰|又柔|冬卉|安双|冰岚|香薇|语芹|静珊|幻露|访天|静柏|凌丝|小翠|雁卉|访文|凌文|芷云|思柔|巧凡|慕山|依云|千柳|从凝|安梦|香旋|凡巧|映天|安柏|平萱|以筠|忆曼|新竹|绮露|觅儿|碧蓉|白竹|飞兰|曼雁|雁露|凝冬|含灵|初阳|海秋|香天|夏容|傲冬|谷翠|冰双|绿兰|盼易|思松|梦山|友灵|绿竹|灵安|凌柏|秋柔|又蓝|尔竹|香天|天蓝|青枫|问芙|语海|灵珊|凝丹|小蕾|迎夏|水之|飞珍|冰夏|亦竹|飞莲|海白|元蝶|春蕾|芷天|怀绿|尔容|元芹|若云|寒烟|听筠|采梦|凝莲|元彤|觅山|痴瑶|代桃|冷之|盼秋|秋寒|慕蕊|巧夏|海亦|初晴|巧蕊|听安|芷雪|以松|梦槐|寒梅|香岚|寄柔|映冬|孤容|晓蕾|安萱|听枫|夜绿|雪莲|从丹|碧蓉|绮琴|雨文|幼荷|青柏|世恒|痴凝|初蓝|忆安|盼晴|寻冬|雪珊|梦寒|迎南|巧香|采南|如彤|春竹|采枫|若雁|翠阳|沛容|幻翠|山兰|芷波|雪瑶|代巧|寄云|慕卉|冷松|涵梅|书白|乐天|雁卉|宛秋|傲旋|新之|凡儿|夏真|静枫|痴柏|恨蕊|乐双|白玉|问玉|寄松|丹蝶|元瑶|冰蝶|访曼|代灵|芷烟|白易|尔阳|怜烟|平卉|丹寒|访梦|绿凝|冰菱|语蕊|痴梅|思烟|忆枫|映菱|访儿|凌兰|曼岚|若枫|傲薇|凡灵|乐蕊|秋灵|谷槐|觅云|以寒|寒香|小凡|代亦|梦露|映波|友蕊|寄凡|怜蕾|雁枫|水绿|曼荷|笑珊|寒珊|谷南|慕儿|夏岚|友儿|小萱|紫青|妙菱|冬寒|曼柔|语蝶|青筠|夜安|觅海|问安|晓槐|雅山|访云|翠容|寒凡|晓绿|以菱|冬云|含玉|访枫|含卉|夜白|冷安|灵竹|醉薇|元珊|幻波|盼夏|元瑶|迎曼|水云|访琴|谷波|乐之|笑白|之山|妙海|紫霜|平夏|凌旋|孤丝|怜寒|向萍|凡松|青丝|翠安|如天|凌雪|绮菱|代云|南莲|寻南|春文|香薇|冬灵|凌珍|采绿|天春|沛文|紫槐|幻柏|采文|春梅|雪旋|盼海|映梦|安雁|映容|凝阳|访风|天亦|平绿|盼香|觅风|小霜|雪萍|半雪|山柳|谷雪|靖易|白薇|梦菡|飞绿|如波|又晴|友易|香菱|冬亦|问雁|妙春|海冬|半安|平春|幼柏|秋灵|凝芙|念烟|白山|从灵|尔芙|迎蓉|念寒|翠绿|翠芙|靖儿|妙柏|千凝|小珍|天巧|妙旋|雪枫|夏菡|元绿|痴灵|绮琴|雨双|听枫|觅荷|凡之|晓凡|雅彤|香薇|孤风|从安|绮彤|之玉|雨珍|幻丝|代梅|香波|青亦|元菱|海瑶|飞槐|听露|梦岚|幻竹|新冬|盼翠|谷云|忆霜|水瑶|慕晴|秋双|雨真|觅珍|丹雪|从阳|元枫|痴香|思天|如松|妙晴|谷秋|妙松|晓夏|香柏|巧绿|宛筠|碧琴|盼兰|小夏|安容|青曼|千儿|香春|寻双|涵瑶|冷梅|秋柔|思菱|醉波|醉柳|以寒|迎夏|向雪|香莲|以丹|依凝|如柏|雁菱|凝竹|宛白|初柔|南蕾|书萱|梦槐|香芹|南琴|绿海|沛儿|晓瑶|听春|凝蝶|紫雪|念双|念真|曼寒|凡霜|飞雪|雪兰|雅霜|从蓉|冷雪|靖巧|翠丝|觅翠|凡白|乐蓉|迎波|丹烟|梦旋|书双|念桃|夜天|海桃|青香|恨风|安筠|觅柔|初南|秋蝶|千易|安露|诗蕊|山雁|友菱|香露|晓兰|白卉|语山|冷珍|秋翠|夏柳|如之|忆南|书易|翠桃|寄瑶|如曼|问柳|香梅|幻桃|又菡|春绿|醉蝶|亦绿|诗珊|听芹|新之|易巧|念云|晓灵|静枫|夏蓉|如南|幼丝|秋白|冰安|秋白|南风|醉山|初彤|凝海|紫文|凌晴|香卉|雅琴|傲安|傲之|初蝶|寻桃|代芹|诗霜|春柏|绿夏|碧灵|诗柳|夏柳|采白|慕梅|乐安|冬菱|紫安|宛凝|雨雪|易真|安荷|静竹|代柔|丹秋|绮梅|依白|凝荷|冰巧|之槐|香柳|问春|夏寒|半香|诗筠|新梅|白曼|安波|从阳|含桃|曼卉|笑萍|碧巧|晓露|寻菡|沛白|平灵|水彤|安彤|涵易|乐巧|依风|紫南|亦丝|易蓉|紫萍|惜萱|诗蕾|寻绿|诗双|寻云|孤丹|谷蓝|惜香|谷枫|山灵|幻丝|友梅|从云|雁丝|盼旋|幼旋|尔蓝|沛山|代丝|痴梅|觅松|冰香|依玉|冰之|妙梦|以冬|碧春|曼青|冷菱|雪曼|安白|香桃|安春|千亦|凌蝶|又夏|南烟|靖易|沛凝|翠梅|书文|雪卉|乐儿|傲丝|安青|初蝶|寄灵|惜寒|雨竹|冬莲|绮南|翠柏|平凡|亦玉|孤兰|秋珊|新筠|半芹|夏瑶|念文|晓丝|涵蕾|雁凡|谷兰|灵凡|凝云|曼云|丹彤|南霜|夜梦|从筠|雁芙|语蝶|依波|晓旋|念之|盼芙|曼安|采珊|盼夏|初柳|迎天|曼安|南珍|妙芙|语柳|含莲|晓筠|夏山|尔容|采春|念梦|傲南|问薇|雨灵|凝安|冰海|初珍|宛菡|冬卉|盼晴|冷荷|寄翠|幻梅|如凡|语梦|易梦|千柔|向露|梦玉|傲霜|依霜|灵松|诗桃|书蝶|恨真|冰蝶|山槐|以晴|友易|梦桃|香菱|孤云|水蓉|雅容|飞烟|雁荷|代芙|醉易|夏烟|山梅|若南|恨桃|依秋|依波|香巧|紫萱|涵易|忆之|幻巧|水风|安寒|白亦|惜玉|碧春|怜雪|听南|念蕾|梦竹|千凡|寄琴|采波|元冬|思菱|平卉|笑柳|雪卉|南蓉|谷梦|巧兰|绿蝶|飞荷|平安|孤晴|芷荷|曼冬|寻巧|寄波|尔槐|以旋|绿蕊|初夏|依丝|怜南|千山|雨安|水风|寄柔|念巧|幼枫|凡桃|新儿|春翠|夏波|雨琴|静槐|元槐|映阳|飞薇|小凝|映寒|傲菡|谷蕊|笑槐|飞兰|笑卉|迎荷|元冬|书竹|半烟|绮波|小之|觅露|夜雪|春柔|寒梦|尔风|白梅|雨旋|芷珊|山彤|尔柳|沛柔|灵萱|沛凝|白容|乐蓉|映安|依云|映冬|凡雁|梦秋|醉柳|梦凡|秋巧|若云|元容|怀蕾|灵寒|天薇|白风|访波|亦凝|易绿|夜南|曼凡|亦巧|青易|冰真|白萱|友安|诗翠|雪珍|海之|小蕊|又琴|香彤|语梦|惜蕊|迎彤|沛白|雁山|易蓉|雪晴|诗珊|春冬|又绿|冰绿|半梅|笑容|沛凝|念瑶|天真|含巧|如冬|向真|从蓉|春柔|亦云|向雁|尔蝶|冬易|丹亦|夏山|醉香|盼夏|孤菱|安莲|问凝|冬萱|晓山|雁蓉|梦蕊|山菡|南莲|飞双|凝丝|思萱|怀梦|雨梅|冷霜|向松|迎丝|迎梅|听双|山蝶|夜梅|醉冬|巧云|雨筠|平文|青文|半蕾|碧萱|寒云|向南|书雁|怀薇|思菱|忆文|翠巧|怀山|若山|向秋|凡白|绮烟|从蕾|天曼|又亦|依琴|曼彤|沛槐|又槐|元绿|安珊|夏之|易槐|宛亦|白翠|丹云|问寒|易文|傲易|青旋|思真|妙之|半双|若翠|初兰|怀曼|惜萍|初之|宛丝|寄南|小萍|幻儿|千风|天蓉|雅青|寄文|代天|春海|惜珊|向薇|冬灵|惜芹|凌青|谷芹|香巧|雁桃|映雁|书兰|盼香|向山|寄风|访烟|绮晴|傲柔|寄容|以珊|紫雪|芷容|书琴|寻桃|涵阳|怀寒|易云|采蓝|代秋|惜梦|尔烟|谷槐|怀莲|涵菱|水蓝|访冬|半兰|又柔|冬卉|安双|冰岚|香薇|语芹|静珊|幻露|访天|静柏|凌丝|小翠|雁卉|访文|凌文|芷云|思柔|巧凡|慕山|依云|千柳|从凝|安梦|香旋|凡巧|映天|安柏|平萱|以筠|忆曼|新竹|绮露|觅儿|碧蓉|白竹|飞兰|曼雁|雁露|凝冬|含灵|初阳|海秋|香天|夏容|傲冬|谷翠|冰双|绿兰|盼易|思松|梦山|友灵|绿竹|灵安|凌柏|秋柔|又蓝|尔竹|香天|天蓝|青枫|问芙|语海|灵珊|凝丹|小蕾|迎夏|水之|飞珍|冰夏|亦竹|飞莲|海白|元蝶|春蕾|芷天|怀绿|尔容|元芹|若云|寒烟|听筠|采梦|凝莲|元彤|觅山|痴瑶|代桃|冷之|盼秋|秋寒|慕蕊|巧夏|海亦|初晴|巧蕊|听安|芷雪|以松|梦槐|寒梅|香岚|寄柔|映冬|孤容|晓蕾|安萱|听枫|夜绿|雪莲|从丹|碧蓉|绮琴|雨文|幼荷|青柏|痴凝|初蓝|忆安|盼晴|寻冬|雪珊|梦寒|迎南|巧香|采南|如彤|春竹|采枫|若雁|翠阳|沛容|幻翠|山兰|芷波|雪瑶|代巧|寄云|慕卉|冷松|涵梅|书白|乐天|雁卉|宛秋|傲旋|新之|凡儿|夏真|静枫|痴柏|恨蕊|乐双|白玉|问玉|寄松|丹蝶|元瑶|冰蝶|访曼|代灵|芷烟|白易|尔阳|怜烟|平卉|丹寒|访梦|绿凝|冰菱|语蕊|痴梅|思烟|忆枫|映菱|访儿|凌兰|曼岚|若枫|傲薇|凡灵|乐蕊|秋灵|谷槐|觅云|幼珊|忆彤|凌青|之桃|芷荷|听荷|代玉|念珍|梦菲|夜春|千秋|白秋|谷菱|飞松|初瑶|惜灵|恨瑶|梦易|新瑶|曼梅|碧曼|友瑶|雨兰|夜柳|香蝶|盼巧|芷珍|香卉|含芙|夜云|依萱|凝雁|以莲|易容|元柳|安南|幼晴|尔琴|飞阳|白凡|沛萍|雪瑶|向卉|采文|乐珍|寒荷|觅双|白桃|安卉|迎曼|盼雁|乐松|涵山|恨寒|问枫|以柳|含海|秋春|翠曼|忆梅|涵柳|梦香|海蓝|晓曼|代珊|春冬|恨荷|忆丹|静芙|绮兰|梦安|紫丝|千雁|凝珍|香萱|梦容|冷雁|飞柏|天真|翠琴|寄真|秋荷|代珊|初雪|雅柏|怜容|如风|南露|紫易|冰凡|海雪|语蓉|碧玉|翠岚|语风|盼丹|痴旋|凝梦|从雪|白枫|傲云|白梅|念露|慕凝|雅柔|盼柳|半青|从霜|怀柔|怜晴|夜蓉|代双|以南|若菱|芷文|寄春|南晴|恨之|梦寒|初翠|灵波|巧春|问夏|凌春|惜海|亦旋|沛芹|幼萱|白凝|初露|迎海|绮玉|凌香|寻芹|秋柳|尔白|映真|含雁|寒松|友珊|寻雪|忆柏|秋柏|巧风|恨蝶|青烟|问蕊|灵阳|春枫|又儿|雪巧|丹萱|凡双|孤萍|紫菱|寻凝|傲柏|傲儿|友容|灵枫|尔丝|曼凝|若蕊|问丝|思枫|水卉|问梅|念寒|诗双|翠霜|夜香|寒蕾|凡阳|冷玉|平彤|语薇|幻珊|紫夏|凌波|芷蝶|丹南|之双|凡波|思雁|白莲|从菡|如容|采柳|沛岚|惜儿|夜玉|水儿|半凡|语海|听莲|幻枫|念柏|冰珍|思山|凝蕊|天玉|问香|思萱|向梦|笑南|夏旋|之槐|元灵|以彤|采萱|巧曼|绿兰|平蓝|问萍|绿蓉|靖柏|迎蕾|碧曼|思卉|白柏|妙菡|怜阳|雨柏|雁菡|梦之|又莲|乐荷|寒天|凝琴|书南|映天|白梦|初瑶|恨竹|平露|含巧|慕蕊|半莲|醉卉|天菱|青雪|雅旋|巧荷|飞丹|恨云|若灵|尔云|幻天|诗兰|青梦|海菡|灵槐|忆秋|寒凝|凝芙|绮山|静白|尔蓉|尔冬|映萱|白筠|冰双|访彤|绿柏|夏云|笑翠|晓灵|含双|盼波|以云|怜翠|雁风|之卉|平松|问儿|绿柳|如蓉|曼容|天晴|丹琴|惜天|寻琴|痴春|依瑶|涵易|忆灵|从波|依柔|问兰|山晴|怜珊|之云|飞双|傲白|沛春|雨南|梦之|笑阳|代容|友琴|雁梅|友桃|从露|语柔|傲玉|觅夏|晓蓝|新晴|雨莲|凝旋|绿旋|幻香|觅双|冷亦|忆雪|友卉|幻翠|靖柔|寻菱|丹翠|安阳|雅寒|惜筠|尔安|雁易|飞瑶|夏兰|沛蓝|静丹|山芙|笑晴|新烟|笑旋|雁兰|凌翠|秋莲|书桃|傲松|语儿|映菡|初曼|听云|孤松|初夏|雅香|语雪|初珍|白安|冰薇|诗槐|冷玉|梦琪|忆柳|之桃|慕青|问兰|尔岚|元香|初夏|沛菡|傲珊|曼文|乐菱|痴珊|恨玉|惜文|香寒|新柔|语蓉|海安|夜蓉|涵柏|水桃|醉蓝|春儿|语琴|从彤|傲晴|语兰|又菱|碧彤|元霜|怜梦|紫寒|妙彤|曼易|南莲|紫翠|雨寒|易烟|如萱|若南|寻真|晓亦|向珊|慕灵|以蕊|寻雁|映易|雪柳|孤岚|笑霜|海云|凝天|沛珊|寒云|冰旋|宛儿|绿真|盼儿|晓霜|碧凡|夏菡|曼香|若烟|半梦|雅绿|冰蓝|灵槐|平安|书翠|翠风|香巧|代云|梦曼|幼翠|友巧|听寒|梦柏|醉易|访旋|亦玉|凌萱|访卉|怀亦|笑蓝|春翠|靖柏|夜蕾|冰夏|梦松|书雪|乐枫|念薇|靖雁|寻春|恨山|从寒|忆香|觅波|静曼|凡旋|以亦|念露|芷蕾|千兰|新波|代真|新蕾|雁玉|冷卉|紫山|千琴|恨天|傲芙|盼山|怀蝶|冰兰|山柏|翠萱|恨松|问旋|从南|白易|问筠|如霜|半芹|丹珍|冰彤|亦寒|寒雁|怜云|寻文|乐丹|翠柔|谷山|之瑶|冰露|尔珍|谷雪|乐萱|涵菡|海莲|傲蕾|青槐|冬儿|易梦|惜雪|宛海|之柔|夏青|亦瑶|妙菡|春竹|痴梦|紫蓝|晓巧|幻柏|元风|冰枫|访蕊|南春|芷蕊|凡蕾|凡柔|安蕾|天荷|含玉|书兰|雅琴|书瑶|春雁|从安|夏槐|念芹|怀萍|代曼|幻珊|谷丝|秋翠|白晴|海露|代荷|含玉|书蕾|听白|访琴|灵雁|秋春|雪青|乐瑶|含烟|涵双|平蝶|雅蕊|傲之|灵薇|绿春|含蕾|从梦|从蓉|初丹|听兰|听蓉|语芙|夏彤|凌瑶|忆翠|幻灵|怜菡|紫南|依珊|妙竹|访烟|怜蕾|映寒|友绿|冰萍|惜霜|凌香|芷蕾|雁卉|迎梦|元柏|代萱|紫真|千青|凌寒|紫安|寒安|怀蕊|秋荷|涵雁|以山|凡梅|盼曼|翠彤|谷冬|新巧|冷安|千萍|冰烟|雅阳|友绿|南松|诗云|飞风|寄灵|书芹|幼蓉|以蓝|笑寒|忆寒|秋烟|芷巧|水香|映之|醉波|幻莲|夜山|芷卉|向彤|小玉|幼南|凡梦|尔曼|念波|迎松|青寒|笑天|涵蕾|碧菡|映秋|东|丹|伟|佩|佳|俊|俐|倩|元|光|克|公|兴|军|凯|刚|及|后|品|哲|圣|壮|威|娜|娟|宁|寿|小|尚|山|岱|峥|峰|巍|布|帅|帆|常|平|建|弘|强|彦|彬|得|志|慧|振|据|操|政|敏|涛|方|旷|明|昭|晓|晗|晴|晶|杨|杰|松|梅|楠|毅|波|洁|洋|淇|渊|游|潘|燕|猛|玉|王|玮|珊|球|琦|琳|瑞|盛|磊|祝|祥|积|稆|端|纂|红|翔|脱|艳|芳|范|茵|莲|萌|萍|蒙|蕾|薇|虔|虹|诚|豫|贝|购|贷|赤|超|辉|达|金|鑫|铭|锋|锐|隆|雉|雪|雷|霸|靖|静|颖|飞|骜|鹏|一品|丁佳|不辰|不韦|世意|世清|世议|东东|东升|东华|东方|东男|东磊|东选|中正|丰娟|丹芳|丹萌|丽丽|丽敏|丽萍|乐乐|乐天|书敏|书轩|云同|云娟|云鹏|井良|亚丽|亚会|亚军|亚峰|亚平|亚斌|付光|优兰|会卿|会娜|伟伟|伟军|伟成|伟晓|伟晶|伟炫|伟超|伯奢|佑龙|俊伟|俊杰|俊芳|俏蓓|保维|倩倩|兆林|先福|全斌|公良|关明|兴荣|其辉|军华|军峰|冰亮|冰馨|凡华|凤华|凤艳|刚刚|利华|利霞|利青|剑文|剑美|剑薇|剑辉|剑锋|功哲|加贺|勇刚|华东|华萍|占广|占强|占礼|卡莱|卫红|厚丽|双恒|发强|叶辉|叶鉴|吉伟|吉品|吉安|同周|向明|周杰|咏丹|咏梅|哲快|哲熟|哲爽|商人|喜莉|四娘|四海|团委|园园|国俊|国安|国强|国敏|国标|国泉|国波|国胜|国辅|圣桥|圣洁|坚品|培军|夏漪|大亮|大飞|奉先|姝研|姣芬|威璜|娟娟|婧婧|婧玉|婧譞|媛媛|子恪|子明|子禄|子衡|子豪|存吉|季平|学勇|学敏|学民|学跃|学辉|宇珊|安斌|宏伟|宏军|宏帅|宏钊|定公|宜辰|宜香|宝乾|宝生|宝骅|宴宾|家兴|家虎|家辉|家锋|富娟|小东|小云|小伟|小军|小妹|小娟|小子|小宣|小平|小彪|小忠|小波|小洁|小涛|小琴|小瑞|小白|小红|小翠|小舜|小营|少妍|少宾|少怡|少栋|少炯|少辉|少锐|尚业|尧财|山峰|岩中|岩岩|师囊|希亮|常旭|平华|平波|平福|幸华|广友|广法|庆丰|庆琳|庆福|庭玉|延杰|廷良|建云|建仑|建伟|建功|建华|建峰|建忠|建超|建锋|建飞|建鹏|彤蕾|彦茹|彩妹|征桅|微波|德忠|德龙|志强|思远|惠平|惠强|慈仙|慈母|慧冬|慧娜|慧娟|慧才|慧琴|慧芳|成明|成春|所华|拴景|振利|振勇|振华|振潭|擎虹|改良"
	arr := make([]string, 0, 100)
	for _, char := range s {
		arr = append(arr, string(char))
	}
	return arr
})()

func RandName() string {
	rand.Seed(time.Now().UnixNano())
	surname := surName[rand.Intn(len(surName))]
	firstName := first[rand.Intn(len(first))]
	return surname + firstName
}
