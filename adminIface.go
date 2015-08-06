package main
import (
	"net/http"
	"strconv"
	"bytes"
	"io/ioutil"
	"log"
	"net/http/cookiejar"
	"net/url"
	"time"
	"golang.org/x/net/html"
)

type adminIface struct {
	cookies []*http.Cookie
	loggedIn bool
	client *http.Client
}

func (ai *adminIface) postNewPage(){

	page := MPage{
		DocTitle: "Сообщение",
		ParentId: "0",
		LangId: "ru",
		DocIdent: "contacts",
		DocContent: "<h1>Привет Мир</h1>",
		Cat: "",
		DocId: "",
		Title: "",
		Descr: "",
		Kw: "",
	}

	page.AddTrueLogo()
	page.AddYellowPageRemover()

	data := page.GetUrlValues()


	req, _ := http.NewRequest("POST", "http://www.scircus.ru/admin/editor/index.php?save", bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	resp, err := ai.client.Do(req)
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	log.Println(string(body))

}

func (ai *adminIface) createNewPage(){
	log.Println("Need to create New Page")
	req, _ := http.NewRequest("GET", "http://www.scircus.ru/admin/?a=adddoc", nil)
	resp, err := ai.client.Do(req)
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	log.Println(string(body))

	r := bytes.NewReader(body)
	doc, err := html.Parse(r)
	if err != nil {
		log.Fatal(err)
	}
	var f func(string, *html.Node)
	f = func(pre string, n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "input" || n.Data == "option" {
				log.Println(pre,n.Data, n.Attr)
			} else {

			}
			//log.Println(pre,n.Data)

		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(pre+">", c)
		}
	}
	f(">", doc)
}

func (ai *adminIface) updateOldPage(url string){
	log.Printf("Need to update Old Page:%s", url)
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

func (ai *adminIface) getLogin() (bool, []byte) {
	jar, _ := cookiejar.New(nil)
	cookieURL, _ := url.Parse("/") // http://www.scircus.ru/admin/index.php
	jar.SetCookies(cookieURL, ai.cookies)

	ai.client = &http.Client{
		Jar: jar,
	}

	data := url.Values{"login": {"srZkg5eg"}, "pass": {"YcTbIJeg"}}

	r, _ := http.NewRequest("POST", "http://www.scircus.ru/admin/index.php", bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := ai.client.Do(r)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200{
		ai.loggedIn = true

		//log.Printf("COOKIES: %#v",resp.Cookies())
		//ai.processBody(body)
		return true, body
	} else {
		log.Println(resp)
		log.Println(body)
		return false, nil
	}
}


func (ai *adminIface) getEditorContents(getUrl string) []byte{
	time.Sleep(500 * time.Millisecond)

	req, _ := http.NewRequest("GET", getUrl, nil)
	resp, err := ai.client.Do(req)
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return body
}

func (ai *adminIface) DoTheJob(){
	time.Sleep(1 * time.Second)
	ai.getSession()
	time.Sleep(1 * time.Second)
	ai.getLogin()
	time.Sleep(1 * time.Second)
	ai.postNewPage()
	/*
	if !ai.loggedIn{
		log.Println("Not Logged In")
		ai.getLogin()
	} else {
		log.Println("Logged In")

	}
	*/
}
