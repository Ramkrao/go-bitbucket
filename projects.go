package bitbucket

import (
	"fmt"
	"net/url"
	"strconv"
)

//"github.com/k0kubun/pp"

func (r *Repositories) ListForProject(ro *RepositoriesOptions) (*RepositoriesRes, error) {
	urlStr := r.c.requestUrl("/projects/%s/repos", ro.Owner)
	if ro.Role != "" {
		urlStr += "?role=" + ro.Role
	}
	if ro.Limit != "" {
		urlStr += "?limit=" + ro.Limit
	}
	fmt.Println(urlStr)
	repos, err := r.c.execute("GET", urlStr, "")
	fmt.Println(repos)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return decodeRepositorys(repos)
}

func (r *Repository) ListBranchesForProject(rbo *RepositoryBranchOptions) (*RepositoryBranches, error) {

	params := url.Values{}
	if rbo.Query != "" {
		params.Add("q", rbo.Query)
	}

	if rbo.Sort != "" {
		params.Add("sort", rbo.Sort)
	}

	if rbo.PageNum > 0 {
		params.Add("page", strconv.Itoa(rbo.PageNum))
	}

	if rbo.Pagelen > 0 {
		params.Add("pagelen", strconv.Itoa(rbo.Pagelen))
	}

	urlStr := r.c.requestUrl("/projects/%s/repos/%s/branches?%s", rbo.Owner, rbo.RepoSlug, params.Encode())
	fmt.Println(urlStr)
	response, err := r.c.executeRaw("GET", urlStr, "")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return decodeRepositoryBranches(response)
}

func (cm *Commits) GetCommitsForProject(cmo *CommitsOptions) (interface{}, error) {
	urlStr := cm.c.requestUrl("/projects/%s/repos/%s/commits/%s", cmo.Owner, cmo.RepoSlug, cmo.Branchortag)
	urlStr += cm.buildCommitsQuery(cmo.Include, cmo.Exclude)
	return cm.c.execute("GET", urlStr, "")
}

func (p *PullRequests) GetsForProject(po *PullRequestsOptions) (interface{}, error) {
	urlStr := GetApiBaseURL() + "/projects/" + po.Owner + "/repos/" + po.RepoSlug + "/pull-requests/"

	if po.States != nil && len(po.States) != 0 {
		parsed, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		query := parsed.Query()
		for _, state := range po.States {
			query.Set("state", state)
		}
		parsed.RawQuery = query.Encode()
		urlStr = parsed.String()
	}

	if po.Query != "" {
		parsed, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		query := parsed.Query()
		query.Set("q", po.Query)
		parsed.RawQuery = query.Encode()
		urlStr = parsed.String()
	}

	if po.Sort != "" {
		parsed, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		query := parsed.Query()
		query.Set("sort", po.Sort)
		parsed.RawQuery = query.Encode()
		urlStr = parsed.String()
	}

	return p.c.execute("GET", urlStr, "")
}
