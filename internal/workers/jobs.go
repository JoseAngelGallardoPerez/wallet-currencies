package workers

import (
	"sync"

	"github.com/inconshreveable/log15"
	"github.com/jasonlvhit/gocron"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	workerServices "github.com/Confialink/wallet-currencies/internal/services/workers"
)

const atDynamicCetTime = "5:00PM"

type Jobs struct {
	db           *gorm.DB
	logger       log15.Logger
	scheduler    *gocron.Scheduler
	ratesUpdater *UpdateRates
	mutex        sync.Mutex
}

func NewJobs(db *gorm.DB, scheduler *gocron.Scheduler, ratesUpdater *UpdateRates, logger log15.Logger) *Jobs {
	return &Jobs{db, logger.New("service", "Jobs"), scheduler, ratesUpdater, sync.Mutex{}}
}

// Start initializes jobs and starts scheduler
func (s *Jobs) Start() error {
	if err := s.scheduler.Every(1).Day().At(s.localTime()).Do(s.updateRates); err != nil {
		return errors.Wrap(err, "cannot schedule updating rates")
	}

	s.scheduler.Start()
	s.logger.Info("scheduler is started")

	return nil
}

func (s *Jobs) updateRates() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.logger.Info("job 'updateRates' is started")
	tx := s.db.Begin()
	ratesUpdater := s.ratesUpdater.WrapContext(tx)
	if err := ratesUpdater.Update(); err != nil {
		s.logger.Error("cannot update rates", "err", err)
		tx.Rollback()
	} else {
		tx.Commit()
	}

	s.logger.Info("job 'updateRates' is finished")
}

func (s *Jobs) localTime() string {
	result, _ := workerServices.DynamicCetToLocal(atDynamicCetTime)
	return result
}
