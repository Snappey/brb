package core

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/rs/zerolog/log"
    "net/http"
)

const (
    minute = 60
    hour   = minute * 60
    day    = hour * 24
)

var (
    bucketDistribution = []float64{
        minute,
        minute * 5,
        minute * 15,
        minute * 30,
        hour,
        hour * 3,
        hour * 6,
        hour * 12,
        day,
    }
)

var (
    registry = prometheus.NewRegistry()

    brbSessionActive = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "brb_session_active",
        Help: "currently active sessions",
    }, []string{"guild_id", "user_id"})

    brbSessionDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "brb_session_finished_duration",
        Help:    "brb duration from start to finish recorded in seconds",
        Buckets: bucketDistribution,
    }, []string{"guild_id", "user_id"})

    brbSessionTargetDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "brb_session_target_duration",
        Help:    "brb target duration set by the reporting user recorded in seconds",
        Buckets: bucketDistribution,
    }, []string{"guild_id", "user_id"})

    brbSessionLateDifferenceDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "brb_session_late_duration",
        Help:    "brb difference between finished duration and target duration >0 recorded in seconds",
        Buckets: bucketDistribution,
    }, []string{"guild_id", "user_id"})

    brbSessionEarlyDifferenceDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "brb_session_early_duration",
        Help:    "brb difference between finished duration and target duration <0 recorded in seconds",
        Buckets: bucketDistribution,
    }, []string{"guild_id", "user_id"})
)

func init() {
    registry.MustRegister(
        brbSessionActive,
        brbSessionDuration,
        brbSessionTargetDuration,
        brbSessionLateDifferenceDuration,
        brbSessionEarlyDifferenceDuration,
    )

    go func() {
        log.Info().Msg("started promhttp at :8080")
        http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))
        if err := http.ListenAndServe(":8080", nil); err != nil {
            log.Fatal().Err(err).Msg("promhttp listener has exited with error")
        } else {
            log.Info().Msg("promhttp has finished")
        }
    }()
}
