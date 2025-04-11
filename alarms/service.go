// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package alarms

import (
	"context"
	"time"

	"github.com/absmach/supermq"
	"github.com/absmach/supermq/pkg/authn"
)

type service struct {
	idp  supermq.IDProvider
	repo Repository
}

var _ Service = (*service)(nil)

func NewService(idp supermq.IDProvider, repo Repository) Service {
	return &service{
		idp:  idp,
		repo: repo,
	}
}

func (s *service) CreateAlarm(ctx context.Context, alarm Alarm) error {
	id, err := s.idp.ID()
	if err != nil {
		return err
	}
	alarm.ID = id
	alarm.CreatedAt = time.Now()

	if err := alarm.Validate(); err != nil {
		return err
	}

	lastAlarms, err := s.repo.ListAlarms(ctx, PageMetadata{Limit: 1, Offset: 0})
	if err != nil {
		return err
	}

	if len(lastAlarms.Alarms) > 0 {
		lastAlarm := lastAlarms.Alarms[0]
		if lastAlarm.Severity == alarm.Severity {
			return nil
		}
	}

	_, err = s.repo.CreateAlarm(ctx, alarm)

	return err
}

func (s *service) ViewAlarm(ctx context.Context, session authn.Session, alarmID string) (Alarm, error) {
	return s.repo.ViewAlarm(ctx, alarmID, session.DomainID)
}

func (s *service) ListAlarms(ctx context.Context, session authn.Session, pm PageMetadata) (AlarmsPage, error) {
	return s.repo.ListAlarms(ctx, pm)
}

func (s *service) DeleteAlarm(ctx context.Context, session authn.Session, alarmID string) error {
	return s.repo.DeleteAlarm(ctx, alarmID)
}

func (s *service) UpdateAlarm(ctx context.Context, session authn.Session, alarm Alarm) (Alarm, error) {
	alarm.UpdatedAt = time.Now()
	alarm.UpdatedBy = session.UserID

	if err := alarm.Validate(); err != nil {
		return Alarm{}, err
	}

	return s.repo.UpdateAlarm(ctx, alarm)
}
