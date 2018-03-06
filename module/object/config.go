package object

type (
	// Config is all config from config file
	Config struct {
		Meta      *Meta       `toml:"meta" ini:"meta"`
		Navs      NavGroup    `toml:"nav" ini:"nav"`
		Authors   AuthorGroup `toml:"author" ini:"author"`
		Comment   *Comment    `toml:"comment" ini:"comment"`
		Analytics *Analytics  `toml:"analytics" ini:"analytics"`
		SrcFile   string      `toml:"-" ini:"-"`
	}
)

// Meta is metadata info for site
type Meta struct {
	Title    string `toml:"title" ini:"title"`
	Subtitle string `toml:"subtitle" ini:"subtitle"`
	Keyword  string `toml:"keyword" ini:"keyword"`
	Desc     string `toml:"desc" ini:"desc"`
	Domain   string `toml:"domain" ini:"domain"`
	Root     string `toml:"root" ini:"root"`
	Lang     string `toml:"lang" ini:"lang"`
}

// Nav is one item in navigator
type Nav struct {
	Link    string `toml:"link" ini:"link"`
	Title   string `toml:"title" ini:"title"`
	I18n    string `toml:"i18n" ini:"i18n"`
	Hover   string `toml:"hover" ini:"hover"`
	IsBlank bool   `toml:"is_blank" ini:"is_blank"`
	Icon    string `toml:"icon" ini:"icon"`
}

// NavGroup is navigator for site
type NavGroup []*Nav

// Author is author info
type Author struct {
	Name  string `toml:"name" ini:"name"`
	Nick  string `toml:"nick" ini:"nick"`
	Email string `toml:"email" ini:"email"`
	URL   string `toml:"url" ini:"url"`
}

// AuthorGroup are authors for site
type AuthorGroup []*Author

// Comment is comment settings for site
type Comment struct {
	Disqus string `toml:"disqus" ini:"disqus"`
}

// Analytics is third-party analytics settings for site
type Analytics struct {
	Baidu   string `toml:"baidu" ini:"baidu"`
	Google  string `toml:"google" ini:"google"`
	Tencent string `toml:"tencent" ini:"tencent"`
	Cnzz    string `toml:"cnzz" ini:"cnzz"`
}
