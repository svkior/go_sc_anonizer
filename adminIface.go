package main
import (
	"net/http"
	"strconv"
	"bytes"
	"io/ioutil"
	"log"
	"net/http/cookiejar"
	"net/url"
	"fmt"
	"time"
	"golang.org/x/net/html"
	"strings"
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

func (ai *adminIface) processBody(b []byte) bool{
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
					if n.FirstChild.Data == PAGE_ID{
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
		ai.updateOldPage(editUrl)
	} else {
		log.Println("Page is not found")
		ai.createNewPage()
	}
	return found
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
		ai.processBody(body)
		return true
	} else {
		log.Println(resp)
		log.Println(body)
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
