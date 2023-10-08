package storage

type BaseStorage interface {
	AddNote() error
	EditNote() error
	GetNotes() error
}
