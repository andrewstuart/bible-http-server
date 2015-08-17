# osis
--
    import "github.com/andrewstuart/bible-http-server/osis"


## Usage

```go
var (
	NoSuchVerse = fmt.Errorf("No such verse")
	InvalidRef  = fmt.Errorf("Invalid verse reference")
)
```

#### type Bible

```go
type Bible struct {
	Books      []*Book        `json:"-"xml:"-"`
	BooksById  map[string]int `json:"-"xml:"-"`
	XMLName    xml.Name       `xml:"osis"json:"-"`
	Testaments []Testament    `xml:"osisText>div"json:"testaments"`
	Version    Version        `xml:"osisText>header>work"json:"version"`
}
```


#### func  NewBible

```go
func NewBible(r io.Reader) (*Bible, error)
```

#### func (*Bible) GetVerse

```go
func (b *Bible) GetVerse(ref string) (*Verse, error)
```

#### func (*Bible) GetVerses

```go
func (b *Bible) GetVerses(from, to string) ([]*Verse, error)
```
GetVerses takes two reference strings and returns a slice of verses

#### type Book

```go
type Book struct {
	Chs []*Chapter `xml:"chapter"`
	ID  string     `xml:"osisID,attr"`
}
```


#### type Chapter

```go
type Chapter struct {
	Vrs []*Verse `xml:"verse"`
	ID  string   `xml:"osisWork,attr"`
}
```


#### type Reference

```go
type Reference struct {
	Text  string `xml:",chardata"json:"text"json:"comment"`
	RefID string `xml:"osisRef,attr"json:"refid"`
}
```


#### type Testament

```go
type Testament struct {
	Books []Book `xml:"div"`
}
```


#### type Verse

```go
type Verse struct {
	Text  string      `xml:",chardata"json:"text"`
	ID    string      `xml:"osisID,attr"json:"id"`
	Words []Word      `xml:"w"json:"words,omitempty"`
	Refs  []Reference `xml:"note>reference"json:"references,omitempty"`
}
```


#### type VerseRef

```go
type VerseRef struct {
	BookID         string
	Chapter, Verse int
}
```


#### func  NewVerseRef

```go
func NewVerseRef(ref string) (*VerseRef, error)
```

#### type Version

```go
type Version struct {
	Title     string `xml:"title"json:"title"`
	ID        string `xml:"identifier"json:"id"`
	RefSystem string `xml:"refSystem"json:"refSystem"`
}
```


#### type Word

```go
type Word struct {
	N          int    `xml:"n,attr"`
	Text       string `xml:",chardata"`
	Translit   string `xml:"xlit,attr"`
	Lemma      string `xml:"lemma,attr"`
	Morphology string `xml:"morph,attr"`
}
```
