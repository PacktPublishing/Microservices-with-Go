package metadata
 
import (
   "context"
   "errors"
   "testing"
 
   "github.com/golang/mock/gomock"
   gen "movieexample.com/gen/mock/metadata/repository"
   "movieexample.com/metadata/internal/repository"   
   "movieexample.com/metadata/pkg/model"
)
 
func TestGetErrNotFound(t *testing.T) {
   ctrl := gomock.NewController(t)
   defer ctrl.Finish()
   repoMock := gen.NewMockmetadataRepository(ctrl)
   c := New(repoMock)
   ctx := context.Background()
   id := "id"
   repoMock.EXPECT().Get(ctx, id).Return(nil, repository.ErrNotFound)
   _, err := c.Get(ctx, id)
   if got, want := err, ErrNotFound; got != want {
       t.Errorf("Get: got %v, want %v", got, want)
   }
}

func TestGetUnexpectedErr(t *testing.T) {
   ctrl := gomock.NewController(t)
   defer ctrl.Finish()
   repoMock := gen.NewMockmetadataRepository(ctrl)
   c := New(repoMock)
   ctx := context.Background()
   id := "id"
   repoErr := errors.New("unexpected error")
   repoMock.EXPECT().Get(ctx, id).Return(nil, repoErr)
   _, err := c.Get(ctx, id)
   if !errors.Is(err, repoErr) {
       t.Errorf("Get: got %v want %v", err, repoErr)
   }
}

func TestGetSuccess(t *testing.T) {
   ctrl := gomock.NewController(t)
   defer ctrl.Finish()
   repoMock := gen.NewMockmetadataRepository(ctrl)
   c := New(repoMock)
   ctx := context.Background()
   id := "id"
   m := &model.Metadata{}
   repoMock.EXPECT().Get(ctx, id).Return(m, nil)
   res, err := c.Get(ctx, id)
   if err != nil {
       t.Errorf("Get: got %v, want %v", err, nil)
   }
   if got, want := res, m; got != want {
       t.Errorf("Get: got %v, want %v", got, want)
   }
}