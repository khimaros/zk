package note

import (
	"fmt"
	"testing"

	"github.com/mickael-menu/zk/core/templ"
	"github.com/mickael-menu/zk/core/zk"
	"github.com/mickael-menu/zk/util/opt"
	"github.com/mickael-menu/zk/util/test/assert"
)

func TestCreate(t *testing.T) {
	filenameTemplate := NewRendererSpyString("filename")
	bodyTemplate := NewRendererSpyString("body")

	res, err := create(
		CreateOpts{
			Dir: zk.Dir{
				Name: "log",
				Path: "/test/log",
				Config: zk.DirConfig{
					Extension: "md",
					Extra: map[string]string{
						"hello": "world",
					},
				},
			},
			Title:   opt.NewString("Note title"),
			Content: opt.NewString("Note content"),
		},
		createDeps{
			filenameTemplate: filenameTemplate,
			bodyTemplate:     bodyTemplate,
			genId:            func() string { return "abc" },
			validatePath:     func(path string) (bool, error) { return true, nil },
			now:              Now,
		},
	)

	// Check the created note.
	assert.Nil(t, err)
	assert.Equal(t, res, &createdNote{
		path:    "/test/log/filename.md",
		content: "body",
	})

	// Check that the templates received the proper render contexts.
	assert.Equal(t, filenameTemplate.Contexts, []interface{}{renderContext{
		ID:      "abc",
		Title:   "Note title",
		Content: "Note content",
		Dir:     "log",
		Extra: map[string]string{
			"hello": "world",
		},
		Now: Now,
	}})
	assert.Equal(t, bodyTemplate.Contexts, []interface{}{renderContext{
		ID:           "abc",
		Title:        "Note title",
		Content:      "Note content",
		Dir:          "log",
		Filename:     "filename.md",
		FilenameStem: "filename",
		Extra: map[string]string{
			"hello": "world",
		},
		Now: Now,
	}})
}

func TestCreateTriesUntilValidPath(t *testing.T) {
	filenameTemplate := NewRendererSpy(func(context interface{}) string {
		return context.(renderContext).ID
	})
	bodyTemplate := NewRendererSpyString("body")

	res, err := create(
		CreateOpts{
			Dir: zk.Dir{
				Name: "log",
				Path: "/test/log",
				Config: zk.DirConfig{
					Extension: "md",
				},
			},
			Title: opt.NewString("Note title"),
		},
		createDeps{
			filenameTemplate: filenameTemplate,
			bodyTemplate:     bodyTemplate,
			genId:            incrementingID(),
			validatePath: func(path string) (bool, error) {
				return path == "/test/log/3.md", nil
			},
			now: Now,
		},
	)

	// Check the created note.
	assert.Nil(t, err)
	assert.Equal(t, res, &createdNote{
		path:    "/test/log/3.md",
		content: "body",
	})

	assert.Equal(t, filenameTemplate.Contexts, []interface{}{
		renderContext{
			ID:    "1",
			Title: "Note title",
			Dir:   "log",
			Now:   Now,
		},
		renderContext{
			ID:    "2",
			Title: "Note title",
			Dir:   "log",
			Now:   Now,
		},
		renderContext{
			ID:    "3",
			Title: "Note title",
			Dir:   "log",
			Now:   Now,
		},
	})
}

func TestCreateErrorWhenNoValidPaths(t *testing.T) {
	_, err := create(
		CreateOpts{
			Dir: zk.Dir{
				Name: "log",
				Path: "/test/log",
				Config: zk.DirConfig{
					Extension: "md",
				},
			},
		},
		createDeps{
			filenameTemplate: templ.RendererFunc(func(context interface{}) (string, error) {
				return "filename", nil
			}),
			bodyTemplate: templ.NullRenderer,
			genId:        func() string { return "abc" },
			validatePath: func(path string) (bool, error) { return false, nil },
			now:          Now,
		},
	)

	assert.Err(t, err, "/test/log/filename.md: note already exists")
}

// incrementingID returns a generator of incrementing string ID.
func incrementingID() func() string {
	i := 0
	return func() string {
		i++
		return fmt.Sprintf("%d", i)
	}
}
