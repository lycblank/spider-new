package github

import (
	"context"
	"fmt"
	githubsdk "github.com/google/go-github/v35/github"
	"github.com/jasonlvhit/gocron"
	"github.com/lycblank/spider-new/pkg/notify"
	"time"
)

var defaultNotify notify.Notify
func Init(n notify.Notify) {
	defaultNotify = n
	s := gocron.NewScheduler()
	s.Every(1).Day().At("09:05").Do(Search, n)
	s.Start()
}

func Search(n notify.Notify) {
	time.Sleep(5*time.Second)
	client := githubsdk.NewClient(nil)
	ret, _, _ := client.Search.Repositories(context.Background(), "language:go", &githubsdk.SearchOptions{
		Sort:"updated",
		Order: "desc",
		ListOptions:githubsdk.ListOptions{
			Page: 1,
			PerPage: 10,
		},
	})
	arg := notify.NotifyArg{
		Title:"github活跃库",
		Contents: make([]string, 0, 10),
	}
	for _, result := range ret.Repositories {
		arg.Contents = append(arg.Contents, fmt.Sprintf("[%s](%s)", result.GetFullName(),
			result.GetSVNURL()))
	}
	n.Send(context.Background(), arg)
}



