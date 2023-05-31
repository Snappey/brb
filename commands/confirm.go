package commands

import (
    "brb/manager"
    "fmt"
    "github.com/bwmarrin/discordgo"
    "github.com/rs/zerolog/log"
    "time"
)

type BackConfirm string

const (
    BackConfirmYes      BackConfirm = "back_confirm_yes"
    BackConfirmNo       BackConfirm = "back_confirm_no"
    BackConfirmGone     BackConfirm = "back_confirm_gone"
    BackConfirmTimedOut BackConfirm = "back_confirm_timeout"
    BackConfirmError    BackConfirm = "back_confirm_error"
)

var ComponentsHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}

func createCustomId(brb *manager.BrbSession, key BackConfirm) string {
    return fmt.Sprintf("%s:%s", brb.Id, key)
}

func createComponentInteraction(brb *manager.BrbSession, channel chan<- BackConfirm, result BackConfirm) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
    return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
        if i.Member.User.Bot {
            return
        }
        buttonUserId := i.Member.User.ID

        if buttonUserId == brb.UserId {
            _ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
                Type: discordgo.InteractionResponseChannelMessageWithSource,
                Data: &discordgo.InteractionResponseData{
                    Content: fmt.Sprintf("you can't confirm your own return <@%s>", brb.UserId),
                    Flags:   discordgo.MessageFlagsEphemeral,
                },
            })
        }

        nonReportingUserTime := brb.LastUpdated.Add(time.Minute * 2)
        if buttonUserId != brb.ReportingUserId && nonReportingUserTime.Before(time.Now().UTC()) {
            _ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
                Type: discordgo.InteractionResponseChannelMessageWithSource,
                Data: &discordgo.InteractionResponseData{
                    Content: fmt.Sprintf("you can't confirm <@%s> return yet, needs <@%s> confirmation. (try again in %s)",
                        brb.UserId,
                        brb.ReportingUserId,
                        nonReportingUserTime.Sub(time.Now().UTC()).String(),
                    ),
                    Flags: discordgo.MessageFlagsEphemeral,
                },
            })
        }

        channel <- result
    }
}

func cleanupHandlers(Ids []string) {
    for _, id := range Ids {
        delete(ComponentsHandler, id)
    }
}

func createBackConfirmationAndWait(s *discordgo.Session, m *discordgo.Message, brb *manager.BrbSession) (BackConfirm, error) {
    if !brb.IsActive() {
        return BackConfirmError, fmt.Errorf("session is not active")
    }

    err := brb.AwaitConfirmation()
    if err != nil {
        return BackConfirmError, fmt.Errorf("session failed to update span")
    }

    var callback chan BackConfirm
    yesId, noId, goneId := createCustomId(brb, BackConfirmYes), createCustomId(brb, BackConfirmNo), createCustomId(brb, BackConfirmGone)
    ComponentsHandler[yesId] = createComponentInteraction(brb, callback, BackConfirmYes)
    ComponentsHandler[noId] = createComponentInteraction(brb, callback, BackConfirmNo)
    ComponentsHandler[goneId] = createComponentInteraction(brb, callback, BackConfirmGone)

    defer cleanupHandlers([]string{yesId, noId, goneId})

    returnMsg := fmt.Sprintf("confirm <@%s> is back (you have 5 minutes)", brb.UserId)
    if brb.UserId != brb.ReportingUserId {
        returnMsg = fmt.Sprintf("%s (only <@%s> can confirm in the first 2 minutes)", returnMsg, brb.ReportingUserId)
    }

    msg, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
        Content: returnMsg,
        TTS:     false,
        Components: []discordgo.MessageComponent{
            discordgo.ActionsRow{
                Components: []discordgo.MessageComponent{
                    discordgo.Button{
                        Label:    "Yes",
                        Style:    discordgo.SuccessButton,
                        Disabled: false,
                        CustomID: yesId,
                    },
                    discordgo.Button{
                        Label:    "No",
                        Style:    discordgo.DangerButton,
                        Disabled: false,
                        CustomID: noId,
                    },
                    discordgo.Button{
                        Label:    "They're gone again",
                        Style:    discordgo.SecondaryButton,
                        Disabled: false,
                        CustomID: goneId,
                    },
                },
            },
        },
        Reference: m.Reference(),
    })
    if err != nil {
        log.Err(err).Interface("brb_session", brb).Msg("failed to send confirmation interaction")
        return BackConfirmError, err
    }

    defer func() {
        _ = s.ChannelMessageDelete(msg.ChannelID, m.ID)
    }()

    select {
    case res := <-callback:
        return res, nil
    case <-time.After(time.Minute * 5):
        return BackConfirmTimedOut, nil
    }
}
