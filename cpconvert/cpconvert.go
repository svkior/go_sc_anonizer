package cpconvert
import (
	"io/ioutil"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"strings"
)

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
