package object

// Archive is archive set for posts
type Archive struct {
	Year  int // each list by year
	Posts []*Post
}
