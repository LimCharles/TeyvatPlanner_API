package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.42

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"teyvat_planner_api/auth"
	"teyvat_planner_api/graph/model"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, username string, email string, password string) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var user model.User

	err = r.DB.QueryRow(
		ctx,
		"INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, username, email, password",
		username, email, string(hashedPassword),
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateCommission is the resolver for the createCommission field.
func (r *mutationResolver) CreateCommission(ctx context.Context, name string, category string) (*model.Commission, error) {
	authUser := auth.ForContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("Access Denied")
	}

	var commission model.Commission

	err := r.DB.QueryRow(
		ctx,
		"INSERT INTO commissions (name, category, userID) VALUES ($1, $2, $3) RETURNING id, name, category, completed",
		name, category, authUser.ID,
	).Scan(&commission.ID, &commission.Name, &commission.Category, &commission.Completed)
	if err != nil {
		return nil, err
	}

	commission.User = &model.User{
		ID:       authUser.ID,
		Username: authUser.Username,
		Email:    authUser.Email,
	}

	return &commission, nil
}

// CreateDomain is the resolver for the createDomain field.
func (r *mutationResolver) CreateDomain(ctx context.Context, name string) (*model.Domain, error) {
	authUser := auth.ForContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("Access Denied")
	}

	var domain model.Domain

	err := r.DB.QueryRow(
		ctx,
		"INSERT INTO domains (name, userID) VALUES ($1, $2) RETURNING id, name, completed",
		name, authUser.ID,
	).Scan(&domain.ID, &domain.Name, &domain.Completed)
	if err != nil {
		return nil, err
	}

	domain.User = &model.User{
		ID:       authUser.ID,
		Username: authUser.Username,
		Email:    authUser.Email,
	}

	return &domain, nil
}

// CreateWeeklyBoss is the resolver for the createWeeklyBoss field.
func (r *mutationResolver) CreateWeeklyBoss(ctx context.Context, name string) (*model.WeeklyBoss, error) {
	authUser := auth.ForContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("Access Denied")
	}

	var weeklyBoss model.WeeklyBoss

	err := r.DB.QueryRow(
		ctx,
		"INSERT INTO weeklybosses (name, userID) VALUES ($1, $2) RETURNING id, name, completed",
		name, authUser.ID,
	).Scan(&weeklyBoss.ID, &weeklyBoss.Name, &weeklyBoss.Completed)
	if err != nil {
		return nil, err
	}

	weeklyBoss.User = &model.User{
		ID:       authUser.ID,
		Username: authUser.Username,
		Email:    authUser.Email,
	}

	return &weeklyBoss, nil
}

