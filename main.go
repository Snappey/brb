package main

import (
    "brb/commands"
    "brb/config"
    "fmt"
    "github.com/bwmarrin/discordgo"
    "github.com/rs/zerolog/log"
    "os"
    "os/signal"
)

func main() {
    cfg := config.Get()

    discord, err := discordgo.New(fmt.Sprintf("Bot %s", cfg.AuthToken))
    if err != nil {
        log.Fatal().
            Err(err).
            Msg("failed to create bot")
    }

    discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
        log.Info().Msg("Bot is ready!")
    })

    discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
        if handler, ok := commands.CommandHandlers[i.ApplicationCommandData().Name]; ok {
            handler(s, i)
        }
    })

    discord.AddHandler(func(s *discordgo.Session, i *discordgo.MessageCreate) {
        if commands.MentionedBot(s, i) {
            commands.MentionHandler(s, i)
        }
    })

    cmdIds := commands.Register(discord, cfg)

    err = discord.Open()
    if err != nil {
        log.Fatal().
            Err(err).
            Msg("failed to open session")
    }
    defer discord.Close()

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt)
    <-stop
    log.Info().Msg("interrupt received.. shutting down")

    commands.UnRegister(discord, cfg, cmdIds)
    log.Info().Msg("unregistered commands")
}
