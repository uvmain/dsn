package services

import (
	"database/sql"
	"dsn/core/database"
	"dsn/core/types"
	"fmt"
)

type TagService struct {
	db *sql.DB
}

func NewTagService() *TagService {
	return &TagService{db: database.DB}
}

func (s *TagService) Create(req types.CreateTagRequest) (*types.Tag, error) {
	query := `
		INSERT INTO tags (name, color) 
		VALUES (?, ?)
		RETURNING id, created_at
	`

	color := req.Color
	if color == "" {
		color = "#e0e0e0"
	}

	var tag types.Tag
	err := s.db.QueryRow(query, req.Name, color).Scan(&tag.ID, &tag.CreatedAt)
	if err != nil {
		return nil, err
	}

	tag.Name = req.Name
	tag.Color = color

	return &tag, nil
}

func (s *TagService) GetAll() ([]types.Tag, error) {
	query := `
		SELECT id, name, color, created_at 
		FROM tags 
		ORDER BY name ASC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]types.Tag, 0)
	for rows.Next() {
		var tag types.Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Color, &tag.CreatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (s *TagService) GetByID(id int) (*types.Tag, error) {
	query := `
		SELECT id, name, color, created_at 
		FROM tags 
		WHERE id = ?
	`

	var tag types.Tag
	err := s.db.QueryRow(query, id).Scan(&tag.ID, &tag.Name, &tag.Color, &tag.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func (s *TagService) Update(id int, req types.UpdateTagRequest) (*types.Tag, error) {
	var setParts []string
	var args []interface{}

	if req.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *req.Name)
	}
	if req.Color != nil {
		setParts = append(setParts, "color = ?")
		args = append(args, *req.Color)
	}

	if len(setParts) == 0 {
		return s.GetByID(id)
	}

	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE tags 
		SET %s 
		WHERE id = ?
	`, fmt.Sprintf("%s", setParts[0]))

	for i := 1; i < len(setParts); i++ {
		query = fmt.Sprintf("%s, %s", query, setParts[i])
	}

	result, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("tag with id %d not found", id)
	}

	return s.GetByID(id)
}

func (s *TagService) Delete(id int) error {
	query := "DELETE FROM tags WHERE id = ?"
	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tag with id %d not found", id)
	}

	return nil
}

func (s *TagService) AssignToNote(noteID, tagID int) error {
	query := `
		INSERT OR IGNORE INTO note_tags (note_id, tag_id) 
		VALUES (?, ?)
	`
	_, err := s.db.Exec(query, noteID, tagID)
	return err
}

func (s *TagService) RemoveFromNote(noteID, tagID int) error {
	query := "DELETE FROM note_tags WHERE note_id = ? AND tag_id = ?"
	result, err := s.db.Exec(query, noteID, tagID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tag %d not assigned to note %d", tagID, noteID)
	}

	return nil
}

func (s *TagService) SetNoteTags(noteID int, tagIDs []int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM note_tags WHERE note_id = ?", noteID)
	if err != nil {
		return err
	}

	for _, tagID := range tagIDs {
		_, err = tx.Exec("INSERT INTO note_tags (note_id, tag_id) VALUES (?, ?)", noteID, tagID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
