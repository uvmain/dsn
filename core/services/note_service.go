package services

import (
	"database/sql"
	"dsn/core/types"
	"fmt"
	"strings"
)

type NoteService struct {
	db *sql.DB
}

func NewNoteService(db *sql.DB) *NoteService {
	return &NoteService{db: db}
}

func (s *NoteService) Create(userID int, req types.CreateNoteRequest) (*types.Note, error) {
	query := `
		INSERT INTO notes (user_id, title, content, color, pinned, archived, order_position) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING id, created_at, updated_at
	`

	color := req.Color
	if color == "" {
		color = "#ffffff"
	}

	var note types.Note
	err := s.db.QueryRow(query, userID, req.Title, req.Content, color, req.Pinned, req.Archived, req.Order).
		Scan(&note.ID, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, err
	}

	note.UserID = userID
	note.Title = req.Title
	note.Content = req.Content
	note.Color = color
	note.Pinned = req.Pinned
	note.Archived = req.Archived
	note.Order = req.Order

	return &note, nil
}

func (s *NoteService) GetByID(id, userID int) (*types.Note, error) {
	query := `
		SELECT id, user_id, title, content, color, pinned, archived, order_position, created_at, updated_at 
		FROM notes 
		WHERE id = ? AND user_id = ?
	`

	var note types.Note
	err := s.db.QueryRow(query, id, userID).Scan(
		&note.ID, &note.UserID, &note.Title, &note.Content, &note.Color,
		&note.Pinned, &note.Archived, &note.Order, &note.CreatedAt, &note.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Load tags
	tags, err := s.getNoteTags(note.ID)
	if err != nil {
		return nil, err
	}
	note.Tags = tags

	return &note, nil
}

func (s *NoteService) GetByUserID(userID int, includeArchived bool) ([]types.Note, error) {
	query := `
		SELECT id, user_id, title, content, color, pinned, archived, order_position, created_at, updated_at 
		FROM notes 
		WHERE user_id = ?
	`
	args := []interface{}{userID}

	if !includeArchived {
		query += " AND archived = FALSE"
	}

	query += " ORDER BY pinned DESC, order_position ASC, updated_at DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]types.Note, 0)
	for rows.Next() {
		var note types.Note
		err := rows.Scan(
			&note.ID, &note.UserID, &note.Title, &note.Content, &note.Color,
			&note.Pinned, &note.Archived, &note.Order, &note.CreatedAt, &note.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Load tags for each note
		tags, err := s.getNoteTags(note.ID)
		if err != nil {
			return nil, err
		}
		note.Tags = tags

		notes = append(notes, note)
	}

	return notes, nil
}

func (s *NoteService) Update(id, userID int, req types.UpdateNoteRequest) (*types.Note, error) {
	// Build dynamic update query
	var setParts []string
	var args []interface{}

	if req.Title != nil {
		setParts = append(setParts, "title = ?")
		args = append(args, *req.Title)
	}
	if req.Content != nil {
		setParts = append(setParts, "content = ?")
		args = append(args, *req.Content)
	}
	if req.Color != nil {
		setParts = append(setParts, "color = ?")
		args = append(args, *req.Color)
	}
	if req.Pinned != nil {
		setParts = append(setParts, "pinned = ?")
		args = append(args, *req.Pinned)
	}
	if req.Archived != nil {
		setParts = append(setParts, "archived = ?")
		args = append(args, *req.Archived)
	}
	if req.Order != nil {
		setParts = append(setParts, "order_position = ?")
		args = append(args, *req.Order)
	}

	if len(setParts) == 0 {
		return s.GetByID(id, userID)
	}

	setParts = append(setParts, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id, userID)

	query := fmt.Sprintf(`
		UPDATE notes 
		SET %s 
		WHERE id = ? AND user_id = ?
	`, strings.Join(setParts, ", "))

	result, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("note with id %d not found", id)
	}

	return s.GetByID(id, userID)
}

func (s *NoteService) UpdateOrder(userID int, noteOrders map[int]int) error {
	// Begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update order for each note
	for noteID, order := range noteOrders {
		_, err := tx.Exec(`
			UPDATE notes 
			SET order_position = ?, updated_at = CURRENT_TIMESTAMP 
			WHERE id = ? AND user_id = ?
		`, order, noteID, userID)
		if err != nil {
			return err
		}
	}

	// Commit transaction
	return tx.Commit()
}

func (s *NoteService) Delete(id, userID int) error {
	query := "DELETE FROM notes WHERE id = ? AND user_id = ?"
	result, err := s.db.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("note with id %d not found", id)
	}

	return nil
}

func (s *NoteService) Search(userID int, query string) ([]types.Note, error) {
	searchQuery := `
		SELECT id, user_id, title, content, color, pinned, archived, order_position, created_at, updated_at 
		FROM notes 
		WHERE user_id = ? AND (title LIKE ? OR content LIKE ?)
		ORDER BY pinned DESC, order_position ASC, updated_at DESC
	`

	searchTerm := "%" + query + "%"
	rows, err := s.db.Query(searchQuery, userID, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]types.Note, 0)
	for rows.Next() {
		var note types.Note
		err := rows.Scan(
			&note.ID, &note.UserID, &note.Title, &note.Content, &note.Color,
			&note.Pinned, &note.Archived, &note.Order, &note.CreatedAt, &note.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Load tags for each note
		tags, err := s.getNoteTags(note.ID)
		if err != nil {
			return nil, err
		}
		note.Tags = tags

		notes = append(notes, note)
	}

	return notes, nil
}

func (s *NoteService) TogglePin(id, userID int, pinned bool) (*types.Note, error) {
	query := `
		UPDATE notes 
		SET pinned = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ? AND user_id = ?
	`

	result, err := s.db.Exec(query, pinned, id, userID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("note with id %d not found", id)
	}

	return s.GetByID(id, userID)
}

func (s *NoteService) ToggleArchive(id, userID int, archived bool) (*types.Note, error) {
	query := `
		UPDATE notes 
		SET archived = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ? AND user_id = ?
	`

	result, err := s.db.Exec(query, archived, id, userID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("note with id %d not found", id)
	}

	return s.GetByID(id, userID)
}

func (s *NoteService) getNoteTags(noteID int) ([]types.Tag, error) {
	query := `
		SELECT t.id, t.name, t.color, t.created_at
		FROM tags t
		JOIN note_tags nt ON t.id = nt.tag_id
		WHERE nt.note_id = ?
		ORDER BY t.name
	`

	rows, err := s.db.Query(query, noteID)
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

	// Check for iteration errors
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}
