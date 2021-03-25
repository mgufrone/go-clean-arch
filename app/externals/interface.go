package externals

import (
	"context"
	"firebase.google.com/go/auth"
)

type IFirebaseAuth interface {
	VerifyIDTokenAndCheckRevoked(ctx context.Context, idToken string) (*auth.Token, error)
	GetUser(ctx context.Context, uid string) (*auth.UserRecord, error)
}
