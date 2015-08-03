package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/parnurzeal/gorequest"
	"net/url"
	"fmt"
	"bytes"
	"net/http/cookiejar"
	"golang.org/x/net/html"
	"strconv"
)




/*
Форма создания файла

<form action="/admin/editor/index.php?save" method="post">
<table width="100%" border="0" cellpadding="0" cellspacing="0">
<tbody><tr>
<td width="100%">
<input style="background-color:#FFFFFF; color:#0000FF; width:100%;" type="text" id="doc_title" name="doc_title" value="">
</td>
<td><select name="parent_id"><option value="0">Место размещения</option><option value="68">Карта сайта</option><option value="48">О компании</option><option value="61">Прайс</option><option value="62">Продукция</option><option value="97">Связь</option></select></td>
<td><select name="lang_id" id="lang_id"><option value="ru" selected="">Русский</option></select></td>
<td><input type="text" name="doc_ident" value=""></td>
</tr>
</tbody></table>
<input type="hidden" id="DocumentContent" name="DocumentContent" value="" style="display:none"><input type="hidden" id="DocumentContent___Config" value="AutoDetectLanguage=false&amp;DefaultLanguage=ru" style="display:none"><iframe id="DocumentContent___Frame" src="editor/fckeditor.html?InstanceName=DocumentContent&amp;Toolbar=Default" width="100%" height="500" frameborder="0" scrolling="no" style="margin: 0px; padding: 0px; border: 0px; width: 100%; height: 500px; background-image: none; background-color: transparent;"></iframe>
<input type="hidden" name="cat" value="">
<input type="hidden" name="doc_id" value="">
<table width="100%" border="0" cellpadding="4" cellspacing="2">
<tbody><tr><td colspan="2"><h3>Для раскрутки страницы</h3><hr></td></tr>
<tr>
<td>Заголовок страницы (title)</td>
<td width="90%"><input style="color:#0000FF;width:100%;" type="text" name="title" value=""></td>
</tr>
<tr>
<td>Описание (description)</td>
<td width="90%"><input style="color:#0000FF;width:100%;" type="text" name="descr" value=""></td>
</tr>
<tr>
<td>Ключевые слова (keywords)</td>
<td width="90%"><input style="color:#0000FF;width:100%;" type="text" name="kw" value=""></td>
</tr>
</tbody></table></form>


 */


/*
Форма логина
<form action="" method="POST">
<input class="login_text_field" style="width:100%;" name="login" onfocus="onLoginFocus();" onblur="onLoginBlur();" type="text" value="" id="l_field">
<input class="login_text_field" id="p_field" name="pass" type="password" onfocus="onPasswordFocus();" onblur="onPasswordBlur();" value="">
<input class="button" type="submit" value="Авторизоваться">
</form>
 */

var bod string

type adminIface struct {
	cookies []*http.Cookie
	loggedIn bool
	client *http.Client
}

