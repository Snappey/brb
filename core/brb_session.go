package core

import (
    "fmt"
    "github.com/google/uuid"
    "github.com/samber/lo"
    "time"
)

const InvalidDuration = -1

type BrbState int

const (
    BrbInitial BrbState = iota
    BrbActive
    BrbFinished
    BrbCancelled
    BrbAwaitingConfirmation
)

type BrbSpan struct {
    StartedAt time.Time
    Active    bool
    Duration  time.Duration
}

type BrbSession struct {
    Id               uuid.UUID
    UserId           string
    ReportingUserId  string
    State            BrbState
    CreatedAt        time.Time
    LastUpdated      time.Time
    TargetDuration   time.Duration
    FinishedDuration time.Duration
    Spans            []BrbSpan
}

func createNewBrb(reportingUserId string, targetUserId string, targetDuration time.Duration) *BrbSession {
    return &BrbSession{
        Id:               uuid.New(),
        UserId:           targetUserId,
        ReportingUserId:  reportingUserId,
        CreatedAt:        time.Now().UTC(),
        LastUpdated:      time.Now().UTC(),
        TargetDuration:   targetDuration,
        State:            BrbInitial,
        FinishedDuration: InvalidDuration,
        Spans:            []BrbSpan{},
    }
}

func createNewBrbAndStart(reportingUserId string, targetUserId string, targetDuration time.Duration) (*BrbSession, error) {
    brb := createNewBrb(reportingUserId, targetUserId, targetDuration)
    return brb, brb.Start()
}

func (b *BrbSession) startSpan() error {
    hasActiveSpan := lo.ContainsBy(b.Spans, func(item BrbSpan) bool {
        return item.Active
    })

    if hasActiveSpan {
        return fmt.Errorf("brb session has active span")
    }

    b.Spans = append(b.Spans, BrbSpan{
        StartedAt: time.Now().UTC(),
        Active:    true,
        Duration:  InvalidDuration,
    })

    return nil
}

func (b *BrbSession) stopSpan() error {
    if len(b.Spans) == 0 {
        return fmt.Errorf("brb session has no recorded spans")
    }

    latestSpan := &b.Spans[len(b.Spans)-1]

    if !latestSpan.Active {
        return nil // If we try to stop a span that's already stopped ignore it
    }

    latestSpan.Active = false
    latestSpan.Duration = time.
        Now().
        UTC().
        Sub(latestSpan.StartedAt)

    return nil
}

func (b *BrbSession) calculateDuration() time.Duration {
    return lo.Reduce(b.Spans, func(agg time.Duration, item BrbSpan, i int) time.Duration {
        if item.Active {
            return agg + (time.Now().UTC().Sub(item.StartedAt))
        }
        return agg + item.Duration
    }, 0)
}

func (b *BrbSession) stop(state BrbState) time.Duration {
    if b.IsFinished() {
        return b.FinishedDuration
    }

    _ = b.stopSpan()

    b.FinishedDuration = b.calculateDuration()
    b.setState(state)

    return b.FinishedDuration
}

func (b *BrbSession) isState(state BrbState) bool {
    return b.State == state
}

func (b *BrbSession) setState(state BrbState) {
    b.State = state
    b.LastUpdated = time.Now().UTC()
}

func (b *BrbSession) IsActive() bool {
    return b.isState(BrbActive)
}

func (b *BrbSession) IsAwaitingConfirmation() bool {
    return b.isState(BrbAwaitingConfirmation)
}

func (b *BrbSession) IsFinished() bool {
    return b.isState(BrbFinished) || b.isState(BrbCancelled)
}

func (b *BrbSession) IsCancelled() bool {
    return b.isState(BrbCancelled)
}

func (b *BrbSession) GetDuration() time.Duration {
    if b.IsFinished() {
        return b.FinishedDuration
    }
    return b.calculateDuration()
}

func (b *BrbSession) Finish() time.Duration {
    return b.stop(BrbFinished)
}

func (b *BrbSession) Cancel() time.Duration {
    return b.stop(BrbCancelled)
}

func (b *BrbSession) Start() error {
    if b.IsActive() {
        return fmt.Errorf("brb session is already active")
    }

    if err := b.startSpan(); err != nil {
        return err
    }

    b.setState(BrbActive)

    return nil
}

func (b *BrbSession) AwaitConfirmation() error {
    if b.IsAwaitingConfirmation() {
        return fmt.Errorf("brb session is already awaiting confirmation")
    }

    if err := b.stopSpan(); err != nil {
        return err
    }

    b.setState(BrbAwaitingConfirmation)

    return nil
}
