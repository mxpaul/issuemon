package main

// ############################################################################
//
//   This code is a result of RnD for gitlab API
//   It's purpose is to list groups for current user
//   This code is far from perfect. There is no parallel requests, timeouts, etc
//   It is not covered with tests, so use it only as an example of working with
//   gitlab API from golang
//
// ############################################################################

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/xanzy/go-gitlab"
)

// ############################################################################
// Parse flags
// ############################################################################

type CmdOptions struct {
	TokenFilePath string
}

func ParseCommandLineOrDie() *CmdOptions {
	opt := &CmdOptions{}

	flag.StringVar(&opt.TokenFilePath, "token-file", "./gitlab.token.txt", "read gitlab secret token from this file")

	flag.Parse()

	return opt
}

// ############################################################################
// Gitlab API
// ############################################################################
func GitlabToken(opt *CmdOptions) (string, error) {
	gitlabTokenBytes, err := ioutil.ReadFile(opt.TokenFilePath)
	if err != nil {
		return "", fmt.Errorf("Failed to read gitlab token from file %q: %v", opt.TokenFilePath, err)
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

//func UserIssues(client *gitlab.Client) ([]*gitlab.Issue, error) {
//	searchOpt := gitlab.ListIssuesOptions{}
//
//	issues, resp, err := client.Issues.ListIssues(&searchOpt)
//	if err != nil {
//		return nil, fmt.Errorf("ListIssues error: %v", err)
//	}
//	_ = resp
//
//	return issues, nil
//}

func UserGroups(client *gitlab.Client) ([]*gitlab.Group, error) {
	groups, resp, err := client.Groups.ListGroups(nil)
	_ = resp
	if err != nil {
		return nil, fmt.Errorf("ListGroups error: %v", err)
	}
	return groups, nil
}

func DumpGroups(groups []*gitlab.Group) {
	log.Printf("ListGroups TotalItems: %v", len(groups))
	for i, group := range groups {
		//log.Printf("group[%v]: ID %v Name: %v URL: %v", i, group.ID, group.Name, group.WebURL)
		log.Printf("group[%v]: ID %v Name: %v", i, group.ID, group.Name)
	}
}

//func GroupIssueBoards(client *gitlab.Client, group *gitlab.Group) ([]*gitlab.GroupIssueBoard, error) {
//	boards, resp, err := client.GroupIssueBoards.ListGroupIssueBoards(group.ID, nil)
//	if err != nil {
//		return nil, fmt.Errorf("ListGroupIssueBoards[ID=%v][%v] error: %v", group.ID, group.Name, err)
//	}
//	_ = resp
//	return boards, nil
//}

//func DumpIssueBoards(group *gitlab.Group, boards []*gitlab.GroupIssueBoard) {
//	log.Printf("ListGroupIssueBoards[ID=%v][%v] TotalItems: %v", group.ID, group.Name, len(boards))
//	for i, board := range boards {
//		log.Printf("Group [ID=%v][%v] IssueBoard[%v]: [%+v] %v", group.ID, group.Name, i, board.ID, board.Name)
//	}
//}

// ############################################################################
// main()
// ############################################################################
func main() {
	opt := ParseCommandLineOrDie()

	client, err := GitlabClient(opt)
	if err != nil {
		log.Fatalf("gitlab client error: %v", err)
	}

	//if _, err = UserIssues(client); err != nil {
	//	log.Fatalf("UserIssues error: %v", err)
	//}

	groups, err := UserGroups(client)
	if err != nil {
		log.Fatalf("UserGroups error: %v", err)
	}
	DumpGroups(groups)

	//for _, group := range groups {
	//	boards, err := GroupIssueBoards(client, group)
	//	if err != nil {
	//		log.Printf("GroupIssueBoards [ID=%v][%v] error: %v", group.ID, group.Name, err)
	//	}
	//	DumpIssueBoards(group, boards)
	//}
}
