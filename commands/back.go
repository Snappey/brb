package commands

import (
    "brb/manager"
    "fmt"
    "github.com/bwmarrin/discordgo"
    "github.com/rs/zerolog/log"
    "time"
)

type HandlerInput struct {
    User    *discordgo.User
    Message *discordgo.Message
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

    confirmResult, err := createBackConfirmationAndWait(s, i.Message, activeBrb)
    log.Printf("back confirmation result %s", confirmResult)
    switch confirmResult {
    case BackConfirmYes:
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
    case BackConfirmNo:
        err := activeBrb.Start()
        if err != nil {
            return HandlerOutput{
                Content: fmt.Sprintf("failed to restart %s brb session", i.User.Mention()),
            }, err
        }
        return HandlerOutput{
            Content: fmt.Sprintf("%s is not actually back!? restarting ession",
                i.User.Mention(),
            ),
            Flags: 0,
        }, nil
    case BackConfirmGone:
        err := activeBrb.Start()
        if err != nil {
            return HandlerOutput{
                Content: fmt.Sprintf("failed to restart %s brb session", i.User.Mention()),
            }, err
        }
        return HandlerOutput{
            Content: fmt.Sprintf("%s was back now they're gone again!? restarting ession",
                i.User.Mention(),
            ),
            Flags: 0,
        }, nil
    case BackConfirmError:
        return HandlerOutput{
            Content: fmt.Sprintf("back confirmation for %s failed, %s",
                i.User.Mention(),
                err,
            ),
            Flags: 0,
        }, nil
    case BackConfirmTimedOut:
        activeBrb.Finish()

        response := fmt.Sprintf("no one confirmed %s :( we'll trust you and finish anyways.. you took %s (target: %s, difference: %s)",
            i.User.Mention(),
            activeBrb.FinishedDuration.String(),
            activeBrb.TargetDuration.String(),
            (activeBrb.FinishedDuration - activeBrb.TargetDuration).String(),
        )

        if activeBrb.UserId != activeBrb.ReportingUserId {
            err := Manager.CreateBrb(manager.CreateBrbInput{
                TargetUserId:    activeBrb.ReportingUserId,
                ReportingUserId: activeBrb.ReportingUserId,
                TargetDuration:  time.Minute * 5,
            })
            if err != nil {
                log.Err(err).Msg("failed to create brb session for reporting user who did not confirm back")
            } else {
                response = fmt.Sprintf("%s p.s we've started a brb session for <@%s>", response, activeBrb.ReportingUserId)
            }
        }

        return HandlerOutput{
            Content: response,
            Flags:   0,
        }, nil
    default:
        log.Warn().Str("back_confirmation", string(confirmResult)).Msg("unimplemented back confirmation type")
        return HandlerOutput{
            Content: fmt.Sprintf("back confirmation for %s failed, something went very wrong",
                i.User.Mention(),
            ),
        }, nil
    }
}

func BackChatHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
    out, err := BackHandler(s, HandlerInput{User: i.Member.User, Message: i.Message})
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
    out, err := BackHandler(s, HandlerInput{User: i.Author, Message: i.Message})
    if err != nil {
        log.Printf("`back` command handler failed err=%v", err)
    }

    _, sentErr := s.ChannelMessageSendReply(i.ChannelID, out.Content, i.Message.Reference())
    if sentErr != nil {
        log.Printf("error sending `back` response to user channel_id=%s err=%s", i.ChannelID, sentErr)
    }
}
