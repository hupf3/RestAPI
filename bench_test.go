package mux

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
)

var muxT *Mux

type api struct {
	method    string
	routeItem string // 路由项
	test      string // 测试地址
}

// 测试GITHUB官方API
var apis = []*api{
	{method: http.MethodGet, routeItem: "/events"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/events"},
	{method: http.MethodGet, routeItem: "/networks/{owner}/{repo}/events"},
	{method: http.MethodGet, routeItem: "/orgs/{org}/events"},
	{method: http.MethodGet, routeItem: "/users/{username}/received_events"},
	{method: http.MethodGet, routeItem: "/users/{username}/received_events/public"},
	{method: http.MethodGet, routeItem: "/users/{username}/events"},
	{method: http.MethodGet, routeItem: "/users/{username}/events/public"},
	{method: http.MethodGet, routeItem: "/users/{username}/events/orgs/{org}"},
	{method: http.MethodGet, routeItem: "/feeds"},
	{method: http.MethodGet, routeItem: "/notifications"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/notifications"},
	{method: http.MethodPut, routeItem: "/notifications"},
	{method: http.MethodGet, routeItem: "/notifications/threads/{id}"},
	{method: http.MethodPatch, routeItem: "/notifications/threads/{id}"},
	{method: http.MethodGet, routeItem: "/notifications/threads/{id}/subscription"},
	{method: http.MethodPut, routeItem: "/notifications/threads/{id}/subscription"},
	{method: http.MethodDelete, routeItem: "/notifications/threads/{id}/subscription"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/stargazers"},
	{method: http.MethodGet, routeItem: "/users/{username}/starred"},
	{method: http.MethodGet, routeItem: "/user/starred"},
	{method: http.MethodGet, routeItem: "/user/starred/{owner}/{repo}"},
	{method: http.MethodPut, routeItem: "/user/starred/{owner}/{repo}"},
	{method: http.MethodDelete, routeItem: "/user/starred/{owner}/{repo}"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/subscribers"},
	{method: http.MethodGet, routeItem: "/users/{username}/subscriptions"},
	{method: http.MethodGet, routeItem: "/user/subscriptions"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/subscription"},
	{method: http.MethodPut, routeItem: "/repos/{owner}/{repo}/subscription"},
	{method: http.MethodDelete, routeItem: "/repos/{owner}/{repo}/subscription"},
	{method: http.MethodGet, routeItem: "/user/subscriptions/{owner}/{repo}"},
	{method: http.MethodPut, routeItem: "/user/subscriptions/{owner}/{repo}"},
	{method: http.MethodDelete, routeItem: "/user/subscriptions/{owner}/{repo}"},
	{method: http.MethodGet, routeItem: "/gists/{gist_id}/comments"},
	{method: http.MethodGet, routeItem: "/gists/{gist_id}/comments/{id}"},
	{method: http.MethodPost, routeItem: "/gists/{gist_id}/comments"},
	{method: http.MethodPatch, routeItem: "/gists/{gist_id}/comments/{id}"},
	{method: http.MethodDelete, routeItem: "/gists/{gist_id}/comments/{id}"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/git/blobs/{sha}"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/git/blobs"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/git/commits/{sha}"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/git/commits"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/git/refs/{ref}"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/git/refs/heads/feature"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/git/refs"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/git/refs/tags"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/git/refs"},
	{method: http.MethodPatch, routeItem: "/repos/{owner}/{repo}/git/refs/{ref}"},
	{method: http.MethodDelete, routeItem: "/repos/{owner}/{repo}/git/refs/{ref}"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/git/tags"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/git/tags/{sha}"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/git/trees/{sha}"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/git/trees"},
	{method: http.MethodGet, routeItem: "/integration/installations"},
	{method: http.MethodGet, routeItem: "/integration/installations/{installation_id}"},
	{method: http.MethodGet, routeItem: "/user/installations"},
	{method: http.MethodPost, routeItem: "/installations/{installation_id}/access_tokens"},
	{method: http.MethodGet, routeItem: "/installation/repositories"},
	{method: http.MethodGet, routeItem: "/user/installations/{installation_id}/repositories"},
	{method: http.MethodPut, routeItem: "/installations/{installation_id}/repositories/{repo}sitory_id"},
	{method: http.MethodDelete, routeItem: "/installations/{installation_id}/repositories/{repo}sitory_id"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/assignees"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/assignees/{assignee}"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/issues/{number}/assignees"},
	{method: http.MethodDelete, routeItem: "/repos/{owner}/{repo}/issues/{number}/assignees"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/issues/{number}/comments"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/issues/comments"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/issues/comments/{id}"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/issues/{number}/comments"},
	{method: http.MethodPatch, routeItem: "/repos/{owner}/{repo}/issues/comments/{id}"},
	{method: http.MethodDelete, routeItem: "/repos/{owner}/{repo}/issues/comments/{id}"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/issues/{issue_number}/events"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/issues/events"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/issues/events/{id}"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/labels"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/labels/{name}"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/labels"},
	{method: http.MethodPatch, routeItem: "/repos/{owner}/{repo}/labels/{name}"},
	{method: http.MethodDelete, routeItem: "/repos/{owner}/{repo}/labels/{name}"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/issues/{number}/labels"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/issues/{number}/labels"},
	{method: http.MethodDelete, routeItem: "/repos/{owner}/{repo}/issues/{number}/labels/{name}"},
	{method: http.MethodPut, routeItem: "/repos/{owner}/{repo}/issues/{number}/labels"},
	{method: http.MethodDelete, routeItem: "/repos/{owner}/{repo}/issues/{number}/labels"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/milestones/{number}/labels"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/milestones"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/milestones/{number}"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/milestones"},
	{method: http.MethodPatch, routeItem: "/repos/{owner}/{repo}/milestones/{number}"},
	{method: http.MethodDelete, routeItem: "/repos/{owner}/{repo}/milestones/{number}"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/issues/{issue_number}/timeline"},
	{method: http.MethodPost, routeItem: "/orgs/{org}/migrations"},
	{method: http.MethodGet, routeItem: "/orgs/{org}/migrations"},
	{method: http.MethodGet, routeItem: "/orgs/{org}/migrations/{id}"},
	{method: http.MethodGet, routeItem: "/orgs/{org}/migrations/{id}/archive"},
	{method: http.MethodDelete, routeItem: "/orgs/{org}/migrations/{id}/archive"},
	{method: http.MethodDelete, routeItem: "/orgs/{org}/migrations/{id}/repos/{repo}_name/lock"},
	{method: http.MethodPut, routeItem: "/repos/{owner}/{repo}/import"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/import"},
	{method: http.MethodPatch, routeItem: "/repos/{owner}/{repo}/import"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/import/authors"},
	{method: http.MethodPatch, routeItem: "/repos/{owner}/{repo}/import/authors/{author_id}"},
	{method: http.MethodPatch, routeItem: "/{owner}/{name}/import/lfs"},
	{method: http.MethodGet, routeItem: "/{owner}/{name}/import/large_files"},
	{method: http.MethodDelete, routeItem: "/repos/{owner}/{repo}/import"},
	{method: http.MethodGet, routeItem: "/emojis"},
	{method: http.MethodGet, routeItem: "/gitignore/templates"},
	{method: http.MethodGet, routeItem: "/gitignore/templates/C"},
	{method: http.MethodGet, routeItem: "/licenses"},
	{method: http.MethodGet, routeItem: "/licenses/{license}"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/license"},
	{method: http.MethodPost, routeItem: "/markdown"},
	{method: http.MethodPost, routeItem: "/markdown/raw"},
	{method: http.MethodGet, routeItem: "/meta"},
	{method: http.MethodGet, routeItem: "/rate_limit"},
	{method: http.MethodGet, routeItem: "/orgs/{org}/members"},
	{method: http.MethodGet, routeItem: "/orgs/{org}/members/{username}"},
	{method: http.MethodDelete, routeItem: "/orgs/{org}/members/{username}"},
	{method: http.MethodGet, routeItem: "/orgs/{org}/public_members"},
	{method: http.MethodGet, routeItem: "/orgs/{org}/public_members/{username}"},
	{method: http.MethodPut, routeItem: "/orgs/{org}/public_members/{username}"},
	{method: http.MethodDelete, routeItem: "/orgs/{org}/public_members/{username}"},
	{method: http.MethodGet, routeItem: "/orgs/{org}/memberships/{username}"},
	{method: http.MethodPut, routeItem: "/orgs/{org}/memberships/{username}"},
	{method: http.MethodDelete, routeItem: "/orgs/{org}/memberships/{username}"},
	{method: http.MethodGet, routeItem: "/orgs/{org}/invitations"},
	{method: http.MethodGet, routeItem: "/user/memberships/orgs"},
	{method: http.MethodGet, routeItem: "/user/memberships/orgs/{org}"},
	{method: http.MethodPatch, routeItem: "/user/memberships/orgs/{org}"},
	{method: http.MethodGet, routeItem: "/orgs/{org}/outside_collaborators"},
	{method: http.MethodDelete, routeItem: "/orgs/{org}/outside_collaborators/{username}"},
	{method: http.MethodPut, routeItem: "/orgs/{org}/outside_collaborators/{username}"},
	{method: http.MethodGet, routeItem: "/orgs/{org}/teams"},
	{method: http.MethodGet, routeItem: "/teams/{id}"},
	{method: http.MethodPost, routeItem: "/orgs/{org}/teams"},
	{method: http.MethodPatch, routeItem: "/teams/{id}"},
	{method: http.MethodDelete, routeItem: "/teams/{id}"},
	{method: http.MethodGet, routeItem: "/teams/{id}/members/{username}"},
	{method: http.MethodPut, routeItem: "/teams/{id}/members/{username}"},
	{method: http.MethodDelete, routeItem: "/teams/{id}/members/{username}"},
	{method: http.MethodGet, routeItem: "/teams/{id}/memberships/{username}"},
	{method: http.MethodPut, routeItem: "/teams/{id}/memberships/{username}"},
	{method: http.MethodDelete, routeItem: "/teams/{id}/memberships/{username}"},
	{method: http.MethodGet, routeItem: "/teams/{id}/repos"},
	{method: http.MethodGet, routeItem: "/teams/{id}/invitations"},
	{method: http.MethodPut, routeItem: "/teams/{id}/repos/{org}/{repo}"},
	{method: http.MethodGet, routeItem: "/orgs/{org}/hooks"},
	{method: http.MethodGet, routeItem: "/orgs/{org}/hooks/{id}"},
	{method: http.MethodPost, routeItem: "/orgs/{org}/hooks"},
	{method: http.MethodPatch, routeItem: "/orgs/{org}/hooks/{id}"},
	{method: http.MethodPost, routeItem: "/orgs/{org}/hooks/{id}/pings"},
	{method: http.MethodDelete, routeItem: "/orgs/{org}/hooks/{id}"},
	{method: http.MethodGet, routeItem: "/orgs/{org}/blocks"},
	{method: http.MethodGet, routeItem: "/orgs/{org}/blocks/{username}"},
	{method: http.MethodPut, routeItem: "/orgs/{org}/blocks/{username}"},
	{method: http.MethodDelete, routeItem: "/orgs/{org}/blocks/{username}"},
	{method: http.MethodGet, routeItem: "/projects/columns/{column_id}/cards"},
	{method: http.MethodGet, routeItem: "/projects/columns/cards/{id}"},
	{method: http.MethodPost, routeItem: "/projects/columns/{column_id}/cards"},
	{method: http.MethodPatch, routeItem: "/projects/columns/cards/{id}"},
	{method: http.MethodDelete, routeItem: "/projects/columns/cards/{id}"},
	{method: http.MethodPost, routeItem: "/projects/columns/cards/{id}/moves"},
	{method: http.MethodGet, routeItem: "/projects/{project_id}/columns"},
	{method: http.MethodPost, routeItem: "/projects/{project_id}/columns"},
	{method: http.MethodDelete, routeItem: "/projects/columns/{id}"},
	{method: http.MethodPost, routeItem: "/projects/columns/{id}/moves"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/pulls/{number}/reviews"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/pulls/{number}/reviews/{id}"},
	{method: http.MethodDelete, routeItem: "/repos/{owner}/{repo}/pulls/{number}/reviews/{id}"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/pulls/{number}/reviews/{id}/comments"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/pulls/{number}/reviews"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/pulls/{number}/reviews/{id}/events"},
	{method: http.MethodPut, routeItem: "/repos/{owner}/{repo}/pulls/{number}/reviews/{id}/dismissals"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/pulls/{number}/comments"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/pulls/comments"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/pulls/comments/{id}"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/pulls/{number}/comments"},
	{method: http.MethodPatch, routeItem: "/repos/{owner}/{repo}/pulls/comments/{id}"},
	{method: http.MethodDelete, routeItem: "/repos/{owner}/{repo}/pulls/comments/{id}"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/pulls/{number}/requested_reviewers"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/pulls/{number}/requested_reviewers"},
	{method: http.MethodDelete, routeItem: "/repos/{owner}/{repo}/pulls/{number}/requested_reviewers"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/comments/{id}/reactions"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/comments/{id}/reactions"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/issues/{number}/reactions"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/issues/{number}/reactions"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/issues/comments/{id}/reactions"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/issues/comments/{id}/reactions"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/pulls/comments/{id}/reactions"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/pulls/comments/{id}/reactions"},
	{method: http.MethodDelete, routeItem: "/reactions/{id}"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/branches"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/branches/{branch}"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection"},
	{method: http.MethodPut, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection"},
	{method: http.MethodDelete, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection"},
	{method: http.MethodPatch, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection/required_status_checks"},
	{method: http.MethodDelete, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection/required_status_checks"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection/required_status_checks/contexts"},
	{method: http.MethodPut, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection/required_status_checks/contexts"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection/required_status_checks/contexts"},
	{method: http.MethodDelete, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection/required_status_checks/contexts"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection/required_pull_request_reviews"},
	{method: http.MethodPatch, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection/required_pull_request_reviews"},
	{method: http.MethodDelete, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection/required_pull_request_reviews"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection/enforce_admins"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection/enforce_admins"},
	{method: http.MethodDelete, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection/enforce_admins"},
	{method: http.MethodDelete, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection/restrictions"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection/restrictions/teams"},
	{method: http.MethodPut, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection/restrictions/teams"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection/restrictions/teams"},
	{method: http.MethodDelete, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection/restrictions/teams"},
	{method: http.MethodGet, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection/restrictions/users"},
	{method: http.MethodPut, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection/restrictions/users"},
	{method: http.MethodPost, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection/restrictions/users"},
	{method: http.MethodDelete, routeItem: "/repos/{owner}/{repo}/branches/{branch}/protection/restrictions/users"},
}

// 初始化
func init() {
	for _, api := range apis {
		path := strings.Replace(api.routeItem, "}", "", -1)
		api.test = strings.Replace(path, "{", "", -1)
	}

	getStates(func() {
		h := func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(r.URL.Path))
		}

		muxT = New(false, true, false, nil, nil)
		for _, api := range apis {
			if err := muxT.HandleFunc(api.routeItem, h, api.method); err != nil {
				fmt.Println("getStates:", err)
			}
		}
	})
}

func getStates(load func()) {
	stats := &runtime.MemStats{}

	runtime.GC()
	runtime.ReadMemStats(stats)
	before := stats.HeapAlloc

	load()

	runtime.GC()
	runtime.ReadMemStats(stats)
	after := stats.HeapAlloc
	fmt.Printf("%d Bytes\n", after-before)
}

func benchGithubAPI(b *testing.B, srv http.Handler) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		api := apis[i%len(apis)]

		w := httptest.NewRecorder()
		r := httptest.NewRequest(api.method, api.test, nil)
		srv.ServeHTTP(w, r)

		if w.Body.String() != r.URL.Path {
			b.Errorf("%s:%s", w.Body.String(), r.URL.Path)
		}
	}
}

func BenchmarkGithubAPI(b *testing.B) {
	benchGithubAPI(b, muxT)
}
