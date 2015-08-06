package main
import (
	"net/http"
	"log"
	"strings"
	"golang.org/x/net/html"
	"net/url"
	"os"
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
		Ai: adminIface{},
	}

	sw.Ai.getSession()
	sw.Ai.getLogin()
	return *sw
}

type SiteWatcher struct {
	Pages  map[string]MPage
	Ai *adminIface
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


