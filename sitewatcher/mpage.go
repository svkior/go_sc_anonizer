package sitewatcher
import (
	"net/url"
	"bitbucket.org/svkior/go_sc_anonizer/cpconvert"
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

func (m *MPage) GetUrlValues() url.Values{


	return url.Values{
		"doc_title": {cpconvert.ConvertU2CP(m.DocTitle)},
		"parent_id": {cpconvert.ConvertU2CP(m.ParentId)},
		"lang_id": {cpconvert.ConvertU2CP(m.LangId)},
		"doc_ident": {cpconvert.ConvertU2CP(m.DocIdent)},
		"DocumentContent": {cpconvert.ConvertU2CP(m.DocContent)},
		"cat": {cpconvert.ConvertU2CP(m.Cat)},
		"doc_id": {cpconvert.ConvertU2CP(m.DocId)},
		"title": {cpconvert.ConvertU2CP(m.Title)},
		"descr": {cpconvert.ConvertU2CP(m.Descr)},
		"kw": {cpconvert.ConvertU2CP(m.Kw)},
	}
}
