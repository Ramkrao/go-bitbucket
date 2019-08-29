package bitbucket

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

//"github.com/k0kubun/pp"

type ProjectRepositoryBranches struct {
	Page     int
	Pagelen  int
	Size     int
	Next     string
	Branches []ProjectRepositoryBranch
}

type ProjectRepositoryBranch struct {
	Type            string
	DisplayId       string
	ID              string
	IsDefault       bool
	LatestChangeset string
	LatestCommit    string
}

func (r *Repositories) ListForProject(ro *RepositoriesOptions) (*RepositoriesRes, error) {
	urlStr := r.c.requestUrl("/projects/%s/repos", ro.Owner)
	if ro.Role != "" {
		urlStr += "?role=" + ro.Role
	}
	if ro.Limit != "" {
		urlStr += "?limit=" + ro.Limit
	}
	repos, err := r.c.execute("GET", urlStr, "")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return decodeRepositorys(repos)
}

func (r *Repository) ListBranchesForProject(rbo *RepositoryBranchOptions) (*ProjectRepositoryBranches, error) {

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

	urlStr := r.c.requestUrl("/projects/%s/repos/%s/branches", rbo.Owner, rbo.RepoSlug)
	if params.Encode() != "" {
		urlStr += "?" + params.Encode()
	}
	response, err := r.c.executeRaw("GET", urlStr, "")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return decodeProjectRepositoryBranches(response)
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

func decodeProjectRepositoryBranches(branchResponse interface{}) (*ProjectRepositoryBranches, error) {

	var branchResponseMap map[string]interface{}
	err := json.Unmarshal(branchResponse.([]byte), &branchResponseMap)
	if err != nil {
		return nil, err
	}

	branchArray := branchResponseMap["values"].([]interface{})
	var branches []ProjectRepositoryBranch
	for _, branchEntry := range branchArray {
		var branch ProjectRepositoryBranch
		err = mapstructure.Decode(branchEntry, &branch)
		if err == nil {
			branches = append(branches, branch)
		}
	}

	page, ok := branchResponseMap["page"].(float64)
	if !ok {
		page = 0
	}

	pagelen, ok := branchResponseMap["pagelen"].(float64)
	if !ok {
		pagelen = 0
	}
	size, ok := branchResponseMap["size"].(float64)
	if !ok {
		size = 0
	}

	next, ok := branchResponseMap["next"].(string)
	if !ok {
		next = ""
	}

	repositoryBranches := ProjectRepositoryBranches{
		Page:     int(page),
		Pagelen:  int(pagelen),
		Size:     int(size),
		Next:     next,
		Branches: branches,
	}
	return &repositoryBranches, nil
}
