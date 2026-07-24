package ui

import (
	"strings"

	"github.com/RasmusLindroth/go-mastodon"
	"github.com/RasmusLindroth/tut/api"
	"github.com/RasmusLindroth/tut/util"
)

func itemSearchText(item api.Item) string {
	switch item.Type() {
	case api.StatusType:
		return statusSearchText(item.Raw().(*mastodon.Status))
	case api.StatusHistoryType:
		s := item.Raw().(*mastodon.StatusHistory)
		status := mastodon.Status{
			Content:          s.Content,
			SpoilerText:      s.SpoilerText,
			Account:          s.Account,
			MediaAttachments: s.MediaAttachments,
		}
		return statusSearchText(&status)
	case api.UserType, api.ProfileType:
		return userSearchText(item.Raw().(*api.User))
	case api.NotificationType:
		return notificationSearchText(item.Raw().(*api.NotificationData))
	case api.ListsType:
		return item.Raw().(*mastodon.List).Title
	case api.TagType:
		return item.Raw().(*mastodon.Tag).Name
	}
	return ""
}

func statusSearchText(status *mastodon.Status) string {
	var b strings.Builder
	s := util.StatusOrReblog(status)
	if status.Reblog != nil {
		b.WriteString(util.FormatUsername(status.Account))
		b.WriteString(" ")
	}
	b.WriteString(util.FormatUsername(s.Account))
	b.WriteString(" ")
	b.WriteString(s.SpoilerText)
	b.WriteString(" ")
	content, _ := util.CleanHTML(s.Content)
	b.WriteString(content)
	for _, a := range s.MediaAttachments {
		b.WriteString(" ")
		b.WriteString(a.Description)
	}
	if s.Poll != nil {
		for _, o := range s.Poll.Options {
			b.WriteString(" ")
			b.WriteString(o.Title)
		}
	}
	return b.String()
}

func userSearchText(user *api.User) string {
	var b strings.Builder
	b.WriteString(util.FormatUsername(*user.Data))
	b.WriteString(" ")
	note, _ := util.CleanHTML(user.Data.Note)
	b.WriteString(note)
	for _, f := range user.Data.Fields {
		b.WriteString(" ")
		val, _ := util.CleanHTML(f.Value)
		b.WriteString(val)
	}
	return b.String()
}

func notificationSearchText(nd *api.NotificationData) string {
	var b strings.Builder
	b.WriteString(util.FormatUsername(nd.Item.Account))
	switch nd.Item.Type {
	case "follow", "follow_request":
		b.WriteString(" ")
		b.WriteString(userSearchText(nd.User.Raw().(*api.User)))
	case "favourite", "reblog", "mention", "update", "status", "poll":
		b.WriteString(" ")
		b.WriteString(statusSearchText(nd.Item.Status))
	}
	return b.String()
}

func (tv *TutView) PaneSearchCommand(term string) {
	tv.paneSearchTerm = strings.ToLower(term)
	tv.paneSearchIndex = -1
	tv.NextPaneSearch(true)
}

func (tv *TutView) NextPaneSearch(forward bool) {
	if tv.paneSearchTerm == "" {
		return
	}
	f := tv.GetCurrentFeed()
	items := f.Data.List()
	count := len(items)
	if count == 0 {
		return
	}
	start := f.List.Text.GetCurrentItem()
	if tv.paneSearchIndex >= 0 && tv.paneSearchIndex < count {
		start = tv.paneSearchIndex
	}
	var indices []int
	if forward {
		for i := 1; i <= count; i++ {
			indices = append(indices, (start+i)%count)
		}
	} else {
		for i := 1; i <= count; i++ {
			indices = append(indices, (start-i+count)%count)
		}
	}
	for _, i := range indices {
		text := strings.ToLower(itemSearchText(items[i]))
		if strings.Contains(text, tv.paneSearchTerm) {
			tv.paneSearchIndex = i
			f.List.SetCurrentItem(i)
			f.DrawContent()
			tv.Shared.Bottom.Cmd.ClearInput()
			return
		}
	}
	tv.ShowError("No matches found")
}
