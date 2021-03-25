package externals

import (
	"context"
	"firebase.google.com/go/auth"
	"github.com/stretchr/testify/mock"
)

type FbAuth struct {
	mock.Mock
}

func (f *FbAuth) VerifyIDTokenAndCheckRevoked(ctx context.Context, idToken string) (*auth.Token, error) {
	args := f.Called(ctx, idToken)
	return args.Get(0).(*auth.Token), args.Error(1)
}

func (f *FbAuth) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	args := f.Called(ctx, uid)
	return args.Get(0).(*auth.UserRecord), args.Error(1)
}

func NewMockFirebaseAuth() *FbAuth  {
	return &FbAuth{}
}
