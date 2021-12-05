package checkpoint

import (
	"httpInterceptor/config"
	"time"
)

func Monitor() {
	period := config.GetSnapshotPeriodicity()
	periodicity := time.Duration(period) * time.Millisecond
	for _ = range time.Tick(periodicity) {
		generateSnapshot()
	}
}

func generateSnapshot() {

}
