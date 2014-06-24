package notify

import (
	"fmt"
	"net/url"
)

func getBuildUrl(context *Context) string {
	branchQuery := url.Values{}
	if context.Commit.Branch != "" {
		branchQuery.Set("branch", context.Commit.Branch)
	}

	return fmt.Sprintf("%s/%s/commit/%s?%s", context.Host, context.Repo.Slug, context.Commit.Hash, branchQuery.Encode())
}
