package storage

import "context"

func (s *Store) CreateSkill(ctx context.Context, skill Skill) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO skills (id, user_id, name, slug, description, content, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`, skill.ID, skill.UserID, skill.Name, skill.Slug, skill.Description, skill.Content, skill.Status)
	return err
}

func (s *Store) ListSkillsByUser(ctx context.Context, userID string) ([]Skill, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, user_id, name, slug, description, content, status, created_at, updated_at
		FROM skills WHERE user_id = ? ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSkills(rows)
}

func (s *Store) ListEnabledSkillsByUser(ctx context.Context, userID string) ([]Skill, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, user_id, name, slug, description, content, status, created_at, updated_at
		FROM skills WHERE user_id = ? AND status = 'enabled' ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSkills(rows)
}

func (s *Store) GetSkillByID(ctx context.Context, skillID string) (Skill, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, user_id, name, slug, description, content, status, created_at, updated_at
		FROM skills WHERE id = ?
	`, skillID)
	var skill Skill
	if err := row.Scan(&skill.ID, &skill.UserID, &skill.Name, &skill.Slug, &skill.Description, &skill.Content, &skill.Status, &skill.CreatedAt, &skill.UpdatedAt); err != nil {
		return Skill{}, err
	}
	return skill, nil
}

func (s *Store) UpdateSkill(ctx context.Context, skill Skill) error {
	_, err := s.DB.ExecContext(ctx, `
		UPDATE skills
		SET name = ?, slug = ?, description = ?, content = ?, status = ?, updated_at = NOW()
		WHERE id = ?
	`, skill.Name, skill.Slug, skill.Description, skill.Content, skill.Status, skill.ID)
	return err
}

func (s *Store) DeleteSkill(ctx context.Context, skillID string) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM skills WHERE id = ?`, skillID)
	return err
}
