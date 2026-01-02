package utils

import (
	"context"
	"errors"

	"google.golang.org/api/idtoken"
)

type GoogleUserClaims struct {
	Sub     string // Google user ID (unique identifier)
	Email   string
	Name    string
	Picture string // URL to the profile picture
}

// NOTE: Replace this with the actual CLIENT ID of your Mobile Application
// (e.g., the Android/iOS client ID from your Google Cloud Console OAuth credentials).

const MobileClientID = "YOUR_MOBILE_APP_CLIENT_ID_HERE"




// VerifyGoogleIDToken securely validates the JWT ID Token and extracts user claims.
func VerifyGoogleIDToken(ctx context.Context, idToken string) (*GoogleUserClaims, error) {

	payload, err := idtoken.Validate(ctx, idToken, MobileClientID)
	if err != nil {
		return nil, errors.New("google ID token verification failed: " + err.Error())
	}

	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		return nil, errors.New("email claim is missing or invalid")
	}

	name, ok := payload.Claims["name"].(string)
	if !ok || name == "" {
		return nil, errors.New("name claim is missing or invalid")
	}

	// 5. Success: Map the verified claims to your local struct
	claims := &GoogleUserClaims{
		Sub:     payload.Subject, // User ID is exposed as UserId in Tokeninfo response
		Email:   email,
		Name:    name,
		Picture: payload.Claims["picture"].(string),
	}

	return claims, nil
}

// NOTE: The GoogleOAuth struct and its methods (NewGoogleOAuth, AuthURL, Exchange)
// have been REMOVED as they are no longer needed for the mobile ID Token flow.
