package manager

import "time"

const InvalidDuration = -1

type Brb struct {
    UserId           string
    ReportingUserId  string
    TargetDuration   time.Duration
    StartedAt        time.Time
    Cancelled        bool
    Finished         bool
    FinishedDuration time.Duration
}

func CreateNewBrb(reportingUserId string, targetUserId string, targetDuration time.Duration) *Brb {
    return &Brb{
        UserId:           targetUserId,
        ReportingUserId:  reportingUserId,
        StartedAt:        time.Now().UTC(),
        TargetDuration:   targetDuration,
        Cancelled:        false,
        Finished:         false,
        FinishedDuration: InvalidDuration,
    }
}

func (b *Brb) GetDuration() time.Duration {
    if b.Finished {
        return b.FinishedDuration
    }
    return time.Now().UTC().Sub(b.StartedAt)
}

func (b *Brb) Finish() {
    if b.Finished {
        return
    }

    b.FinishedDuration = b.GetDuration()
    b.Finished = true
}

func (b *Brb) Cancel() {
    if b.Finished {
        return
    }

    b.Cancelled = true
    b.Finish()
}
