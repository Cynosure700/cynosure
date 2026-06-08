package storage

import (
	"context"
	"database/sql"
)

func (s *Store) GetUserProfile(ctx context.Context, userID string) (UserProfile, bool, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT user_id, profile_json, updated_at
		FROM user_profiles WHERE user_id = ?
	`, userID)
	var p UserProfile
	if err := row.Scan(&p.UserID, &p.ProfileJSON, &p.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return UserProfile{}, false, nil
		}
		return UserProfile{}, false, err
	}
	return p, true, nil
}

func (s *Store) UpsertUserProfile(ctx context.Context, profile UserProfile) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO user_profiles (user_id, profile_json, updated_at)
		VALUES (?, ?, NOW(6))
		ON DUPLICATE KEY UPDATE profile_json = VALUES(profile_json)
	`, profile.UserID, profile.ProfileJSON)
	return err
}
