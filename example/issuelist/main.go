package main

// ############################################################################
//
//   This code is a result of RnD for gitlab API
//   It's purpose is to list issues required for assignee attention
//   criteria:
//     * open issue
//     * assigned to specific gitlab username
//     * have no notes or last note is not from assignee
//   This code is far from perfect. There is no parallel requests, timeouts, etc
//   It is not covered with tests, so use it only as an example of working with
//   gitlab API from golang
//
// ############################################################################

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/xanzy/go-gitlab"
)

// ############################################################################
// Parse flags
// ############################################################################
type GroupArray []string

func (array *GroupArray) String() string {
	return fmt.Sprintf("[%s]", strings.Join(*array, ","))
}

func (array *GroupArray) Set(value string) error {
	*array = append(*array, value)
	return nil
}
func (array *GroupArray) Type() string {
	return "GroupArray"
}

type CmdOptions struct {
	TokenFilePath string
	GroupIDList   GroupArray
	AssignUser    string
}

func ParseCommandLineOrDie() *CmdOptions {
	opt := &CmdOptions{}

	flag.Var(&opt.GroupIDList, "group", "gitlab group id to look for issues in")
	flag.StringVar(&opt.TokenFilePath, "token-file", "~/.issuemon/gitlab-token", "read gitlab secret token from this file")
	flag.StringVar(&opt.AssignUser, "username", "mxpatlas", "search issues assigned to this gitlab username")

	flag.Parse()

	if len(opt.GroupIDList) == 0 {
		opt.GroupIDList = GroupArray{"11827547", "11536035", "12255376"}
	}

	return opt
}

// ############################################################################
// Gitlab API
// ############################################################################
func GitlabToken(opt *CmdOptions) (string, error) {
	path := opt.TokenFilePath
	if strings.HasPrefix(path, "~/") {
		curUser, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("user.Current() error: %v", err)
		}
		path = filepath.Join(curUser.HomeDir, path[2:])
	}

	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", fmt.Errorf("Failed to make absolute path from %q: %v", path, err)
		}
		path = absPath
	}

	gitlabTokenBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("Failed to read gitlab token from file %q: %v", path, err)
	}
	gitlabToken := strings.TrimSuffix(string(gitlabTokenBytes), "\n")
	return gitlabToken, nil
}

func GitlabClient(opt *CmdOptions) (*gitlab.Client, error) {
	gitlabToken, err := GitlabToken(opt)
	if err != nil {
		return nil, err
	}

	client, err := gitlab.NewClient(gitlabToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %v", err)
	}
	return client, nil
}

func UserIssues(client *gitlab.Client) ([]*gitlab.Issue, error) {
	searchOpt := gitlab.ListIssuesOptions{}

	issues, resp, err := client.Issues.ListIssues(&searchOpt)
	if err != nil {
		return nil, fmt.Errorf("ListIssues error: %v", err)
	}
	_ = resp

	return issues, nil
}

func GroupIssues(client *gitlab.Client, gid interface{}, opt *CmdOptions) ([]*gitlab.Issue, error) {
	state := "opened"

	searchOpt := &gitlab.ListGroupIssuesOptions{
		AssigneeUsername: &opt.AssignUser,
		State:            &state,
		ListOptions:      gitlab.ListOptions{PerPage: 100},
	}
	allIssues := []*gitlab.Issue{}
	for {
		issues, resp, err := client.Issues.ListGroupIssues(gid, searchOpt)
		if err != nil {
			return nil, fmt.Errorf("ListGroupIssues[%d] error: %v", searchOpt.Page, err)
		}
		//log.Printf("group issues resp: %+v", resp)
		allIssues = append(allIssues, issues...)
		if resp.CurrentPage == resp.TotalPages {
			break
		}
		searchOpt.ListOptions.Page = resp.NextPage
	}
	return allIssues, nil
}

type IssueTree struct {
	Issue *gitlab.Issue
	Notes []*gitlab.Note
}

func (it IssueTree) Print(prefix string) {
	log.Printf("%s%+v", prefix, it.Issue.WebURL)
	log.Printf("%s  Wait time: %.0fh; Author: %q; UserNotesCount: %d",
		prefix,
		time.Since(*it.Issue.UpdatedAt).Hours(),
		it.Issue.Author.Name,
		it.Issue.UserNotesCount,
	)
}

func (it *IssueTree) GetNotes(client *gitlab.Client) error {
	if it.Issue.UserNotesCount == 0 {
		return nil
	}
	searchOpt := &gitlab.ListIssueNotesOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100},
	}
	allNotes := []*gitlab.Note{}
	for {
		notes, resp, err := client.Notes.ListIssueNotes(it.Issue.ProjectID, it.Issue.IID, searchOpt)
		if err != nil {
			return fmt.Errorf("ListIssueNotes[ProjectId=%v][IID=%v][Page=%v] error: %v",
				it.Issue.ProjectID,
				it.Issue.IID,
				searchOpt.Page,
				err,
			)
		}
		allNotes = append(allNotes, notes...)
		if resp.CurrentPage == resp.TotalPages {
			break
		}
		searchOpt.ListOptions.Page = resp.NextPage
	}
	sort.Slice(allNotes, func(i, j int) bool {
		return allNotes[i].CreatedAt.Before(*allNotes[j].CreatedAt)
	})
	it.Notes = allNotes
	return nil
}

// ############################################################################
// main()
// ############################################################################
func main() {
	opt := ParseCommandLineOrDie()

	client, err := GitlabClient(opt)
	if err != nil {
		log.Fatalf("gitlab client error: %v", err)
	}

	allUserIssues := []*gitlab.Issue{}

	for _, groupID := range opt.GroupIDList {
		issues, err := GroupIssues(client, groupID, opt)
		if err != nil {
			log.Fatalf("GroupIssues error: %v", err)
		}
		allUserIssues = append(allUserIssues, issues...)
	}
	sort.Slice(allUserIssues, func(i, j int) bool {
		return allUserIssues[i].UpdatedAt.After(*allUserIssues[j].UpdatedAt)
	})

	issueTrees := make([]*IssueTree, 0, len(allUserIssues))
	for i, issue := range allUserIssues {
		it := IssueTree{Issue: issue}
		issuePrefix := fmt.Sprintf("issue[%d]: ", i)
		if err := it.GetNotes(client); err != nil {
			log.Fatalf("%s GetNotes() error: %v", issuePrefix, err)
		}
		if len(it.Notes) == 0 || it.Notes[len(it.Notes)-1].Author.Username != opt.AssignUser {
			issueTrees = append(issueTrees, &it)
		}
	}

	if len(issueTrees) > 0 {
		for i, it := range issueTrees {
			issuePrefix := fmt.Sprintf("issue[%d]: ", i)
			it.Print(issuePrefix)
			for noteIndex, note := range it.Notes {
				log.Printf("%s Note[%d] %v %+v", issuePrefix, noteIndex, note.CreatedAt, note.Author.Username)
			}
		}
	} else {
		log.Printf("No issues waiting for user %q in groups %s", opt.AssignUser, opt.GroupIDList.String())
	}
}