// CreateRandomQuest is the resolver for the createRandomQuest field.
func (r *mutationResolver) CreateRandomQuest(ctx context.Context, name string, longitude *float64, latitude *float64) (*model.RandomQuest, error) {
	authUser := auth.ForContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("Access Denied")
	}

	var randomQuest model.RandomQuest

	err := r.DB.QueryRow(
		ctx,
		"INSERT INTO randomquests (name, longitude, latitude, userID) VALUES ($1, $2, $3, $4) RETURNING id, name, longitude, latitude, completed",
		name, longitude, latitude, authUser.ID,
	).Scan(&randomQuest.ID, &randomQuest.Name, &randomQuest.Longitude, &randomQuest.Latitude, &randomQuest.Completed)
	if err != nil {
		return nil, err
	}

	randomQuest.User = &model.User{
		ID:       authUser.ID,
		Username: authUser.Username,
		Email:    authUser.Email,
	}

	return &randomQuest, nil
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, username *string, email *string, password *string) (*model.User, error) {
	authUser := auth.ForContext(ctx)
	if authUser == nil || authUser.ID != id {
		return nil, fmt.Errorf("Access Denied")
	}

	var user model.User

	query := `
		UPDATE users SET 
			username = COALESCE($1, username), 
			email = COALESCE($2, email), 
			password = COALESCE($3, password) 
		WHERE id = $4
		AND (
			($1 IS NOT NULL AND $1 IS DISTINCT FROM username) OR
			($2 IS NOT NULL AND $2 IS DISTINCT FROM email) OR
			($3 IS NOT NULL AND $3 IS DISTINCT FROM password)
		)
		RETURNING id, username, email, password
	`

	err := r.DB.QueryRow(
		ctx,
		query,
		username, email, password, id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateCommission is the resolver for the updateCommission field.
func (r *mutationResolver) UpdateCommission(ctx context.Context, id string, name *string, category *string, completed *bool) (*model.Commission, error) {
	authUser := auth.ForContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("Access Denied")
	}

	var commission model.Commission

	query := `
		UPDATE commissions SET
			name = COALESCE($1, name),
			category = COALESCE($2, category),
			completed = COALESCE($3, completed)
		WHERE id = $4
		AND userID = $5
		AND (
			($1 IS NOT NULL AND $1 IS DISTINCT FROM name) OR
			($2 IS NOT NULL AND $2 IS DISTINCT FROM category) OR
			($3 IS NOT NULL AND $3 IS DISTINCT FROM completed)
		)
		RETURNING id, name, category, completed
	`

	err := r.DB.QueryRow(
		ctx,
		query,
		name, category, completed, id, authUser.ID,
	).Scan(&commission.ID, &commission.Name, &commission.Category, &commission.Completed)

	if err != nil {
		return nil, err
	}
	commission.User = &model.User{
		ID:       authUser.ID,
		Username: authUser.Username,
		Email:    authUser.Email,
	}

	return &commission, nil
}

// UpdateDomain is the resolver for the updateDomain field.
func (r *mutationResolver) UpdateDomain(ctx context.Context, id string, name *string, completed *bool) (*model.Domain, error) {
	authUser := auth.ForContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("Access Denied")
	}

	var domain model.Domain

	query := `
		UPDATE domains SET
			name = COALESCE($1, name),
			completed = COALESCE($2, completed)
		WHERE id = $3
		AND userID = $4
		AND (
			($1 IS NOT NULL AND $1 IS DISTINCT FROM name) OR
			($2 IS NOT NULL AND $2 IS DISTINCT FROM completed)
		)
		RETURNING id, name, completed
	`

	err := r.DB.QueryRow(
		ctx,
		query,
		name, completed, id, authUser.ID,
	).Scan(&domain.ID, &domain.Name, &domain.Completed)
	if err != nil {
		return nil, err
	}

	domain.User = &model.User{
		ID:       authUser.ID,
		Username: authUser.Username,
		Email:    authUser.Email,
	}

	return &domain, nil
}

// UpdateWeeklyBoss is the resolver for the updateWeeklyBoss field.
func (r *mutationResolver) UpdateWeeklyBoss(ctx context.Context, id string, name *string, completed *bool) (*model.WeeklyBoss, error) {
	authUser := auth.ForContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("Access Denied")
	}

	var weeklyBoss model.WeeklyBoss

	query := `
		UPDATE weeklybosses SET
			name = COALESCE($1, name),
			completed = COALESCE($2, completed)
		WHERE id = $3
		AND userID = $4
		AND (
			($1 IS NOT NULL AND $1 IS DISTINCT FROM name) OR
			($2 IS NOT NULL AND $2 IS DISTINCT FROM completed)
		)
		RETURNING id, name, completed
	`

	err := r.DB.QueryRow(
		ctx,
		query,
		name, completed, id, authUser.ID,
	).Scan(&weeklyBoss.ID, &weeklyBoss.Name, &weeklyBoss.Completed)
	if err != nil {
		return nil, err
	}

	weeklyBoss.User = &model.User{
		ID:       authUser.ID,
		Username: authUser.Username,
		Email:    authUser.Email,
	}

	return &weeklyBoss, nil
}

