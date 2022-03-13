package runtime

import (
	"os"
	"time"

	"github.com/codingsince1985/geo-golang/openstreetmap"
	"github.com/johnjones4/golden-hour-bot/lib/engine"
	"github.com/johnjones4/golden-hour-bot/lib/service"
	"github.com/johnjones4/golden-hour-bot/lib/shared"
	"github.com/johnjones4/golden-hour-bot/lib/telegram"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

type inMemoryState struct {
	state string
	info  interface{}
}
type inMemoryStateEngine struct {
	states map[int]inMemoryState
}

func (e *inMemoryStateEngine) GetChatState(id int) (string, interface{}, error) {
	if state, ok := e.states[id]; ok {
		return state.state, state.info, nil
	}
	return shared.DefaultState, nil, nil
}

func (e *inMemoryStateEngine) SetChatState(id int, state string, info interface{}) error {
	e.states[id] = inMemoryState{state, info}
	return nil
}

type inMemoryReminderStorage struct {
	data    map[string][]shared.Reminder
	regions map[string]shared.Region
}

func (rs *inMemoryReminderStorage) SaveReminder(r shared.Reminder) error {
	key := r.GetRegionKey()

	var reminders []shared.Reminder
	if reminders1, ok := rs.data[key]; ok {
		reminders = reminders1
	} else {
		reminders = make([]shared.Reminder, 0)
	}

	for _, reminder1 := range reminders {
		if reminder1.ChatId == r.ChatId {
			return shared.ErrorDuplicateReminder(key, r.ChatId)

		}
	}

	if _, ok := rs.regions[r.GetRegionKey()]; !ok {
		c := r.GetRegion()
		rs.regions[r.GetRegionKey()] = shared.Region{
			Region:   r.GetRegionKey(),
			Location: c,
		}
	}

	rs.data[key] = append(reminders, r)

	return nil
}

func (rs *inMemoryReminderStorage) GetRegions() ([]shared.Region, error) {
	regions := make([]shared.Region, 0)
	for _, region := range rs.regions {
		regions = append(regions, region)
	}
	return regions, nil
}

func (rs *inMemoryReminderStorage) UpdateRegion(r shared.Region) error {
	rs.regions[r.Region] = r
	return nil
}

func (rs *inMemoryReminderStorage) GetRemindersInRegion(region string) ([]shared.Reminder, error) {
	if reminders, ok := rs.data[region]; ok {
		return reminders, nil
	}
	return []shared.Reminder{}, nil
}

type inMemoryAlertQueue struct {
	Client telegram.Telegram
}

func (q *inMemoryAlertQueue) EnqueueAlerts(predType string, reminders []shared.Reminder) error {
	for _, reminder := range reminders {
		err := service.SendAlert(q.Client, predType, reminder)
		if err != nil {
			return err
		}
	}
	return nil
}

var rs = &inMemoryReminderStorage{
	data:    make(map[string][]shared.Reminder),
	regions: make(map[string]shared.Region),
}

func StartInMemoryRuntime() {
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	geocoder := openstreetmap.Geocoder()

	mq := &service.DirectQueue{
		Client: telegram.Telegram{
			Token: os.Getenv("TELEGRAM_TOKEN"),
		},
		PredictionParser: service.PredictionRequestParser{
			DateParser: w,
			Geocoder:   geocoder,
		},
		ReminderStorage: rs,
		Geocoder:        geocoder,
	}

	e := engine.Engine{
		StateEngine: &inMemoryStateEngine{
			states: make(map[int]inMemoryState),
		},
		Queue: mq,
	}

	runLocalServer(&e)
}

func StartInMemoryAlerter() {
	aq := inMemoryAlertQueue{
		Client: telegram.Telegram{
			Token: os.Getenv("TELEGRAM_TOKEN"),
		},
	}
	for {
		err := service.RunAlertCycle(rs, &aq)
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Second * 30)
	}
}
