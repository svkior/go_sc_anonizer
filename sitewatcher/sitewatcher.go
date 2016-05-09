package sitewatcher

import (
	"net/http"
	"log"
	"strings"
	"golang.org/x/net/html"
	"net/url"
	"os"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"bitbucket.org/tts/go_webtest/artnet_test/element"
	"bitbucket.org/tts/go_webtest/artnet_test/trace"
	"time"
	"bitbucket.org/tts/go_webtest/artnet_test/filewatcher"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/StephanDollberg/go-json-rest-middleware-jwt"
	"sync"
	"html/template"
	"path/filepath"
	"bitbucket.org/tts/go_webtest/artnet_test/fileconfigurator"
	"bitbucket.org/tts/go_webtest/artnet_test/auth"
	"bitbucket.org/svkior/go_sc_anonizer/cpconvert"
	"sort"
)


type SiteWatcher struct {
	element.AbstractElement
	Pages  map[string]*MPage
	Ai *adminIface
}


func NewSiteWatcher() *SiteWatcher{
	sw := SiteWatcher{
		AbstractElement: *element.NewAbstractElement("pages"),
		Pages: make(map [string]*MPage),
		Ai: &adminIface{},
	}
	sw.Handle("get_all", sw.HandleGetPages)
	sw.AbstractElement.OnSubscribe = sw.OnSubscribeFunc
	sw.Ai.getSession()
	_, body := sw.Ai.getLogin()
	sw.ProcessPages(body)
	log.Println("Running timer:")
	go sw.timer()
	return &sw
}


func (sw *SiteWatcher) HandleGetPages(msg *element.Message) (bool, error){
	log.Println("Hello, World!")
	return true, nil
}

func (sw *SiteWatcher) OnSubscribeFunc(client element.Element){
	//TODO: Нужно клиенту отдать список страничек
	var keys []string

	for k := range sw.Pages{
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _,k2 := range keys {
		msgOut := element.GetEmptyMessage("123", false)
		msgOut.Name = "pages"
		msgOut.Type = "page"
		msgOut.Payload = sw.Pages[k2]
		client.GetRecv() <- msgOut
	}
}


func (sw *SiteWatcher) DownloadPage(nam string, url string){
	fullURL := "http://www.scircus.ru/admin/" + url
	log.Printf("Downloading page: %s from url: %s", nam, fullURL)

	body := sw.Ai.getEditorContents(fullURL)

	//log.Println(string(body))
	r := bytes.NewReader(body)
	doc, err := html.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	m := MPage{}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "input" {
				//log.Println(n.Data, n.Attr)
				hasName := false
				hasValue := false
				strName := ""
				strValue := ""
				for _, a := range n.Attr {
					switch a.Key{
					case "name":
						strName = a.Val
						hasName = true
					case "value":
						strValue = cpconvert.ConvertCP2U(a.Val)
						hasValue = true
					}
					if hasName && hasValue {
						break
					}
				}

				if hasName && hasValue {
					//log.Printf("We got param %s : %s", strName, strValue)
					switch strName{
					case "kw":
						m.Kw = strValue
					case "descr":
						m.Descr = strValue
					case "title":
						m.Title = strValue
					case "doc_id":
						m.DocId = strValue
					case "cat":
						m.Cat = strValue
					case "doc_ident":
						m.DocIdent = strValue
					case "doc_title":
						m.DocTitle = strValue
					}
				}
			} else if n.Data == "textarea"{
				w := bytes.NewBuffer(nil)
				html.Render(w, n.FirstChild)
				str := cpconvert.ConvertCP2U(w.String())
				//log.Println(str)
				m.DocContent = str
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}

	}

	f(doc)

	sw.UpdatePage(&m)
}

func (s *SiteWatcher) ProcessPages(b []byte){
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

					//log.Printf("%s : %s ", editUrl, n.FirstChild.Data)
					s.DownloadPage(n.FirstChild.Data, editUrl)
					//fmt.Printf("Child: %#v\n", n.FirstChild.Data)
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

}