// UpdateRandomQuest is the resolver for the updateRandomQuest field.
func (r *mutationResolver) UpdateRandomQuest(ctx context.Context, id string, name *string, longitude *float64, latitude *float64, completed *bool) (*model.RandomQuest, error) {
	authUser := auth.ForContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("Access Denied")
	}

	var randomQuest model.RandomQuest

	query := `
		UPDATE randomquests SET
			name = COALESCE($1, name),
			longitude = COALESCE($2, longitude),
			latitude = COALESCE($3, latitude),
			completed = COALESCE($4, completed)
		WHERE id = $5
		AND userID = $6
		AND (
			($1 IS NOT NULL AND $1 IS DISTINCT FROM name) OR
			($2 IS NOT NULL AND $2 IS DISTINCT FROM longitude) OR
			($3 IS NOT NULL AND $3 IS DISTINCT FROM latitude) OR
			($4 IS NOT NULL AND $4 IS DISTINCT FROM completed)
		)
		RETURNING id, name, longitude, latitude, completed
	`

	err := r.DB.QueryRow(
		ctx,
		query,
		name, longitude, latitude, completed, id, authUser.ID,
	).Scan(&randomQuest.ID, &randomQuest.Name, &randomQuest.Longitude, &randomQuest.Latitude, &randomQuest.Completed)
	if err != nil {
		return nil, err
	}

	randomQuest.User = &model.User{
		ID:       authUser.ID,
		Username: authUser.Username,
		Email:    authUser.Email,
	}

	return &randomQuest, nil
}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (string, error) {
	authUser := auth.ForContext(ctx)
	if authUser == nil || authUser.ID != id {
		return "", fmt.Errorf("Access Denied")
	}

	cmdTag, err := r.DB.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return "", err
	}

	if cmdTag.RowsAffected() == 0 {
		return "", fmt.Errorf("no user found with ID %s", id)
	}

	return id, nil
}

// DeleteCommission is the resolver for the deleteCommission field.
func (r *mutationResolver) DeleteCommission(ctx context.Context, id string) (string, error) {
	authUser := auth.ForContext(ctx)
	if authUser == nil {
		return "", fmt.Errorf("Access Denied")
	}

	cmdTag, err := r.DB.Exec(ctx, `DELETE FROM commissions WHERE id = $1 AND userID = $2`, id, authUser.ID)
	if err != nil {
		return "", err
	}

	if cmdTag.RowsAffected() == 0 {
		return "", fmt.Errorf("No commission found with ID %s or access denied", id)
	}

	return id, nil
}

// DeleteDomain is the resolver for the deleteDomain field.
func (r *mutationResolver) DeleteDomain(ctx context.Context, id string) (string, error) {
	authUser := auth.ForContext(ctx)
	if authUser == nil {
		return "", fmt.Errorf("Access Denied")
	}

	cmdTag, err := r.DB.Exec(ctx, `DELETE FROM domains WHERE id = $1 AND userID = $2`, id, authUser.ID)
	if err != nil {
		return "", err
	}

	if cmdTag.RowsAffected() == 0 {
		return "", fmt.Errorf("No domain found with ID %s or access denied", id)
	}

	return id, nil
}

// DeleteWeeklyBoss is the resolver for the deleteWeeklyBoss field.
func (r *mutationResolver) DeleteWeeklyBoss(ctx context.Context, id string) (string, error) {
	authUser := auth.ForContext(ctx)
	if authUser == nil {
		return "", fmt.Errorf("Access Denied")
	}

	cmdTag, err := r.DB.Exec(ctx, `DELETE FROM weeklybosses WHERE id = $1 AND userID = $2`, id, authUser.ID)
	if err != nil {
		return "", err
	}

	if cmdTag.RowsAffected() == 0 {
		return "", fmt.Errorf("No weekly boss found with ID %s or access denied", id)
	}

	return id, nil
}

