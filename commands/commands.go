package commands

import (
    "brb/config"
    "brb/manager"
    "github.com/bwmarrin/discordgo"
    "github.com/rs/zerolog/log"
)

const brbDurationKey = "duration"
const brbMentionKey = "user"

var Manager = manager.GetInstance()

var ApplicationCommands = []discordgo.ApplicationCommand{
    {
        Name:        "brb",
        Description: "mark yourself as brb",
        Type:        discordgo.ChatApplicationCommand,
        Options: []*discordgo.ApplicationCommandOption{
            {
                Type:        discordgo.ApplicationCommandOptionString,
                Name:        brbDurationKey,
                Description: "set a duration if you know how long you'll be e.g. 2m, 1h 15m",
                Required:    false,
            },
            {
                Type:        discordgo.ApplicationCommandOptionUser,
                Name:        brbMentionKey,
                Description: "who is brb",
                Required:    false,
            },
        },
    },
    {
        Name:        "back",
        Description: "mark yourself as back, finishing your brb",
        Type:        discordgo.ChatApplicationCommand,
    },
    {
        Name:        "gonefor",
        Description: "find out how long your friend has been gone for",
        Type:        discordgo.ChatApplicationCommand,
        Options: []*discordgo.ApplicationCommandOption{
            {
                Type:        discordgo.ApplicationCommandOptionUser,
                Name:        brbMentionKey,
                Description: "who is brb",
                Required:    false,
            },
        },
    },
}

var CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
    "brb":     BrbChatCommandHandler,
    "back":    BackChatHandler,
    "gonefor": GoneChatHandler,
}

func Register(discord *discordgo.Session, cfg config.Config) map[string]string {
    cmdIds := make(map[string]string, len(ApplicationCommands))

    for _, cmd := range ApplicationCommands {
        rcmd, err := discord.ApplicationCommandCreate(cfg.AppId, cfg.GuildId, &cmd)
        if err != nil {
            log.Error().
                Err(err).
                Str("command", cmd.Name).
                Msg("failed to create slash command")
        }

        cmdIds[rcmd.ID] = rcmd.Name
    }

    return cmdIds
}

func UnRegister(discord *discordgo.Session, cfg config.Config, cmdIds map[string]string) {
    for id, name := range cmdIds {
        err := discord.ApplicationCommandDelete(cfg.AppId, cfg.GuildId, id)
        if err != nil {
            log.Error().
                Err(err).
                Str("command", name).
                Msg("failed to delete slash command")
        }
    }
}
