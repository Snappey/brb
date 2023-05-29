package manager

import (
    "fmt"
    "github.com/rs/zerolog/log"
    "time"
)

type Manager struct {
    users map[string][]*Brb
}

var manager *Manager

func GetInstance() *Manager {
    if manager == nil {
        manager = &Manager{
            users: map[string][]*Brb{},
        }
    }
    return manager
}

type CreateBrbInput struct {
    TargetUserId    string
    ReportingUserId string
    TargetDuration  time.Duration
}

func (m *Manager) CreateBrb(input CreateBrbInput) error {
    userBrbs, hasUser := m.users[input.TargetUserId]

    if !hasUser {
        m.users[input.TargetUserId] = make([]*Brb, 1)
    } else {
        // Check no active brbs
        latestBrb := userBrbs[len(userBrbs)-1]
        if !latestBrb.Finished {
            return fmt.Errorf("user has active brb, started %s ago", latestBrb.GetDuration().String())
        }
    }

    m.users[input.TargetUserId] = append(
        userBrbs,
        CreateNewBrb(input.ReportingUserId, input.TargetUserId, input.TargetDuration),
    )

    log.Info().
        Str("reporting_user", input.ReportingUserId).
        Str("target_user", input.TargetUserId).
        Dur("target_duration", input.TargetDuration).
        Msg("created brb entry")

    return nil
}

type FinishBrbInput struct {
    TargetUserId    string
    ReportingUserId string
}

func (m *Manager) FinishBrb(input FinishBrbInput) (time.Duration, error) {
    latestBrb, err := m.GetActiveBrb(input.TargetUserId)
    if err != nil {
        return InvalidDuration, err
    }

    latestBrb.Finish()

    log.Info().
        Str("reporting_user", input.ReportingUserId).
        Str("target_user", input.TargetUserId).
        Dur("finished_duration", latestBrb.FinishedDuration).
        Msg("created brb entry")

    return latestBrb.GetDuration(), nil
}

func (m *Manager) GetActiveBrb(targetUserId string) (*Brb, error) {
    userBrbs, hasUser := m.users[targetUserId]

    if !hasUser || len(userBrbs) == 0 {
        return nil, fmt.Errorf("user has no active brb")
    }

    latestBrb := userBrbs[len(userBrbs)-1]
    if latestBrb.Finished {
        return nil, fmt.Errorf("user has no active brb")
    }

    return latestBrb, nil
}
