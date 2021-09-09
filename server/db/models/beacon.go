package models

/*
	Sliver Implant Framework
	Copyright (C) 2021  Bishop Fox

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"time"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
	"github.com/gofrs/uuid"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

// Beacon - Represents a host machine
type Beacon struct {
	gorm.Model

	CreatedAt time.Time `gorm:"->;<-:create;"`

	ID                uuid.UUID `gorm:"type:uuid;"`
	Name              string
	Hostname          string
	UUID              uuid.UUID `gorm:"type:uuid;"` // Host UUID
	Username          string
	UID               string
	GID               string
	OS                string
	Arch              string
	Transport         string
	RemoteAddress     string
	PID               int32
	Filename          string
	LastCheckin       time.Time
	Version           string
	ReconnectInterval int64
	ProxyURL          string
	PollTimeout       int64

	ImplantBuildID uuid.UUID `gorm:"type:uuid;"`

	Interval    int64
	Jitter      int64
	NextCheckin int64

	Tasks []BeaconTask
}

// BeforeCreate - GORM hook
func (b *Beacon) BeforeCreate(tx *gorm.DB) (err error) {
	b.CreatedAt = time.Now()
	return nil
}

func (b *Beacon) ToProtobuf() *clientpb.Beacon {
	return &clientpb.Beacon{
		ID:            b.ID.String(),
		Name:          b.Name,
		Hostname:      b.Hostname,
		UUID:          b.UUID.String(),
		Username:      b.Username,
		UID:           b.UID,
		GID:           b.GID,
		OS:            b.OS,
		Arch:          b.Arch,
		Transport:     b.Transport,
		RemoteAddress: b.RemoteAddress,
		PID:           b.PID,
		Filename:      b.Filename,
		LastCheckin:   b.LastCheckin.Unix(),
		Version:       b.Version,
		// ReconnectInterval: b.ReconnectInterval,
		ProxyURL: b.ProxyURL,
		// PollTimeout: b.PollTimeout,
		// ImplantBuildID: b.ImplantBuildID.String(),
		Interval:    b.Interval,
		Jitter:      b.Jitter,
		NextCheckin: b.NextCheckin,
	}
}

func (b *Beacon) Task(envelope *sliverpb.Envelope) (*BeaconTask, error) {
	data, err := proto.Marshal(envelope)
	if err != nil {
		return nil, err
	}
	task := &BeaconTask{
		BeaconID: b.ID,
		State:    PENDING,
		Request:  data,
	}
	return task, nil
}

// BeaconTask - Represents a host machine
const (
	PENDING   = "pending"
	SENT      = "sent"
	COMPLETED = "completed"
)

type BeaconTask struct {
	gorm.Model

	ID          uuid.UUID `gorm:"primaryKey;->;<-:create;type:uuid;"`
	BeaconID    uuid.UUID `gorm:"type:uuid;"`
	CreatedAt   time.Time `gorm:"->;<-:create;"`
	State       string
	SentAt      time.Time
	CompletedAt time.Time
	Request     []byte // *sliverpb.Envelope
	Response    []byte // *sliverpb.Envelope
}

// BeforeCreate - GORM hook
func (b *BeaconTask) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}
	b.CreatedAt = time.Now()
	b.State = PENDING
	return nil
}