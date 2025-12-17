package urlsigner

// Signer provides URL signing functionality
type Signer struct {
	Secret []byte
}

// GenerateTokenFromString generates a signed token from a string
func (s *Signer) GenerateTokenFromString(data string) string {
	// Implementation will be added
	return ""
}

// VerifyToken verifies if a token is valid
func (s *Signer) VerifyToken(token string) bool {
	// Implementation will be added
	return false
}

// Expired checks if a token has expired
func (s *Signer) Expired(token string, minutesUntilExpire int) bool {
	// Implementation will be added
	return false
}

