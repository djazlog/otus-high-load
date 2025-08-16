package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"otus-project/internal/config"
	"runtime"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	ctx := context.Background()

	runtime.GOMAXPROCS(6 * runtime.NumCPU())
	err := config.Load(".env")
	if err != nil {
		log.Fatal(fmt.Sprintf("Ошибка при получении env: %v", err))
	}

	dns, err := config.NewPGConfig()
	if err != nil {
		log.Fatal(fmt.Sprintf("Ошибка при получении конфигурации: %v", err))
	}

	// Подключение к БД
	db, err := sql.Open("pgx", dns.DSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	log.Println("starting citus rebalance...")

	if err := StartRebalance(ctx, db); err != nil {
		log.Fatalf("failed to start rebalance: %v", err)
	}

	log.Println("rebalance started, waiting for completion...")

	if err := WaitRebalance(ctx, db, 1*time.Second, true); err != nil {
		log.Fatalf("rebalance failed: %v", err)
	}

	log.Println("rebalance finished successfully ✅")
}

type RebalanceStatus struct {
	JobID       int64           `json:"job_id"`
	State       string          `json:"state"`
	JobType     string          `json:"job_type"`
	Description string          `json:"description"`
	StartedAt   *time.Time      `json:"started_at"`
	FinishedAt  *time.Time      `json:"finished_at"`
	Details     json.RawMessage `json:"details"`
}

// структура для details
type RebalanceDetails struct {
	GroupID     int    `json:"group_id"`
	MovedShards int    `json:"moved_shards"`
	TotalShards int    `json:"total_shards"`
	Phase       string `json:"phase"`
}

func StartRebalance(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `SELECT citus_rebalance_start()`)
	return err
}

func GetRebalanceStatus(ctx context.Context, db *sql.DB) (*RebalanceStatus, error) {
	row := db.QueryRowContext(ctx, `SELECT job_id, state, job_type, description, started_at, finished_at, details
	                                FROM citus_rebalance_status()`)
	var s RebalanceStatus
	if err := row.Scan(&s.JobID, &s.State, &s.JobType, &s.Description, &s.StartedAt, &s.FinishedAt, &s.Details); err != nil {
		return nil, err
	}
	return &s, nil
}

func WaitRebalance(ctx context.Context, db *sql.DB, pollEvery time.Duration, logProgress bool) error {
	t := time.NewTicker(pollEvery)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			st, err := GetRebalanceStatus(ctx, db)
			if err != nil {
				_, waitErr := db.ExecContext(ctx, `SELECT citus_rebalance_wait()`)
				return waitErr
			}

			if logProgress {
				// парсим details
				var det RebalanceDetails
				if len(st.Details) > 0 {
					if err := json.Unmarshal(st.Details, &det); err == nil && det.TotalShards > 0 {
						percent := float64(det.MovedShards) / float64(det.TotalShards) * 100
						log.Printf("[citus] job=%d state=%s phase=%s %d/%d shards (%.1f%%)",
							st.JobID, st.State, det.Phase, det.MovedShards, det.TotalShards, percent)
					} else {
						log.Printf("[citus] job=%d state=%s desc=%s details=%s",
							st.JobID, st.State, st.Description, string(st.Details))
					}
				} else {
					log.Printf("[citus] job=%d state=%s desc=%s", st.JobID, st.State, st.Description)
				}
			}

			switch st.State {
			case "finished":
				return nil
			case "failed", "canceled":
				return errors.New("citus rebalance failed or canceled")
			default:
				// running|waiting — продолжаем
			}
		}
	}
}
