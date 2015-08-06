package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"net/url"
	"golang.org/x/text/transform"
	"golang.org/x/text/encoding/charmap"
)

const PAGE_ID = "contacts1"






func ConvertU2CP(in string) string{
	sr := strings.NewReader(in)
	tr := transform.NewReader(sr, charmap.Windows1251.NewEncoder())
	buf, _ := ioutil.ReadAll(tr)
	return string(buf)
}

func ConvertCP2U(in string) string{
	sr := strings.NewReader(in)
	tr := transform.NewReader(sr, charmap.Windows1251.NewDecoder())
	buf, _ := ioutil.ReadAll(tr)
	return string(buf)
}


func (m *MPage) GetUrlValues() url.Values{


	return url.Values{
		"doc_title": {ConvertU2CP(m.DocTitle)},
		"parent_id": {ConvertU2CP(m.ParentId)},
		"lang_id": {ConvertU2CP(m.LangId)},
		"doc_ident": {ConvertU2CP(m.DocIdent)},
		"DocumentContent": {ConvertU2CP(m.DocContent)},
		"cat": {ConvertU2CP(m.Cat)},
		"doc_id": {ConvertU2CP(m.DocId)},
		"title": {ConvertU2CP(m.Title)},
		"descr": {ConvertU2CP(m.Descr)},
		"kw": {ConvertU2CP(m.Kw)},
	}
}








const PAGE_BODY = `
<div>
<h1 style="text-align: center;"><span style="font-size: medium;">ООО &laquo;Стройцирк&raquo;</span></h1>
<h2 style="text-align: center;"><span style="font-size: medium;">+7 (495) 969-63-56</span></h2>
<h2 style="text-align: center;"><span style="font-size: medium;">+7 (495) 542-17-72</span></h2>
<h2 style="text-align: center;"><span style="font-size: medium;">E-mail: scircus2@list.ru</span></h2>
<p>&nbsp;</p>

<div id="contactsss">&nbsp;</div>
<script>
var inc = 0;
var timer = setInterval(function(){
if(inc < 1){
	document.getElementById("contactsss").innerHTML =
'<h1 style="text-align: center;"><span style="font-size: medium;">ООО &laquo;Театральные Технологические Системы&raquo;</span></h1>'
+'<h2 style="text-align: center;"><span style="font-size: medium;">+7 (495) 730-83-45</span></h2>'
+'<h2 style="text-align: center;"><span style="font-size: medium;">+7 (495) 730-83-46</span></h2>'
+'<h2 style="text-align: center;"><span style="font-size: medium;">+7 (499) 649-29-40</span></h2>'
+'<p>&nbsp;</p>'
+'<h2 style="text-align: center;"><span style="font-size: medium;"><a href="mailto:info@ttsy.ru">info@ttsy.ru</a></span></h2>'
+'<h2 style="text-align: center;"><span style="font-size: medium;"><a href="http://www.ttsy.ru/">www.ttsy.ru</a></span></h2>';
	inc = inc + 1
} else {
   clearInterval(timer);
}
}, 100)


</script></div>
`




func mainOld() {
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan bool)
	go func(){
		for {
			select {
			case <- ticker.C:
				res, err := http.Get("http://scircus.ru/index.php?load=contacts")
				if err != nil {
					log.Fatal(err)
				}
				robots, err := ioutil.ReadAll(res.Body)
				res.Body.Close()
				if err != nil {
					log.Fatal(err)
				}
				//fmt.Printf("%s", robots)
				if strings.Contains(string(robots), "+7 (495) 730-83-45") {
					log.Println("OK")
				} else {
					log.Println("AHTUNG")
				}
			}
		}
	}()
	<- quit
}

func main(){

	sw := NewSiteWatcher()
	sw.readMain()

	//as := adminIface{}
	//go as.DoTheJob()
	//mainOld()
}