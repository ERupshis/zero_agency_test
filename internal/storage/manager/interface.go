package manager

type BaseStorageManager interface {
	AddNote() error
	EditNote() error
	GetNotes() error
}
