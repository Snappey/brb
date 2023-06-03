package core

import (
    "fmt"
    "github.com/rs/zerolog/log"
    "time"
)

type CreateBrbInput struct {
    GuildId         string
    TargetUserId    string
    ReportingUserId string
    TargetDuration  time.Duration
}

type Manager struct {
    users map[string][]*BrbSession
}

var manager *Manager

func GetManagerInstance() *Manager {
    if manager == nil {
        manager = &Manager{
            users: map[string][]*BrbSession{},
        }
    }
    return manager
}

func (m *Manager) CreateBrb(input CreateBrbInput) error {
    userBrbs, hasUser := m.users[input.TargetUserId]

    if !hasUser {
        m.users[input.TargetUserId] = make([]*BrbSession, 1)
    } else {
        // Check no active brbs
        latestBrb := userBrbs[len(userBrbs)-1]
        if !latestBrb.IsFinished() {
            return fmt.Errorf("user has active brbSession, started %s ago", latestBrb.GetDuration().String())
        }
    }

    brbSession, err := createNewBrbAndStart(input.GuildId, input.ReportingUserId, input.TargetUserId, input.TargetDuration)
    if err != nil {
        return err
    }

    m.users[input.TargetUserId] = append(
        userBrbs,
        brbSession,
    )

    log.Info().
        Str("reporting_user", input.ReportingUserId).
        Str("target_user", input.TargetUserId).
        Dur("target_duration", input.TargetDuration).
        Msg("created brbSession entry")

    return nil
}

func (m *Manager) GetActiveBrb(targetUserId string) (*BrbSession, error) {
    userBrbs, hasUser := m.users[targetUserId]

    if !hasUser || len(userBrbs) == 0 {
        return nil, fmt.Errorf("user has no active brb")
    }

    latestBrb := userBrbs[len(userBrbs)-1]
    if latestBrb.IsFinished() {
        return nil, fmt.Errorf("user has no active brb")
    }

    return latestBrb, nil
}
