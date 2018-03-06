package object

import (
	"os"

	"github.com/fuxiaohei/pugov1/module/i18n"
)

// Source is all source data object
type Source struct {
	PostDir string
	PageDir string

	PostFiles []SourceItem
	PageFiles []SourceItem

	CopyFiles     []SourceItem
	RenderedFiles []string

	Posts     []*Post
	Lists     []PostList
	ListSize  int
	TagsLists map[string]*TagList
	Archives  []*Archive

	Pages []*Page
	I18ns map[string]*i18n.File
}

const (
	// SourceOpPost means source item is a post
	SourceOpPost = 1
	// SourceOpPage means source item is a page
	SourceOpPage = 3
	// SourceOpFile means source item is a common file
	SourceOpFile = 5
)

// SourceItem is an item reading from source directry
type SourceItem struct {
	SrcFile string
	File    string
	Info    os.FileInfo
	OpType  int
}