// DeleteRandomQuest is the resolver for the deleteRandomQuest field.
func (r *mutationResolver) DeleteRandomQuest(ctx context.Context, id string) (string, error) {
	authUser := auth.ForContext(ctx)
	if authUser == nil {
		return "", fmt.Errorf("Access Denied")
	}

	cmdTag, err := r.DB.Exec(ctx, `DELETE FROM randomquests WHERE id = $1 AND userID = $2`, id, authUser.ID)
	if err != nil {
		return "", err
	}

	if cmdTag.RowsAffected() == 0 {
		return "", fmt.Errorf("No random quest found with ID %s or access denied", id)
	}

	return id, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, email string, password string) (string, error) {
	var userID string
	var hashedPassword string

	err := r.DB.QueryRow(
		ctx,
		"SELECT id, password FROM users WHERE email = $1",
		email,
	).Scan(&userID, &hashedPassword)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", err
	}

	claims := &jwt.StandardClaims{
		Issuer:  "TeyvatPlanner",
		Subject: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshToken, err := token.SignedString([]byte(os.Getenv("TP_REFRESH_TOKEN_SECRET")))
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

// RequestAccessToken is the resolver for the requestAccessToken field.
func (r *mutationResolver) RequestAccessToken(ctx context.Context, refreshToken string) (string, error) {
	parsedRefreshToken, err := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TP_REFRESH_TOKEN_SECRET")), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := parsedRefreshToken.Claims.(*jwt.StandardClaims)
	if !ok || !parsedRefreshToken.Valid {
		return "", fmt.Errorf("Invalid Refresh Token")
	}

	accessTokenClaims := &jwt.StandardClaims{
		Issuer:    claims.Issuer,
		Subject:   claims.Subject,
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	signedAccessToken, err := accessToken.SignedString([]byte(os.Getenv("TP_ACCESS_TOKEN_SECRET")))
	if err != nil {
		return "", err
	}

	return signedAccessToken, nil
}

// GoogleLogin is the resolver for the googleLogin field.
func (r *mutationResolver) GoogleLogin(ctx context.Context, token string) (string, error) {
	conf := &oauth2.Config{
        ClientID: os.Getenv("GOOGLE_CLIENT_ID"),
        ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        Endpoint: oauth2.Endpoint{
            AuthURL: "https://accounts.google.com/o/oauth2/auth",
            TokenURL: "https://oauth2.googleapis.com/token",
        },
        RedirectURL: "YOUR_REDIRECT_URL", 
    }

	oauth2Token := &oauth2.Token{
        AccessToken: token,
    }

	client := conf.Client(ctx, oauth2Token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
    if err != nil {
        return "", fmt.Errorf("Error fetching user information")
    }
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("Error reading response")
    }
	
	var userInfo struct {
        Email string `json:"email"`
    }
	
    err = json.Unmarshal(data, &userInfo) 
	if err != nil {
        return "", fmt.Errorf("Error parsing user information")
    }

    if userInfo.Email == "" {
        return "", fmt.Errorf("Email not found in user info")
    }
	
	var userID string
	err = r.DB.QueryRow(
		ctx,
		"SELECT id FROM users WHERE email = $1",
		userInfo.Email,
	).Scan(&userID)
	if userID == "" {
		const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"
		var password string
		for i := 0; i < 32; i++ {
			num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
			if err != nil {
				return "", err
			}
			password += string(letters[num.Int64()])
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}

		err = r.DB.QueryRow(
			ctx,
			"INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id",
			userInfo.Email, userInfo.Email, hashedPassword,
		).Scan(&userID)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	claims := &jwt.StandardClaims{
		Issuer:  "TeyvatPlanner",
		Subject: userID,
	}

	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshToken, err := unsignedToken.SignedString([]byte(os.Getenv("TP_REFRESH_TOKEN_SECRET")))
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

// DiscordLogin is the resolver for the discordLogin field.
func (r *mutationResolver) DiscordLogin(ctx context.Context, token string) (string, error) {
	conf := &oauth2.Config{
        ClientID: os.Getenv("DISCORD_CLIENT_ID"),
        ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
        Endpoint: oauth2.Endpoint{
            AuthURL: "https://discord.com/oauth2/authorize",
            TokenURL: "https://discord.com/api/oauth2/token",
        },
        RedirectURL: "YOUR_REDIRECT_URL", 
    }

	oauth2Token := &oauth2.Token{
        AccessToken: token,
    }

	client := conf.Client(ctx, oauth2Token)

	resp, err := client.Get("https://discord.com/api/users/@me")
    if err != nil {
        return "", fmt.Errorf("Error fetching user information")
    }
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("Error reading response")
    }
	
	var userInfo struct {
        Email string `json:"email"`
    }
	
    err = json.Unmarshal(data, &userInfo) 
	if err != nil {
        return "", fmt.Errorf("Error parsing user information")
    }

    if userInfo.Email == "" {
        return "", fmt.Errorf("Email not found in user info")
    }

	var userID string
	err = r.DB.QueryRow(
		ctx,
		"SELECT id FROM users WHERE email = $1",
		userInfo.Email,
	).Scan(&userID)
	if userID == "" {
		const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"
		var password string
		for i := 0; i < 32; i++ {
			num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
			if err != nil {
				return "", err
			}
			password += string(letters[num.Int64()])
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}

		err = r.DB.QueryRow(
			ctx,
			"INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id",
			userInfo.Email, userInfo.Email, hashedPassword,
		).Scan(&userID)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	claims := &jwt.StandardClaims{
		Issuer:  "TeyvatPlanner",
		Subject: userID,
	}

	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshToken, err := unsignedToken.SignedString([]byte(os.Getenv("TP_REFRESH_TOKEN_SECRET")))
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	authUser := auth.ForContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("Access Denied")
	}

	var user model.User
	err := r.DB.QueryRow(
		ctx,
		"SELECT id, username, email, password, createdAt FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)

	commissionRows, err := r.DB.Query(ctx, "SELECT id, name, category, completed FROM commissions WHERE userID = $1", id)
	if err != nil {
		return nil, err
	}

	defer commissionRows.Close()
	for commissionRows.Next() {
		var commission model.Commission
		if err := commissionRows.Scan(&commission.ID, &commission.Name, &commission.Category, &commission.Completed); err != nil {
			return nil, err
		}
		user.Commissions = append(user.Commissions, &commission)
	}

	if err := commissionRows.Err(); err != nil {
		return nil, err
	}

	domainRows, err := r.DB.Query(ctx, "SELECT id, name, completed, createdAt FROM domains WHERE userID = $1", id)
	if err != nil {
		return nil, err
	}

	defer domainRows.Close()
	for domainRows.Next() {
		var domain model.Domain
		if err := domainRows.Scan(&domain.ID, &domain.Name, &domain.Completed, &domain.CreatedAt); err != nil {
			return nil, err
		}
		user.Domains = append(user.Domains, &domain)
	}

	if err := domainRows.Err(); err != nil {
		return nil, err
	}

	weeklyBossRows, err := r.DB.Query(ctx, "SELECT id, name, completed, createdAt FROM weeklybosses WHERE userID = $1", id)
	if err != nil {
		return nil, err
	}

	defer weeklyBossRows.Close()
	for weeklyBossRows.Next() {
		var weeklyBoss model.WeeklyBoss
		if err := weeklyBossRows.Scan(&weeklyBoss.ID, &weeklyBoss.Name, &weeklyBoss.Completed, &weeklyBoss.CreatedAt); err != nil {
			return nil, err
		}
		user.WeeklyBosses = append(user.WeeklyBosses, &weeklyBoss)
	}

	if err := weeklyBossRows.Err(); err != nil {
		return nil, err
	}

	randomQuestRows, err := r.DB.Query(ctx, "SELECT id, name, longitude, latitude, completed, createdAt FROM randomquests WHERE userID = $1", id)
	if err != nil {
		return nil, err
	}

	defer randomQuestRows.Close()
	for randomQuestRows.Next() {
		var randomQuest model.RandomQuest
		if err := randomQuestRows.Scan(&randomQuest.ID, &randomQuest.Name, &randomQuest.Longitude, &randomQuest.Latitude, &randomQuest.Completed, &randomQuest.CreatedAt); err != nil {
			return nil, err
		}
		user.RandomQuests = append(user.RandomQuests, &randomQuest)
	}

	if err := randomQuestRows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

// Commission is the resolver for the commission field.
func (r *queryResolver) Commission(ctx context.Context, id string) (*model.Commission, error) {
	authUser := auth.ForContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("Access Denied")
	}

	var commission model.Commission
	err := r.DB.QueryRow(
		ctx,
		"SELECT id, name, category, completed, createdAt FROM commissions WHERE id = $1 AND userID = $2",
		id, authUser.ID,
	).Scan(&commission.ID, &commission.Name, &commission.Category, &commission.Completed, &commission.CreatedAt)

	if err != nil {
		return nil, err
	}

	commission.User = &model.User{
		ID:       authUser.ID,
		Username: authUser.Username,
		Email:    authUser.Email,
	}

	return &commission, nil
}

// Commissions is the resolver for the commissions field.
func (r *queryResolver) Commissions(ctx context.Context) ([]*model.Commission, error) {
	panic(fmt.Errorf("not implemented: Commissions - commissions"))
}

// Domain is the resolver for the domain field.
func (r *queryResolver) Domain(ctx context.Context, id string) (*model.Domain, error) {
	authUser := auth.ForContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("Access Denied")
	}

	var domain model.Domain
	err := r.DB.QueryRow(
		ctx,
		"SELECT id, name, completed, createdAt FROM domains WHERE id = $1 AND userID = $2",
		id, authUser.ID,
	).Scan(&domain.ID, &domain.Name, &domain.Completed, &domain.CreatedAt)

	if err != nil {
		return nil, err
	}

	domain.User = &model.User{
		ID:       authUser.ID,
		Username: authUser.Username,
		Email:    authUser.Email,
	}

	return &domain, nil
}

// Domains is the resolver for the domains field.
func (r *queryResolver) Domains(ctx context.Context) ([]*model.Domain, error) {
	panic(fmt.Errorf("not implemented: Domains - domains"))
}

// WeeklyBoss is the resolver for the weeklyBoss field.
func (r *queryResolver) WeeklyBoss(ctx context.Context, id string) (*model.WeeklyBoss, error) {
	authUser := auth.ForContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("Access Denied")
	}

	var weeklyBoss model.WeeklyBoss
	err := r.DB.QueryRow(
		ctx,
		"SELECT id, name, completed, createdAt FROM weeklybosses WHERE id = $1 AND userID = $2",
		id, authUser.ID,
	).Scan(&weeklyBoss.ID, &weeklyBoss.Name, &weeklyBoss.Completed, &weeklyBoss.CreatedAt)

	if err != nil {
		return nil, err
	}

	weeklyBoss.User = &model.User{
		ID:       authUser.ID,
		Username: authUser.Username,
		Email:    authUser.Email,
	}

	return &weeklyBoss, nil
}

// WeeklyBosses is the resolver for the weeklyBosses field.
func (r *queryResolver) WeeklyBosses(ctx context.Context) ([]*model.WeeklyBoss, error) {
	panic(fmt.Errorf("not implemented: WeeklyBosses - weeklyBosses"))
}

// RandomQuest is the resolver for the randomQuest field.
func (r *queryResolver) RandomQuest(ctx context.Context, id string) (*model.RandomQuest, error) {
	authUser := auth.ForContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("Access Denied")
	}

	var randomQuest model.RandomQuest
	err := r.DB.QueryRow(
		ctx,
		"SELECT id, name, longitude, latitude, completed, createdAt FROM randomquests WHERE id = $1 AND userID = $2",
		id, authUser.ID,
	).Scan(&randomQuest.ID, &randomQuest.Name, &randomQuest.Longitude, &randomQuest.Latitude, &randomQuest.Completed, &randomQuest.CreatedAt)

	if err != nil {
		return nil, err
	}

	randomQuest.User = &model.User{
		ID:       authUser.ID,
		Username: authUser.Username,
		Email:    authUser.Email,
	}

	return &randomQuest, nil
}

// RandomQuests is the resolver for the randomQuests field.
func (r *queryResolver) RandomQuests(ctx context.Context) ([]*model.RandomQuest, error) {
	panic(fmt.Errorf("not implemented: RandomQuests - randomQuests"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented: Users - users"))
}
