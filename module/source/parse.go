package source

import "github.com/fuxiaohei/pugov1/module/object"

// Parse parse source contents to data
func Parse(s *object.Source, withDraft bool) error {
	var err error
	s.Posts, err = parsePosts(s, withDraft)
	if err != nil {
		return err
	}
	s.Lists = buildPostLists(s.Posts)
	s.ListSize = len(s.Lists) + 1
	s.TagsLists = buildTagLists(s.Posts)
	s.Archives = buildArchive(s.Posts)

	s.Pages, err = parsePages(s, withDraft)
	if err != nil {
		return err
	}
	return nil
}

const (
	pageSize = 5
)

func buildPostLists(posts []*object.Post) []object.PostList {
	var lists []object.PostList
	cursor := object.NewPagerCursor(pageSize, len(posts))
	for i := 1; ; i++ {
		pg := cursor.Page(i)
		if pg == nil {
			break
		}
		list := object.PostList{
			Pager: pg,
			Posts: posts[pg.Begin:pg.End],
		}
		lists = append(lists, list)
	}
	return lists
}

func buildTagLists(posts []*object.Post) map[string]*object.TagList {
	lists := make(map[string]*object.TagList)
	for _, p := range posts {
		for _, tag := range p.Tags {
			list := lists[tag.Name]
			if list == nil {
				list = &object.TagList{
					Tag:   tag,
					Posts: []*object.Post{p},
				}
				lists[tag.Name] = list
				continue
			}
			list.Posts = append(list.Posts, p)
		}
	}
	for _, list := range lists {
		list.Archives = buildArchive(list.Posts)
	}
	return lists
}

func buildArchive(posts []*object.Post) []*object.Archive {
	var archives []*object.Archive
	var (
		last, lastYear int
	)
	for _, p := range posts {
		if len(archives) == 0 {
			archives = append(archives, &object.Archive{
				Year:  p.Created().Year(),
				Posts: []*object.Post{p},
			})
			continue
		}
		last = len(archives) - 1
		lastYear = archives[last].Year
		if lastYear == p.Created().Year() {
			archives[last].Posts = append(archives[last].Posts, p)
			continue
		}
		archives = append(archives, &object.Archive{
			Year:  p.Created().Year(),
			Posts: []*object.Post{p},
		})
	}
	return archives
}