func (sw *SiteWatcher) WriteToFile(page string){
	b, err := json.Marshal(sw.Pages[page])
	if err != nil{
		log.Println(err)
		return
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	ioutil.WriteFile("./pages/" + page + ".json", out.Bytes(), 0644)
}





func (sw *SiteWatcher) StatFile(page string){


	filename := "./pages/" + page + ".json"

	log.Printf("Try: %s", filename)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Printf("Page: %s is new, adding to list of pages", page)
		mp, ok := sw.Pages[page]
		if ok {
			mp.IsNew = true
		} else {
			log.Println("New Not Existing page!!!!")
			mp = &MPage{
				DocIdent: page,
				IsNew: true,
			}
			mp.AddTrueLogo()
			mp.AddYellowPageRemover()
			mp.AddGoogleAnalytics()
			sw.Pages[page] = mp
			// TODO: Здесь делаем запись в файл
		}
		sw.WriteToFile(page)
	} else {
		log.Printf("File exists. We do not need to create one")
	}

}




func (sw *SiteWatcher) CheckForUpdates(){
	for _, mp := range sw.Pages{
		log.Printf("%v : %v", mp.DocId, mp.DocIdent)
		
	}
}

func (sw *SiteWatcher) UpdatePage(page *MPage){
	sw.Pages[page.DocIdent]= page
	sw.StatFile(page.DocIdent)
}


func (s *SiteWatcher) AddPage(page string){
	//log.Printf("Looking for page %s", page)
	_, ok := s.Pages[page]
	if !ok {
		log.Printf("Adding new page %s", page)
		s.StatFile(page)
	} else {
		//log.Printf("Page %s is already exists", page)
		// TODO: Зделать разбор полетов
	}
}

func (s *SiteWatcher) timer(){
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan bool)
	go func(){
		for {
			select {
			case <- ticker.C:
				log.Println("Tick")


				// TODO: Разбор нужных
				s.CheckForUpdates()
				/*

				res, err := http.Get("http://scircus.ru/")
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

				*/
			}
		}
	}()
	<- quit
}

func (s *SiteWatcher) readMain() {
	res, err := http.Get("http://scircus.ru/")
	if err != nil {
		log.Fatal(err)
	}

	doc, err := html.Parse(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {

			if n.Data == "a"{
				//log.Printf("AA : %#v", n)
				for _, a := range n.Attr {
					if a.Key == "href" {

						if strings.Contains(a.Val, "load"){
							log.Printf("LINK: %s", a.Val)
							qv, _ := url.ParseQuery(a.Val)
							link := qv["index.php?load"]
							if strings.Contains(link[0], "?") == false {
								//log.Printf("Parse: %#v", link[0])
								s.AddPage(link[0])
							}
						}
						break
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}

type templateHandler struct {
	once sync.Once
	filename string
	templ *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	t.once.Do(func(){
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}


func NewRestInterface(device *element.Device){
	jwt_middleware := &jwt.JWTMiddleware{
		Key:        []byte("secret key"),
		Realm:      "jwt auth",
		Timeout:    time.Hour,
		MaxRefresh: time.Hour * 24,
		Authenticator: func(userId string, password string) bool {
			return userId == "admin" && password == "admin"
		}}
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			return request.URL.Path != "/login"
		},
		IfTrue: jwt_middleware,
	})

	router, err := rest.MakeRouter(
		rest.Post("/login", jwt_middleware.LoginHandler),
		rest.Get("/refresh_token", jwt_middleware.RefreshHandler),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	http.Handle("/", &templateHandler{filename:"main.html"})
	http.Handle("/device", device)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func WebInterfaceRun(){


	log.Println("*** Starting Device ***")
	// Создаем новое устройство
	device := element.NewDevice()
	device.Tracer = trace.New(os.Stdout)

	// Запускаем на выполнение
	device.Run()
	// Добавляем конфигуратор
	device.AddElement(
		fileconfigurator.NewFileConfig("test.json"),
	)
	// Запускаем конфигуратор на выполнение
	device.SendMessage(element.GetEmptyMessage("load", true))

	// Запускаем FileWatcher
	fw := filewatcher.NewFileWatcher("./assets/main.js")
	device.AddElement(fw)

	device.AddElement(auth.NewAuth("hello, world"))

	sw := NewSiteWatcher()
	device.AddElement(sw)


	// Уходим в главный цикл программы

	go NewRestInterface(device)

	device.Wait()
	log.Println("Quit()")

}

