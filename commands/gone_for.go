package commands

import (
    "brb/util"
    "fmt"
    "github.com/bwmarrin/discordgo"
)

func GoneChatHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
    options := i.ApplicationCommandData().Options
    optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
    for _, opt := range options {
        optionMap[opt.Name] = opt
    }

    targetUser := i.Member.User
    if goneUser, ok := optionMap[brbMentionKey]; ok {
        targetUser = goneUser.UserValue(s)
    }

    brbSession, err := Manager.GetActiveBrb(targetUser.ID)
    if err != nil {
        _ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: fmt.Sprintf("user does not have an active brb session"),
                Flags:   discordgo.MessageFlagsEphemeral,
            },
        })
        return
    }

    _ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Content: fmt.Sprintf("<@%s> has been gone for %s",
                brbSession.UserId,
                util.HumanizeDuration(brbSession.GetDuration())),
            Flags: 0,
        },
    })
}
