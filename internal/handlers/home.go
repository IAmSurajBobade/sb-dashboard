package handlers

import (
	"embed"
	"html/template"
	"net/http"
)

func (c *Controller) HomeHandler(content embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.mutex.Lock()
		c.userCount++
		cnt := c.userCount
		c.mutex.Unlock()

		data := map[string]interface{}{
			"Count":      cnt,
			"TodaysDate": c.now().Format("2006-01-02"),
		}

		tmpl, _ := template.ParseFS(content, "templates/home.html")
		tmpl.Execute(w, data)
	}
}
