package commands

import (
    "github.com/bwmarrin/discordgo"
    "github.com/samber/lo"
)

func MentionedBot(s *discordgo.Session, i *discordgo.MessageCreate) bool {
    if s.State == nil && s.State.User == nil {
        return false
    }

    if len(i.Mentions) == 0 || i.Author.Bot {
        return false
    }

    return lo.ContainsBy(i.Mentions, func(item *discordgo.User) bool {
        return item.ID == s.State.User.ID
    })
}

func MentionHandler(s *discordgo.Session, i *discordgo.MessageCreate) {
    if i.Author.ID == s.State.User.ID {
        return // Dont respond to messages from self
    }

    targetUserId := i.Author.ID

    activeBrb, err := Manager.GetActiveBrb(targetUserId)
    if err != nil || activeBrb.IsFinished() {
        BrbMentionHandler(s, i)
    } else {
        BackMentionHandler(s, i)
    }
}
