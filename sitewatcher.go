package main
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
)


type MPage struct {
	DocTitle string
	ParentId string
	LangId string
	DocIdent string
	DocContent string
	Cat string
	DocId string
	Title string
	Descr string
	Kw string

	IsNew bool
}

func (m *MPage) AddTrueLogo(){
	// TODO: Сделать ссылку на www.ttsy.ru
	m.DocContent += `
<script>
qq = document.querySelector("#topright");
qq.style.backgroundImage = "url(http://support.ttsy.ru/topright.gif)";
</script>
`
}

func (m *MPage) AddYellowPageRemover(){
	m.DocContent += `
<script>
	setTimeout(function(){
		foot = document.getElementById("footer");
		foot.parentNode.parentNode.removeChild(foot.parentNode);
	}, 30);
</script>
	`
}


func (m *MPage) AddGoogleAnalytics(){
	m.DocContent +=`
<script>
  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
  })(window,document,'script','//www.google-analytics.com/analytics.js','ga');

  ga('create', 'UA-66058123-1', 'auto');
  ga('send', 'pageview');

</script>
`
}

func NewSiteWatcher() *SiteWatcher{
	sw := SiteWatcher{
		Pages: make(map [string]*MPage),
		Ai: &adminIface{},
	}

	sw.Ai.getSession()
	_, body := sw.Ai.getLogin()
	sw.ProcessPages(body)
	//log.Println(string(body))
	return &sw
}

type SiteWatcher struct {
	Pages  map[string]*MPage
	Ai *adminIface
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
						strValue = ConvertCP2U(a.Val)
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
				str := ConvertCP2U(w.String())
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
	if err == nil {
		err = ioutil.WriteFile("./pages/" + page + ".json", b, 0644)
	}
}


func (sw *SiteWatcher) StatFile(page string){
	filename := "./pages/" + page + ".json"
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
	}

}

func (sw *SiteWatcher) UpdatePage(page *MPage){
	sw.Pages[page.DocIdent]= page
	sw.StatFile(page.DocIdent)
}


func (s *SiteWatcher) AddPage(page string){

	_, ok := s.Pages[page]
	if !ok {
		log.Printf("Adding new page %s", page)
		s.StatFile(page)
	} else {
		//log.Printf("Page %s is already exists", page)
	}
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
							//log.Printf("LINK: %s", a.Val)
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


