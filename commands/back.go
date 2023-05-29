package commands

import (
    "fmt"
    "github.com/bwmarrin/discordgo"
    "github.com/rs/zerolog/log"
)

type HandlerInput struct {
    User *discordgo.User
}

type HandlerOutput struct {
    Content string
    Flags   discordgo.MessageFlags
}

func BackHandler(s *discordgo.Session, i HandlerInput) (HandlerOutput, error) {
    // Check if user has an active brb
    activeBrb, err := Manager.GetActiveBrb(i.User.ID)
    if err != nil {
        return HandlerOutput{
            Content: fmt.Sprintf("failed to mark as back, %s", err),
            Flags:   discordgo.MessageFlagsEphemeral,
        }, err
    }

    // Finish brb
    activeBrb.Finish()

    return HandlerOutput{
        Content: fmt.Sprintf("welcome back %s, you took %s (target: %s, difference: %s)",
            i.User.Mention(),
            activeBrb.FinishedDuration.String(),
            activeBrb.TargetDuration.String(),
            (activeBrb.FinishedDuration - activeBrb.TargetDuration).String(),
        ),
        Flags: 0,
    }, err
}

func BackChatHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
    out, err := BackHandler(s, HandlerInput{User: i.Member.User})
    if err != nil {
        log.Printf("`back` command handler failed err=%v", err)
    }

    _ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Content: out.Content,
            Flags:   out.Flags,
        },
    })
}

func BackMentionHandler(s *discordgo.Session, i *discordgo.MessageCreate) {
    out, err := BackHandler(s, HandlerInput{User: i.Author})
    if err != nil {
        log.Printf("`back` command handler failed err=%v", err)
    }

    _, sentErr := s.ChannelMessageSendReply(i.ChannelID, out.Content, i.Message.Reference())
    if sentErr != nil {
        log.Printf("error sending `back` response to user channel_id=%s err=%s", i.ChannelID, sentErr)
    }
}
