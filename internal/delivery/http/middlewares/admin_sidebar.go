package middlewares

import (
	"net/http"

	"github.com/devetek/go-core/render"
)

type SidebarItem struct {
	Name string
	Icon string
	Link string
}

// If you need to add extra context in the request, please create public context KEY
// https://vishnubharathi.codes/blog/context-with-value-pitfall/
func AdminSidebar(view *render.Engine) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: move sidebar creation to external data source (CMS)
			var sidebarItem []SidebarItem

			sidebarItem = append(sidebarItem, SidebarItem{
				Name: "Dashboard",
				Icon: `<svg class="w-6 h-6 text-gray-800 dark:text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
  <path fill="currentColor" fill-rule="evenodd" d="M3.005 12 3 6.408l6.8-.923v6.517H3.005ZM11 5.32 19.997 4v8H11V5.32ZM20.067 13l-.069 8-9.065-1.275L11 13h9.067ZM9.8 19.58l-6.795-.931V13H9.8v6.58Z" clip-rule="evenodd"/>
</svg>`,
				Link: "/admin",
			}, SidebarItem{
				Name: "Member",
				Icon: `<svg class="w-6 h-6 text-gray-800 dark:text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 24 24">
  <path fill-rule="evenodd" d="M12 4a4 4 0 1 0 0 8 4 4 0 0 0 0-8Zm-2 9a4 4 0 0 0-4 4v1a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2v-1a4 4 0 0 0-4-4h-4Z" clip-rule="evenodd"/>
</svg>`,
				Link: "/admin/member",
			})

			view.Set("sidebars", sidebarItem)

			next.ServeHTTP(w, r)
		})
	}
}
