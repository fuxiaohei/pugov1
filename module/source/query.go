package source

import "github.com/fuxiaohei/pugov1/module/object"

// QueryPost query post with slug and created time
func QueryPost(posts []*object.Post, slug string, ct int64) *object.Post {
	for _, post := range posts {
		if post.Slug == slug && post.Created().Unix() == ct {
			return post
		}
	}
	return nil
}

// QueryPage query page with slug and created time
func QueryPage(pages []*object.Page, slug string, ct int64) *object.Page {
	for _, page := range pages {
		if page.URL == slug && page.Created().Unix() == ct {
			return page
		}
	}
	return nil
}

// CheckPostSlugConflict check post-slug is conflict with other post
func CheckPostSlugConflict(posts []*object.Post, slug string, ct int64) bool {
	for _, post := range posts {
		if post.Slug == slug {
			if ct == 0 {
				return true
			}
			if post.Created().Unix() == ct {
				return false // itself
			}
			return true
		}
	}
	return false
}
