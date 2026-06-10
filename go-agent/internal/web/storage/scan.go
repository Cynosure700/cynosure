package storage

import "database/sql"

func scanUser(row interface{ Scan(dest ...any) error }) (User, error) {
	var user User
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.MemoryEnabled, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return User{}, err
	}
	return user, nil
}

func scanSkills(rows *sql.Rows) ([]Skill, error) {
	var skills []Skill
	for rows.Next() {
		var skill Skill
		if err := rows.Scan(&skill.ID, &skill.UserID, &skill.Name, &skill.Slug, &skill.Description, &skill.Content, &skill.Status, &skill.CreatedAt, &skill.UpdatedAt); err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}
	return skills, rows.Err()
}
