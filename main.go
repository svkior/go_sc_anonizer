package main

import (
	"bitbucket.org/svkior/go_sc_anonizer/sitewatcher"
)

const PAGE_ID = "contacts1"

/*

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

*/

func main(){
	sitewatcher.WebInterfaceRun()

	//as := adminIface{}
	//go as.DoTheJob()
	//mainOld()
}