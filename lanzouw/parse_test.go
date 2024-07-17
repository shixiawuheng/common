package lanzouw

//func TestRequest(t *testing.T) {
//
//	uri := "https://qf1.lanzouw.com/id8lE0v4c1ti"
//	pwd := "59i4"
//	var f File
//	err := Parse(uri, pwd, &f)
//	if err != nil {
//		return
//	}
//	err = Download(f, filepath.Join("D:\\", "test.zip"))
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	//os.Open()
//
//}
//
//func TestCap(t *testing.T) {
//	html := `<!DOCTYPE HTML>
//<html>
//<head>
//<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
//<meta name="viewport" content="width=device-width,initial-scale=1.0,maximum-scale=1.0,user-scalable=0" />
//<script type="text/javascript" src="https://apps.bdimg.com/libs/jquery/1.6.4/jquery.min.js"></script>
//</head>
//<body>
//<style>
//body{font-size: 16px;margin: 0;font-family: \5FAE\8F6F\96C5\9ED1;color: #333;background: #fff; border-color:transparent;-webkit-appearance: none;-webkit-tap-highlight-color:rgba(0,0,0,0);-webkit-tap-highlight-color:rgba(0,0,0,0.0);}
//a{text-decoration: inherit;color: #999;}
//.box{width:260px;margin:auto;    text-align: center;}
//.box1{    margin-top: 35%;}
//.box3{    margin-top: 50px;}
//.box2_2{    padding: 25px;}
//.submit{    width: 90%;
//    margin: auto;
//    height: 42px;
//    line-height: 42px;
//    border: 1px solid #d5d5d5;
//    cursor: pointer;
//        box-shadow: 0 1px 2px rgba(0, 0, 0, 0.075);border-radius: 2px;    color: #777;}
//.submit:hover{    opacity: .7;}
//.box2_1{    width: 60px;
//    height: 60px;
//    font-size: 30px;
//    line-height: 60px;
//    border-radius: 50%;
//    background: #09f;
//    color: #fff;
//    margin: auto;}
//#ok{background: #d4ffa3;
//    width: 60px;
//    height: 50px;
//    border-radius: 50%;
//    margin: auto;
//    padding-top: 10px;    display: none;}
//.ok{    height: 20px;display: none;
//    width: 30px;
//    border-left: 2px solid #84c33c;
//    border-bottom: 2px solid #84c33c;
//    transform: rotate(-45deg);
//    margin-left: 15px;   }
//#ok2{display: none;}
//.goto{    border: 1px solid #fff;    color: #fff;box-shadow: inherit;
//    background-color: #fff;
//    background: #fff;}
//#top{    font-size: 22px;padding: 20px;}
//@media (max-width: 768px){
//	.box{width:100%}
//	#top{font-size: 0;padding: 20px;}
//	}
//
//.load2 .loader,
//.load2 .loader:before,
//.load2 .loader:after {border-radius: 50%;}
//.load2 .loader:before,
//.load2 .loader:after {position: absolute;content: '';}
//.load2 .loader:before {width: 16px;height: 32px;background: #fff;border-radius: 32px 0 0 32px;top: -1px;left: -1px;transform-origin: 16px 16px;animation: load2 2s infinite ease 1.5s;}
//.load2 .loader {font-size: 11px;text-indent: -99999em;margin: auto;position: relative;width: 30px;height: 30px;box-shadow: inset 0 0 0 4px #09f;border-radius: 50%;}
//.load2 .loader:after {width: 16px;height: 32px;background: #fff;border-radius: 0 31px 31px 0;top: -1px;left: 15px;transform-origin: 0px 16px;animation: load2 2s infinite ease;}
//@keyframes load2 {
//  0% {
//    transform: rotate(0deg);
//  }
//  100% {
//    transform: rotate(360deg);
//  }
//}
//.load2{    position: fixed;
//    width: 100%;
//    background: #fff;
//    height: 100%;
//    z-index: 999;
//    top: 0;
//    text-align: center;    padding-top: 100px;}
//.loader2{    padding-top: 30px;}
//</style>
//<div id="top">下载文件</div>
//<div class="load2" id="load2"><div class="loader"></div><div class="loader2">努力加载中...</div></div>
//
//<div class="box">
//<div id="box1" class="box1">
//<div id="ok"><div class="ok"></div></div>
//<div id="ok2" class="box2_2">已经通过验证</div>
//<div id="go" class="submit goto"><div onclick="down_r(1);">验证并下载</div></div>
//</div>
//
//<div id="box2" class="box2">
//<div class="box2_1">!</div>
//<div class="box2_2">系统发现您的网络异常，需要验证后下载文件</div>
//<div id="sub" class="submit"><div onclick="down_r(2);">验证并下载</div></div>
//</div>
//
//<div id="box3" class="box3">
//<div id="sub2" class="submit goto"><div onclick="down_r(3);">验证并下载</div></div>
//</div>
//
//</div>
//<script type="text/javascript">
//	setTimeout('document.getElementById("load2").style.display="none";',2000);
//	function down_r(el){
//		$('#sub').text("提交中...");
//		$.ajax({
//			type : 'post',
//			url : 'ajax.php',
//			data : { 'file':'AGhbeFpzDysBJFA8A3hdcwQ3AD0HaQAsC2BbPAEzVS4GaVNxADpSKFFkBTxUbVN9VT8OPgRqADJQZVRiVDhbagAxWz1aMw9vAWFQMgNmXW0EZgBnBzoALQs+W20Bb1VnBilTNABoUilRNgVjVC9TZlVrDmkEPwBjUG9UYFRrW2kAMVs/WjQPPQFkUDADNF1pBG8AMAc9AGcLPFttAW9VNwYwU2UAZFJgUTIFMFRiU3xVdQ5iBCsAPVAkVCZUNVsTAG1bNFpvDwIBE1BsAzFdBQQTADAHHgBgC2FbCQEIVR0Gc1NcAD9SPlFwBXVUZVNvVT4OPQRjADFQYFRgVDxbaAAzWz1aIQ85AWpQRwMAXRkEEABVBwkASwt7Wx8BGVUeBlRTSQALUldRfQUQVDJTIVVHDkYELABUUDhUBVRiWwgAWFtWWkYPEgFmUEcDAF0JBBIATgcZAFgLO1sfARpVbQZRU2AAHFJZUWQFDFRjU3RVaQ5iBGYAM1BmVGZUPFtsADRbPVo2D2MBcVB2Az5dOARjADcHagAxCyFbZQFvVXkGN1M1AGhSK1E0BWtUJlMnVX8ONgRpACRQOlQiVDVbagAmW29aaA9mAWY=','el':el,'sign':'UmNRMA5rDmkCYwJmUDZcbwNqAjYFOw==' },
//			dataType : 'json',
//			success:function(msg){
//				if(msg.zt == '1'){
//					$("#box2").hide(150);
//					$("#ok").show(150);
//					$("#ok2").show(150);
//					$(".ok").show(500);
//					$("#go").removeClass("goto");
//					$("#go").html("<a href="+ msg.url +"><span class=tc2>立即下载</span></a>");
//				}else{
//					//$('#sub').val("验证并下载");
//					setTimeout('window.location.reload();',1000);
//				}
//			},
//			error:function(){
//				//$("pwderr").text("获取失败，请重试");
//			}
//
//		});
//	}
//</script>
//</body>
//<div style="display:none"><script src="https://s6.cnzz.com/z_stat.php?id=5289169&web_id=5289169" language="JavaScript"></script></div>
//</html>`
//	param := common.GetBetweenStr(html, "data :", "},") + "}"
//	param = strings.Replace(param, "'", `"`, -1)
//	param = strings.Replace(param, "el,", `2,`, -1)
//	var data map[string]interface{}
//	json.Unmarshal([]byte(param), &data)
//	js, _ := json.Marshal(data)
//	fmt.Println(data)
//	client := &http.Client{
//		CheckRedirect: func(req *http.Request, via []*http.Request) error {
//			return http.ErrUseLastResponse
//		},
//	}
//	req, err := http.NewRequest("POST", "", bytes.NewReader(js))
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	req.Header = Header
//	req.Header.Set("Content-Type", "application/json")
//	resp, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer resp.Body.Close()
//	body, _ := ioutil.ReadAll(resp.Body)
//	fmt.Println(string(body))
//}
