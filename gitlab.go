package main

import (
	"time"

	"github.com/dghubble/sling"
)

type Author struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	UserName  string `json:"username"`
	State     string `json:"state"`
	AvatarURL string `json:"avatar_url"`
	WebURL    string `json:"web_url"`
}

type Snippet struct {
	ID          int       `json:"id,omitempty"`
	Title       string    `json:"title"`
	FileName    string    `json:"file_name"`
	Description string    `json:"description"`
	Author      Author    `json:"author,omitempty"`
	UpdateAt    time.Time `json:"updated_at,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	ProjectID   int       `json:"project_id,omitempty"`
	WebURL      string    `json:"web_url,omitempty"`
	RawURL      string    `json:"raw_url,omitempty"`
}

type GitlabPrivateToken struct {
	PrivateToken string `url:"private_token,omitempty"`
}

type GitlabSnippetService struct {
	sling *sling.Sling
}

func newGitlabSnippetService(sling *sling.Sling, privateToken string) *GitlabSnippetService {
	return &GitlabSnippetService{
		sling: sling.
			Base("https://gitlab.com/api/v4/snippets/").
			QueryStruct(&GitlabPrivateToken{PrivateToken: privateToken}),
	}
}

func (s *GitlabSnippetService) getSnippet(snippetID int) (*Snippet, error) {
	return nil, nil
}