func (ai *adminIface) processBody(b []byte){
//	log.Println(string(b))

	r := bytes.NewReader(b)
	doc, err := html.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	var state_m int
	var editUrl string
	var found bool

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode{
			if found {
				return
			}
			switch state_m{
			case 0: // Поиск a
				if  n.Data == "a" {
					for _, a := range n.Attr {
						if a.Key == "href" {
							if strings.Contains(a.Val, "editor"){
								editUrl = a.Val
								//fmt.Println("OPENING ",a.Val)
								state_m  = 1
							}
							break
						}
					}
				}
			case 1:
				if n.Data == "b"{
					//fmt.Println("B1 FOUND")
					state_m = 2
				}
			case 2:
				if n.Data == "b"{
					//fmt.Println("B2 FOUND")
					//fmt.Printf("Node: %#v\n", n)
					fmt.Printf("Child: %#v\n", n.FirstChild.Data)
					if n.FirstChild.Data == "contacts"{
						found = true
						return
					}
					state_m = 3
				}
			case 3:
				if  n.Data == "a" {
					for _, a := range n.Attr {
						if a.Key == "href" {
							if strings.Contains(a.Val, "editor"){
								//fmt.Println("CLOSING ",a.Val)
								state_m  = 0
							}
							break
						}
					}
				}
			default:
				state_m = 0
			}

			//log.Println(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)


	if found {
		log.Printf("Found page: %s", editUrl)
	}
}

func (ai *adminIface) getSession() bool{
	resp, err := http.Get("http://www.scircus.ru/admin/index.php")
	if err != nil {
		log.Println(err.Error())
		return false
	}
	resp.Body.Close()
	log.Printf("Cookies: %v",resp.Cookies())
	ai.cookies = resp.Cookies()
	return true
}

func (ai *adminIface) getLogin() bool{
	jar, _ := cookiejar.New(nil)
	cookieURL, _ := url.Parse("/") // http://www.scircus.ru/admin/index.php
	jar.SetCookies(cookieURL, ai.cookies)

	client := &http.Client{
		Jar: jar,
	}

	data := url.Values{"login": {"srZkg5eg"}, "pass": {"YcTbIJeg"}}

	r, _ := http.NewRequest("POST", "http://www.scircus.ru/admin/index.php", bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200{
		ai.cookies = resp.Cookies()
		log.Println(ai.cookies)
		ai.loggedIn = true
		ai.processBody(body)
		return true
	} else {
		log.Println(resp)
		return false
	}
}

func (ai *adminIface) getListOfPages() {
	jar, _ := cookiejar.New(nil)
	cookieURL, _ := url.Parse("/") // http://www.scircus.ru/admin/index.php

	jar.SetCookies(cookieURL, ai.cookies)

	// sanity check
	fmt.Println(jar.Cookies(cookieURL))
	client := &http.Client{
		Jar: jar,
	}
	req, _ := http.NewRequest("GET", "http://www.scircus.ru/admin/index.php", nil)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	ai.processBody(body)
}

func (ai *adminIface) DoTheJob(){
	time.Sleep(1 * time.Second)
	ai.getSession()
	time.Sleep(1 * time.Second)
	ai.getLogin()
	time.Sleep(1 * time.Second)
	/*
	if !ai.loggedIn{
		log.Println("Not Logged In")
		ai.getLogin()
	} else {
		log.Println("Logged In")

	}
	*/
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


/*
document.cookie = "genipanel_session=7909db8d841240fe1c32d82aab25dd7c"
<form action="/admin/editor/index.php?save" method="post">

<input type="text" id="doc_title" name="doc_title" value="Связь">
<select name="parent_id">
   <option value="0">Место размещения</option>
   <option value="68">Карта сайта</option>
   <option value="93">Наши контакты</option>
   <option value="48">О компании</option>
   <option value="61">Прайс</option>
   <option value="62">Продукция</option>
   <option value="92">Связь</option>
</select>

<select name="lang_id" id="lang_id">
   <option value="ru" selected="">Русский</option>
</select>

<input type="text" name="doc_ident" value="contacts">
<input type="hidden" id="DocumentContent" name="DocumentContent" value="СЮДА ТЕКСТ" style="display:none">

 */



func doHast(){
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "genipanel_session", Value: "7909db8d841240fe1c32d82aab25dd7c", Expires: expiration}

	sa := gorequest.New()
	sa.AddCookie(&cookie).
	Post("/admin/editor/index.php?save").
	Query(`{doc_title: 'Связь'`).
	Query(`{parent_id: '0'}`).
	Query(`{lang_id: 'ru'}`).
	Query(`{doc_ident: 'contacts'}`).
	Query(`{DocumentContent: `+ PAGE_BODY + `}`).
	Query(`{cat: ''`)


/*

<input type="hidden" id="DocumentContent___Config" value="AutoDetectLanguage=false&amp;DefaultLanguage=ru" style="display:none">
<input type="hidden" name="cat" value="">
<input type="hidden" name="doc_id" value="92">

<input style="color:#0000FF;width:100%;" type="text" name="descr" value="">
<input style="color:#0000FF;width:100%;" type="text" name="kw" value="">
</form>
*/
}



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
	as := adminIface{}
	go as.DoTheJob()
	mainOld()
}