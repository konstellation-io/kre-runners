package kre

import (
	"encoding/json"
	"fmt"
	nats2 "github.com/nats-io/nats.go"
	"time"

	"github.com/konstellation-io/kre-runners/kre-go/config"
	"github.com/konstellation-io/kre/libs/simplelogger"
)

type SaveMetricErr string

const (
	classificationMetricsColl               = "classificationMetrics"
	ErrMissingValues          SaveMetricErr = "missing_values"
	ErrNewLabels              SaveMetricErr = "new_labels"
)

func (s SaveMetricErr) IsValid() error {
	switch s {
	case ErrMissingValues, ErrNewLabels:
		return nil
	}
	return fmt.Errorf("invalid SaveMetricErr type: %s", s)
}

type SaveMetricMsg struct {
	Coll string           `json:"coll"`
	Doc  SaveMetricMsgDoc `json:"doc"`
}

type SaveMetricMsgDoc struct {
	Date           string `json:"date"`
	Error          string `json:"error"`
	PredictedValue string `json:"predictedValue"`
	TrueValue      string `json:"trueValue"`
	RuntimeID      string `json:"runtimeId"`
	VersionID      string `json:"versionId"`
	VersionName    string `json:"versionName"`
}

type contextPrediction struct {
	cfg    config.Config
	nc     *nats2.Conn
	logger *simplelogger.SimpleLogger
}

func (c *contextPrediction) Save(date time.Time, predictedValue, trueValue string) {
	msg, err := json.Marshal(SaveMetricMsg{
		Coll: classificationMetricsColl,
		Doc: SaveMetricMsgDoc{
			Date:           date.Format(time.RFC3339),
			PredictedValue: predictedValue,
			TrueValue:      trueValue,
			RuntimeID:      c.cfg.RuntimeID,
			VersionID:      c.cfg.VersionID,
			VersionName:    c.cfg.Version,
		},
	})

	if err != nil {
		c.logger.Infof("Error marshalling SaveMetricMsgDoc: %s", err)
		return
	}

	_, err = c.nc.Request(c.cfg.NATS.MongoWriterSubject, msg, saveMetricTimeout)
	if err != nil {
		c.logger.Infof("Error sending metric to NATS: %s", err)
	}
}

func (c *contextPrediction) SaveError(saveMetricErr SaveMetricErr) {
	if err := saveMetricErr.IsValid(); err != nil {
		c.logger.Error(err.Error())
		return
	}

	msg, err := json.Marshal(SaveMetricMsg{
		Coll: classificationMetricsColl,
		Doc: SaveMetricMsgDoc{
			Date:        time.Now().Format(time.RFC3339),
			RuntimeID:   c.cfg.RuntimeID,
			VersionName: c.cfg.Version,
			Error:       string(saveMetricErr),
		},
	})
	if err != nil {
		c.logger.Infof("Error generating SaveMetricMsg JSON: %s", err)
	}

	_, err = c.nc.Request(c.cfg.NATS.MongoWriterSubject, msg, saveMetricTimeout)
	if err != nil {
		c.logger.Infof("Error sending error metric to NATS: %s", err)
	}
}
