package handler

import (
	"golang.org/x/net/context"

	"github.com/gophercon/gc18/services/profile-srv/db"
	"github.com/gophercon/gc18/services/profile-srv/proto/record"
)

type Record struct{}

func (s *Record) Create(ctx context.Context, req *record.CreateRequest, rsp *record.CreateResponse) error {
	return db.Create(req.Profile)
}

func (s *Record) Read(ctx context.Context, req *record.ReadRequest, rsp *record.ReadResponse) error {
	profile, err := db.Read(req.Id)
	if err != nil {
		return err
	}
	rsp.Profile = profile
	return nil
}

func (s *Record) Update(ctx context.Context, req *record.UpdateRequest, rsp *record.UpdateResponse) error {
	return db.Update(req.Profile)
}

func (s *Record) Delete(ctx context.Context, req *record.DeleteRequest, rsp *record.DeleteResponse) error {
	return db.Delete(req.Id)
}

func (s *Record) Search(ctx context.Context, req *record.SearchRequest, rsp *record.SearchResponse) error {
	profiles, err := db.Search(req.Name, req.Owner, req.Limit, req.Offset)
	if err != nil {
		return err
	}
	rsp.Profiles = profiles
	return nil
}
