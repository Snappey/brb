package commands

import (
    "brb/core"
    "brb/util"
    "fmt"
    "github.com/bwmarrin/discordgo"
    "github.com/rs/zerolog/log"
    "regexp"
    "time"
)

const DefaultDuration = time.Minute * 5
const MinDuration = time.Minute
const MaxDuration = time.Hour * 6

var durationRegex = regexp.MustCompile("(\\d+[hms]+)+")

type BrbHandlerInput struct {
    HandlerInput
    BrbDuration time.Duration
    TargetUser  *discordgo.User
}

func BrbHandler(s *discordgo.Session, i BrbHandlerInput) (HandlerOutput, error) {
    // mark user as BrbSession
    err := Manager.CreateBrb(core.CreateBrbInput{
        TargetUserId:    i.TargetUser.ID,
        ReportingUserId: i.User.ID,
        TargetDuration:  i.BrbDuration,
    })
    if err != nil {
        return HandlerOutput{
            Content: fmt.Sprintf("failed to create brb, %s", err),
            Flags:   discordgo.MessageFlagsEphemeral,
        }, err
    }

    // Let user know they have a brb
    return HandlerOutput{
        Content: fmt.Sprintf("created brb for %s, see you in %s", i.TargetUser.Mention(), util.HumanizeDuration(i.BrbDuration)),
        Flags:   discordgo.MessageFlagsEphemeral,
    }, err
}

func BrbChatCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
    var err error
    input := BrbHandlerInput{
        HandlerInput: HandlerInput{User: i.Member.User},
        BrbDuration:  DefaultDuration,
        TargetUser:   i.Member.User,
    }

    options := i.ApplicationCommandData().Options
    optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
    for _, opt := range options {
        optionMap[opt.Name] = opt
    }

    if duration, ok := optionMap[brbDurationKey]; ok {
        input.BrbDuration, err = time.ParseDuration(duration.StringValue())
        if err != nil {
            _ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
                Type: discordgo.InteractionResponseChannelMessageWithSource,
                Data: &discordgo.InteractionResponseData{
                    Content: fmt.Sprintf("failed to parse duration, %s", err),
                    Flags:   discordgo.MessageFlagsEphemeral,
                },
            })
        }
    }

    if targetUser, ok := optionMap[brbMentionKey]; ok {
        input.TargetUser = targetUser.UserValue(s)
    }

    if input.BrbDuration < MinDuration || input.BrbDuration > MaxDuration {
        err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: fmt.Sprintf("failed to create brb, duration is outside limits of %d and %d", MinDuration, MaxDuration),
                Flags:   discordgo.MessageFlagsEphemeral,
            },
        })
        if err != nil {
            log.Printf("failed to respond to brb interaction, err=%v", err)
            return
        }
    }

    out, err := BrbHandler(s, input)
    if err != nil {
        log.Printf("`brb` command handler failed err=%v", err)
    }

    _ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Content: out.Content,
            Flags:   out.Flags,
        },
    })
}

func BrbMentionHandler(s *discordgo.Session, i *discordgo.MessageCreate) {
    var err error
    var brbDuration time.Duration

    if duration := durationRegex.FindString(i.Message.Content); duration != "" {
        brbDuration, err = time.ParseDuration(duration)
        if err != nil {
            log.Printf("failed to parse duration err=%v", err)
            _, _ = s.ChannelMessageSendReply(i.ChannelID, fmt.Sprintf("i dont understand, %s", err), i.Message.Reference())
            return
        }
    } else {
        brbDuration = DefaultDuration
    }

    var targets []*discordgo.User
    if len(i.Mentions) == 1 {
        targets = append(targets, i.Author)
    } else {
        for _, mention := range i.Mentions {
            if mention.Bot {
                continue
            }

            targets = append(targets, mention)
        }
    }

    if len(targets) == 0 {
        _, _ = s.ChannelMessageSendReply(i.ChannelID, fmt.Sprintf("failed to find valid afker"), i.Message.Reference())
        return
    }

    for _, target := range targets {
        out, err := BrbHandler(s, BrbHandlerInput{
            HandlerInput: HandlerInput{
                User: i.Author,
            },
            BrbDuration: brbDuration,
            TargetUser:  target,
        })
        if err != nil {
            log.Printf("`brb` command handler failed err=%v", err)
        }

        _, sentErr := s.ChannelMessageSendReply(i.ChannelID, out.Content, i.Message.Reference())
        if sentErr != nil {
            log.Printf("error sending `brb` response to user channel_id=%s err=%s", i.ChannelID, sentErr)
        }
    }
}
