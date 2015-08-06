package main
import (
	"net/http"
	"log"
	"strings"
	"golang.org/x/net/html"
	"net/url"
	"os"
	"bytes"
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

func NewSiteWatcher() *SiteWatcher{
	sw := SiteWatcher{
		Pages: make(map [string]MPage),
		Ai: &adminIface{},
	}

	sw.Ai.getSession()
	_, body := sw.Ai.getLogin()
	sw.ProcessPages(body)
	//log.Println(string(body))
	return &sw
}

type SiteWatcher struct {
	Pages  map[string]MPage
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

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			//log.Println(n.Data)
			if n.Data == "input" {

				for _, a := range n.Attr {
					if a.Key == "name" {
						//log.Println(n.Data, a.Key, a.Val)

						switch a.Val{
						case "doc_title":
							log.Println("Doc Title", a.Key)
						}

						break
					}
				}
			} else {

			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}

	}

	f(doc)
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


func (s *SiteWatcher) AddPage(page string){

	_, ok := s.Pages[page]
	if !ok {
		log.Printf("Adding new page %s", page)
		filename := "./pages/" + page
		var notExists bool
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			log.Printf("Page: %s is new, adding to list of pages", page)
			notExists = true
		}
		s.Pages[page] = MPage{
			DocIdent: "page",
			IsNew: notExists,
		}

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


